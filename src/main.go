package main

import (
	"flag"
	"fmt"
	"net/http"
	data "src/data_manager"
	strategy "src/strategy"
	randomWalk "src/strategy/random_walk"

	"github.com/golang/glog"
)

var solver strategy.Strategy

func Ready() bool {
	scenario := data.GetInstance()
	return scenario.Ready()
}
func P1Start(pilots []string) map[string][]string {
	scenario := data.GetInstance()
	scenario.HandleP1(pilots)

	return solver.ExecuteP1()
}
func P2Start(param data.LoadData) map[string][]string {
	scenario := data.GetInstance()
	scenario.HandleP2(param)

	return solver.ExecuteP2()
}

func main() {
	flag.Parse()
	defer glog.Flush()

	// 选择策略
	solver = &randomWalk.RandomWalk{}

	InstallHttpHandlers()

	addr := fmt.Sprintf("0.0.0.0:%d", 3355)
	glog.Infof("@ Start Program: %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		glog.Errorf("ListenAndServe met err %v", err)
	}
}
