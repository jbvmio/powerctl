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

func getAllPods() k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	results, err := rc.GetPods("")
	h(err)
	return results
}

func searchPods(search []string, exact bool) k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	rc.ExactMatches(exact)
	results, err := rc.GetPods(search[:]...)
	h(err)
	return results
}

func nodesToPods(args []string) []k8s.XD {
	var xdata []k8s.XD
	allPods := getAllPods()
	for _, pod := range allPods.XData {
		for _, node := range args {
			if pod.NodeName == node {
				xdata = append(xdata, pod)
			}
		}
	}
	return xdata
}

func rsToPods(args []string) []k8s.XD {
	var xdata []k8s.XD
	allPods := getAllPods()
	for _, pod := range allPods.XData {
		for _, rs := range args {
			for _, p := range pod.OwnerReferences {
				if p.OwnerName == rs {
					xdata = append(xdata, pod)
				}
			}

		}
	}
	return xdata
}

func makePrintPods(xdata []k8s.XD) {
	var pods []Pod
	podChan := make(chan Pod, 100)
	for _, x := range xdata {
		go makePods(x, podChan)
	}
	for i := 0; i < len(xdata); i++ {
		pod := <-podChan
		pods = append(pods, pod)
	}
	sortSlice(pods)
	formatTable(pods)
}

func getPodStatus(cs *k8s.Status) (string, string) {
	var status = cs.Phase
	var message = "PodScheduled"
	if status != "Succeeded" {
		for _, c := range cs.ContainerStatuses {
			if !c.Ready {
				if c.Terminated != nil {
					status = c.Terminated.Reason
				}
				if c.Waiting != nil {
					status = c.Waiting.Reason
				}
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
