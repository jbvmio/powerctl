package cmd

import (
	"time"

	"github.com/jbvmio/k8s"
)

// Node struct definition:
type Node struct {
	NodeName       string
	Status         string
	InternalIP     string
	KernelVersion  string
	KubeletVersion string
	Runtime        string
	Message        string
	Age            string
}

func makeNodes(x k8s.XD, nodeChan chan Node) {

	status, message, ip, _ := getNodeStatus(x.Status)
	node := Node{
		NodeName:       x.Name,
		Status:         status,
		InternalIP:     ip,
		KernelVersion:  x.KernelVersion,
		KubeletVersion: x.KubeletVersion,
		Runtime:        x.ContainerRuntimeVersion,
		Message:        message,
		Age:            roundToDays(time.Now().Sub(x.Created)),
	}
	nodeChan <- node

}

func getNodeStatus(cs *k8s.Status) (string, string, string, string) {
	var status string
	var message string
	var ip string
	var hostname string
	for _, c := range cs.Addresses {
		if c.Type == "InternalIP" {
			ip = c.Address
		}
		if c.Type == "Hostname" {
			hostname = c.Address
		}
	}
	for _, c := range cs.Conditions {
		message = c.Reason
		status = c.Type
		if c.Status == "True" {
			message = c.Reason
			status = c.Type
			return status, message, ip, hostname
		}
	}
	return status, message, ip, hostname
}

func makePrintNodes(xdata []k8s.XD) {
	var nodes []Node
	nodeChan := make(chan Node, 100)
	for _, x := range xdata {
		go makeNodes(x, nodeChan)
	}
	for i := 0; i < len(xdata); i++ {
		node := <-nodeChan
		nodes = append(nodes, node)
	}
	sortSlice(nodes)
	formatTable(nodes)
}
