package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
	//v1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"strings"
	"time"
)

type ExecutorHandler struct {
	Redis        common.RedisClient
	Source       string
	clientSet    *kubernetes.Clientset
	ongoingTasks map[uuid.UUID]operation
}

type operation struct {
	ID     uuid.UUID
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewExecutorHandler(redisClient common.RedisClient, consumerGroup string) (*ExecutorHandler, error) {
	// Create in-cluster config (automatically uses the service account token and API server address)
	config, err := rest.InClusterConfig()
	if err != nil {
		return &ExecutorHandler{}, err
	}

	// Create the clientset (used to interact with the Kubernetes API)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &ExecutorHandler{}, err
	}

	taskMap := make(map[uuid.UUID]operation)
	return &ExecutorHandler{
		Redis:        redisClient,
		Source:       consumerGroup,
		clientSet:    clientset,
		ongoingTasks: taskMap,
	}, nil
}

func (h *ExecutorHandler) Message(msg *message.MessageWithRedisOperations) error {
	switch msg.Message.Type {
	case message.ExperimentStartRequest:
		err := h.handleExperimentStartRequest(msg)

		// ERROR SOMEWHERE AFTER HERE
		slog.Error("handleExperimentStart", "error", err)

		respMsg := message.New(
			message.WithSource(h.Source),
			message.WithType(message.ExperimentStop),
		)

		// if message isn't processed, we can't extract the task info, so return early
		if errors.Is(err, &message.MessageNotProcessedError{}) {
			respMsg.Contents = err
			respMsg.ResponseID = msg.Message.ID
			_ = h.Redis.SendMessageToStream(respMsg, streams.ExperimentControl)
			return err
		}

		// should never not be ok, MessageNotProcessedError should be invoked first
		task, _ := common.GenericUnmarshal[tasks.Task](msg.Message.Contents)
		// copy task
		respTask := task
		// send correct message back to stream
		switch {
		case err == nil:
			respTask.Status = tasks.StatusCompleted
		case errors.Is(err, context.Canceled):
			respTask.Status = tasks.StatusCanceled
		case errors.Is(err, context.DeadlineExceeded):
			respTask.Status = tasks.StatusTimedOut
		default:
			respTask.Status = tasks.StatusFailed
			respTask.Details = map[string]interface{}{
				"error": &tasks.TaskError{
					ID:     task.ID,
					Type:   task.Type,
					Reason: err.Error(),
				},
			}
		}
		respMsg.Contents = respTask
		_ = h.Redis.SendMessageToStream(respMsg, streams.ExperimentControl)

		// return err for logging
		return err

	case message.ExperimentStopRequest:
		experimentStopRequestData, ok := msg.Message.Contents.(tasks.Task)
		if !ok {
			return &message.MessageNotProcessedError{
				ID:     msg.Message.ID,
				Type:   msg.Message.Type,
				Reason: "Unable to extract Experiment Stop Request Data from message",
			}
		}

		// Stop experiment

		// Add jobstop message back to the ExperimentControl stream
		returnMsg := message.New(
			message.WithType(message.ExperimentStop),
			message.WithResponseID(msg.Message.ID),
			message.WithSource(h.Source),
			message.WithContents(tasks.Task{
				ID: experimentStopRequestData.ID,
			}),
		)
		h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl)
		return nil
	default:
		return &message.MessageNotProcessedError{
			ID:     msg.Message.ID,
			Type:   msg.Message.Type,
			Reason: "ExecutorHandler not configured to handle type",
		}
	}
}

// entrypoint for messagestartrequests
func (h *ExecutorHandler) handleExperimentStartRequest(msg *message.MessageWithRedisOperations) error {
	// ensure task contents is of correct type
	task, err := common.GenericUnmarshal[tasks.Task](msg.Message.Contents)
	if err != nil {
		return &message.MessageNotProcessedError{
			ID:     msg.Message.ID,
			Type:   msg.Message.Type,
			Reason: "Unable to extract Experiment Start Task Data from message",
		}
	}

	// validate the contents
	errs, ok := task.Validate()
	if !ok {
		return &message.MessageNotProcessedError{
			ID:     msg.Message.ID,
			Type:   msg.Message.Type,
			Reason: fmt.Sprintf("Unable to validate Request Data from message. Validation errors %s", strings.Join(errs, ", ")),
		}
	}

	// Ack message so it doesn't get reprocessed
	if err := msg.Ack(); err != nil {
		return err
	}

	// create context for the experiment operation
	var ctx context.Context
	var cancel context.CancelFunc
	if task.Timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(task.Timeout))
	}
	defer cancel()

	// add to ongoing task map
	h.ongoingTasks[msg.Message.ID] = operation{
		ID:     msg.Message.ID,
		Ctx:    ctx,
		Cancel: cancel,
	}

	// send jobStart message
	returnMsg := message.New(
		message.WithType(message.ExperimentStart),
		message.WithResponseID(msg.Message.ID),
		message.WithSource(h.Source),
		message.WithContents(&tasks.Task{
			ID:     task.ID,
			Status: tasks.StatusRunning,
		}),
	)
	if err := h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl); err != nil {
		return err
	}

	// run task
	err = h.runtask(ctx, task)

	//remove task from ongoing tasks
	delete(h.ongoingTasks, task.ID)

	return err
}

// handles task Operation contexts to ensure that they time out correctly
func (h *ExecutorHandler) HandleOperationContexts() {
	for _, op := range h.ongoingTasks {
		if err := op.Ctx.Err(); err != nil {
			delete(h.ongoingTasks, op.ID)
			continue
		}
		deadline, ok := op.Ctx.Deadline()
		if !ok {
			continue
		}

		if time.Now().After(deadline) {
			op.Cancel()
			delete(h.ongoingTasks, op.ID)
		}
	}
}

func (h *ExecutorHandler) runtask(ctx context.Context, t tasks.Task) error {
	switch t.Type {
	case tasks.TaskDeletePod:
		deleteOpts := t.Options().Build()
		_, ok := deleteOpts.(metav1.DeleteOptions)
		if !ok {
			return fmt.Errorf("unable to convert task options to DeleteOptions")
		}
		return h.listPods()
		//return h.deletePod(ctx, t.Target, t.NameSpace, opts)
	default:
		return fmt.Errorf("unable to run task of type %v, not configured", t.Type)
	}
}

func (h *ExecutorHandler) deletePod(ctx context.Context, target string, namespace string, options metav1.DeleteOptions) error {
	return h.clientSet.CoreV1().Pods(namespace).Delete(ctx, target, options)
}

func (h *ExecutorHandler) listPods() error {
	// List pods in the "default" namespace
	pods, err := h.clientSet.CoreV1().Pods("chaos-kube").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pods in the 'chaos-kube' namespace:")
	for _, pod := range pods.Items {
		fmt.Println(" -", pod.Name)
	}
	return nil
}

// func (h *ExecutorHandler) evictPod(target string, namespace string, options map[string]interface{} ) error {
// 	ctx := context.TODO()
// 	err := h.clientSet.CoreV1().Pods(namespace).EvictV1(ctx, &v1.Eviction{})
// 	return err
// }
