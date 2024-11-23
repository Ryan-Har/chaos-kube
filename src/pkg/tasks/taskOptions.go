package tasks

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskOptions interface {
	Build() interface{}
}

type DeletePodOptions struct {
	GracePeriodSeconds *int64
}

func (d DeletePodOptions) Build() interface{} {
	return metav1.DeleteOptions{
		GracePeriodSeconds: d.GracePeriodSeconds,
	}
}

func defaultDeletePodOptions() DeletePodOptions {
	gp := int64(0)
	return DeletePodOptions{
		GracePeriodSeconds: &gp,
	}
}
