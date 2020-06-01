package main

import (
	"encoding/json"
	"sort"

	"github.com/golang/glog"
)

func (ths *Player) P1(pilots []string) (string, error) {
	glog.Infof("P1 start: %v\n", pilots)
	ths.Pilots = pilots
	ths.LoadData()
	return ths.arrange()
}

func (ths *Player) P2(params LoadData) (string, error) {
	glog.Infof("P2 start: %v\n", ths.Pilots)
	ths.UpdateParams(params)
	return ths.arrange()
}

func (ths *Player) arrange() (string, error) {
	sort.Slice(ths.AppList, func(i, j int) bool {
		cnt1, cnt2 := ths.AppList[i].SidecarCount, ths.AppList[j].SidecarCount
		len1, len2 := len(ths.AppList[i].ServiceList), len(ths.AppList[j].ServiceList)
		if cnt1 == cnt2 {
			if len1 == len1 {
				return ths.AppList[i].AppName < ths.AppList[j].AppName

			}
			return len1 < len2
		}
		return cnt1 > cnt2
	})
	sz := len(ths.Pilots)
	result := make(map[string][]string)
	for idx, app := range ths.AppList {
		pilot := ths.Pilots[idx%sz]
		result[pilot] = append(result[pilot], app.AppName)
	}
	b, err := json.Marshal(result)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	return string(b), nil
}
