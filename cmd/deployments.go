package cmd

import (
	"time"

	"github.com/jbvmio/k8s"
)

// Deployment struct defined here
type Deployment struct {
	Deployment string
	Desired    int
	Deployed   int
	Ready      int
	Age        string
}

func getAllDeploys() k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	results, err := rc.GetDeployments("")
	h(err)
	return results
}

func makePrintDeploys(xdata []k8s.XD) {
	var deployment []Deployment
	deployChan := make(chan Deployment, 100)
	for _, x := range xdata {
		go makeDeployment(x, deployChan)
	}
	for i := 0; i < len(xdata); i++ {
		deploy := <-deployChan
		deployment = append(deployment, deploy)
	}
	sortSlice(deployment)
	formatTable(deployment)
}

func makeDeployment(x k8s.XD, deployChan chan Deployment) {

	deploy := Deployment{
		Deployment: x.Name,
		Desired:    x.Replicas,
		Deployed:   x.UpdatedReplicas,
		Ready:      x.ReadyReplicas,
		Age:        roundToDays(time.Now().Sub(x.Created)),
	}
	deployChan <- deploy
}

func searchDeploys(args []string) k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	results, err := rc.GetDeployments(args[:]...)
	h(err)
	return results
}

func rsToDeployments(args []string) []k8s.XD {
	var xdata []k8s.XD
	var uids []string
	var deployNames []string
	replicaset := searchRS(args)
	for _, rs := range replicaset.XData {
		for _, dep := range rs.OwnerReferences {
			deployNames = append(deployNames, dep.OwnerName)
			uids = append(uids, dep.OwnerUID)
		}
	}
	xd := searchDeploys(filterUnique(deployNames)).XData
	for _, x := range xd {
		for _, uid := range uids {
			if x.UID == uid {
				xdata = append(xdata, x)
			}
		}
	}
	return xdata
}
