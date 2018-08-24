package cmd

import (
	"time"

	"github.com/jbvmio/k8s"
)

// Pod struct defined here
type Pod struct {
	PodName  string
	Status   string
	Restarts int
	PodIP    string
	Message  string
	NodeName string
	Age      string
}

func getPodStatus(cs *k8s.Status) (string, string) {
	var status = cs.Phase
	var message = "PodScheduled"
	if status != "Succeeded" {
		for _, c := range cs.ContainerStatuses {
			if !c.Ready {
				status = c.Reason
			}
		}
	}
	for _, c := range cs.Conditions {
		if c.Status == "False" {
			message = c.Reason
			return status, message
		}
	}
	return status, message
}

func getPodRestarts(cs []k8s.ContainerStatuses) int {
	var count int
	for _, c := range cs {
		count = count + c.RestartCount
	}
	return count
}

func makePods(x k8s.XD, podChan chan Pod) {

	status, message := getPodStatus(x.Status)
	pod := Pod{
		PodName:  x.Name,
		Status:   status,
		Message:  message,
		Restarts: getPodRestarts(x.ContainerStatuses),
		PodIP:    x.PodIP,
		NodeName: x.NodeName,
		Age:      roundToDays(time.Now().Sub(x.Created)),
	}
	podChan <- pod

}

// Backup
/*
	defer func() {
		if r := recover(); r != nil {

		}
	}
*/
