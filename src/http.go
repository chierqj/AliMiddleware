package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	data "src/data_manager"
)

func InstallHttpHandlers() {
	http.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		ready := Ready()
		statusCode := http.StatusOK
		if !ready {
			statusCode = http.StatusInternalServerError
		}
		http.Error(writer, "", statusCode)
	})

	http.HandleFunc("/p1_start", func(writer http.ResponseWriter, request *http.Request) {
		// Parse pilots from req body.
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var pilots struct {
			Pilots []string `json:"pilots"`
		}
		if err := json.Unmarshal(bs, &pilots); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		res := P1Start(pilots.Pilots)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if bs, err := json.Marshal(res); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, string(bs), http.StatusOK)
		}
	})

	http.HandleFunc("/p2_start", func(writer http.ResponseWriter, request *http.Request) {
		// Parse apps and dependencies from req body.
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var param data.LoadData
		if err := json.Unmarshal(bs, &param); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		res := P2Start(param)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if bs, err := json.Marshal(res); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, string(bs), http.StatusOK)
		}
	})
}
