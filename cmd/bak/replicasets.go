package cmd

import (
	"time"

	"github.com/jbvmio/k8s"
)

// ReplicaSet struct defined here
type ReplicaSet struct {
	ReplicaSet string
	Desired    int
	Deployed   int
	Available  int
	Ready      int
	Age        string
}

func getAllRS() k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	results, err := rc.GetRS("")
	h(err)
	return results
}

func searchRS(args []string, exact bool) k8s.Results {
	rc, err := k8s.NewRawClient(false)
	h(err)
	rc.SetNS(targetNamespace)
	rc.ExactMatches(exact)
	results, err := rc.GetRS(args[:]...)
	h(err)
	return results
}

func podsToRS(args []string) []k8s.XD {
	var xdata []k8s.XD
	var rsNames []string
	pods := searchPods(args, true)
	for _, pod := range pods.XData {
		for _, rs := range pod.OwnerReferences {
			rsNames = append(rsNames, rs.OwnerName)
		}
	}
	xdata = searchRS(filterUnique(rsNames), true).XData
	return xdata
}

func deploysToRS(args []string) []k8s.XD {
	var xdata []k8s.XD
	allRS := getAllRS()
	for _, replicaset := range allRS.XData {
		for _, deployment := range args {
			for _, r := range replicaset.OwnerReferences {
				if r.OwnerName == deployment {
					xdata = append(xdata, replicaset)
				}
			}
		}
	}
	return xdata
}

func makePrintRS(xdata []k8s.XD) {
	var rs []ReplicaSet
	rsChan := make(chan ReplicaSet, 100)
	for _, x := range xdata {
		go makeRS(x, rsChan)
	}
	for i := 0; i < len(xdata); i++ {
		replicaset := <-rsChan
		rs = append(rs, replicaset)
	}
	sortSlice(rs)
	formatTable(rs)
}

func makeRS(x k8s.XD, rsChan chan ReplicaSet) {

	rs := ReplicaSet{
		ReplicaSet: x.Name,
		Desired:    x.Replicas,
		Deployed:   x.FullyLabeledReplicas,
		Available:  x.AvailableReplicas,
		Ready:      x.ReadyReplicas,
		Age:        roundToDays(time.Now().Sub(x.Created)),
	}
	rsChan <- rs

}
