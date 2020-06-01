package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	player := &Player{}

	InstallHttpHandlers(player)

	addr := fmt.Sprintf("0.0.0.0:%d", 3355)
	glog.Infof("@ Start Program: %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		glog.Errorf("ListenAndServe met err %v", err)
	}
}
