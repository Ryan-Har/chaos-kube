package handler

import (
	"context"
	"fmt"
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/message"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"github.com/Ryan-Har/chaos-kube/pkg/tasks"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ExecutorHandler struct {
	Redis     common.RedisClient
	Source    string
	clientSet *kubernetes.Clientset
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

	return &ExecutorHandler{
		Redis:     redisClient,
		Source:    consumerGroup,
		clientSet: clientset,
	}, nil
}

func (h *ExecutorHandler) Message(msg *message.Message) error {
	switch msg.Type {
	case message.ExperimentStartRequest:
		experimentId := uuid.New()
		// generate uuid for experiment
		// start experiment
		// send jobStart message
		h.listPods()
		returnMsg := message.New(
			message.WithType(message.ExperimentStart),
			message.WithResponseID(msg.ID),
			message.WithSource(h.Source),
			message.WithContents(&tasks.Task{
				ID: experimentId,
			}),
		)
		err := h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl)
		return err
	case message.ExperimentStopRequest:
		experimentStopRequestData, ok := msg.Contents.(*tasks.Task)
		if !ok {
			return &message.MessageNotProcessedError{
				ID:     msg.ID,
				Type:   msg.Type,
				Reason: "Unable to extract Experiment Stop Request Data from message",
			}
		}

		// Stop experiment

		// Add jobstop message back to the ExperimentControl stream
		returnMsg := message.New(
			message.WithType(message.ExperimentStop),
			message.WithResponseID(msg.ID),
			message.WithSource(h.Source),
			message.WithContents(tasks.Task{
				ID: experimentStopRequestData.ID,
			}),
		)
		h.Redis.SendMessageToStream(returnMsg, streams.ExperimentControl)
		return nil
	default:
		return &message.MessageNotProcessedError{
			ID:     msg.ID,
			Type:   msg.Type,
			Reason: "ExecutorHandler not configured to handle type",
		}
	}
}

func (h *ExecutorHandler) listPods() {
	// List pods in the "default" namespace
	pods, err := h.clientSet.CoreV1().Pods("chaos-kube").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pods in the 'chaos-kube' namespace:")
	for _, pod := range pods.Items {
		fmt.Println(" -", pod.Name)
	}
}
