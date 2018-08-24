package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

func h(err error) {
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}

func stdinAvailable() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func formatTable(k interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator("")

	switch k := k.(type) {
	case []Pod:
		table.SetHeader([]string{"PodName", "Status", "Restarts", "PodIP", "NodeName", "Message", "Age"})
		for _, v := range k {
			k := []string{v.PodName, truncateString(v.Status, 20), strconv.Itoa(v.Restarts), v.PodIP, v.NodeName, truncateString(v.Message, 20), v.Age}
			table.Append(k)
		}
	case []Node:
		table.SetHeader([]string{"NodeName", "Status", "InternalIP", "Version", "Kernel", "Message", "Age"})
		for _, v := range k {
			k := []string{v.NodeName, truncateString(v.Status, 20), v.InternalIP, v.KubeletVersion, v.KernelVersion, v.Message, v.Age}
			table.Append(k)
		}
	}
	fmt.Println()
	table.Render()
	fmt.Println()
}

// k := []string{v.PodName, truncateString(v.Status, 20), strconv.FormatInt(v.Restarts, 10), v.PodIP, v.NodeName, v.Age}

func roundToDays(d time.Duration) string {
	// Establish seconds in each:
	var (
		//year  float64 = 31207680
		//month float64 = 2600640
		//week  float64 = 604800
		day    float64 = 86400
		hour   float64 = 3600
		minute float64 = 60
	)
	secs := d.Round(time.Second).Seconds()
	if d > time.Hour*24 {
		return fmt.Sprintf("%.2f days", secs/day)
	}
	if d > time.Minute*60 {
		return fmt.Sprintf("%.2f hours", secs/hour)
	}
	if d > time.Second*60 {
		return fmt.Sprintf("%.2f minutes", secs/minute)
	}
	return fmt.Sprintf("%.2f seconds", secs)
}

func truncateString(str string, num int) string {
	s := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		s = str[0:num] + "..."
	}
	return s
}

func sortSlice(sl interface{}) {
	switch sl := sl.(type) {
	case []string:
		sort.Slice(sl, func(i, j int) bool {
			return sl[i] < sl[j]
		})
	case []Pod:
		sort.Slice(sl, func(i, j int) bool {
			return sl[i].PodName < sl[j].PodName
		})
	case []Node:
		sort.Slice(sl, func(i, j int) bool {
			return sl[i].NodeName < sl[j].NodeName
		})
	}
}
