package main

import (
	"encoding/json"
	"sort"

	"github.com/golang/glog"
)

var (
	evolution = &Evolution{}
)

/*
 * @ title:				P1
 * @ description: 响应 /api/p1_start 接口，返回pilots的分配方案
 * @ author: 			FatChier
 * @ param:				pilots ["p1","p2",...]string 需要分配的 pilots
 * @ return:			string: {"p1":["app1","app2",...],"p2":[]...} 分配方案对应的json字符串
 * @ 							error: error or nil
 */
func (ths *Player) P1(pilots []string) (string, error) {
	glog.Infof("P1 start: %v\n", pilots)
	ths.Pilots = pilots
	ths.PilotsSize = len(pilots)
	ths.LoadData()

	evolution.Run(ths)

	return ths.arrange()
}

/*
 * @ title:				P2
 * @ description: 响应 /api/p2_start 接口，更新AppList，返回pilots的分配方案
 * @ author: 			FatChier
 * @ param:				param: LoadData {
 *	 									Apps map[string]int            `json:"apps"`
 *										Dependencies map[string]map[string]int `json:"dependencies"`
 *								}
 * @ return:			string: {"p1":["app1","app2",...],"p2":[]...} 分配方案对应的json字符串
 * @ 							error: error or nil
 */
func (ths *Player) P2(params LoadData) (string, error) {
	glog.Infof("P2 start: %v\n", ths.Pilots)
	ths.UpdateParams(params)
	return ths.arrange()
}

/*
 * @ title:				arrange
 * @ description: 安排pilots算法
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			string: {"p1":["app1","app2",...],"p2":[]...} 分配方案对应的json字符串
 * @ 							error: error or nil
 */
func (ths *Player) arrange() (string, error) {
	sort.Slice(ths.AppList, func(i, j int) bool {
		x1 := ths.AppList[i].ServiceTolMemory
		x2 := ths.AppList[j].ServiceTolMemory
		if x1 == x2 {
			return ths.AppList[i].SidecarCount < ths.AppList[j].SidecarCount
			// return ths.AppList[i].AppName < ths.AppList[j].AppName
		}
		return x1 < x2
	})

	// ths.Print()

	result := make(map[string][]string)

	sz := len(ths.Pilots)
	appNum := len(ths.AppList)

	if appNum%2 == 0 {
		for i := 0; i < appNum/2; i++ {
			pilot := ths.Pilots[i%sz]
			app1, app2 := ths.AppList[i], ths.AppList[appNum-i-1]
			result[pilot] = append(result[pilot], app1.AppName)
			result[pilot] = append(result[pilot], app2.AppName)
		}
	} else {
		for i := 0; i < appNum/2; i++ {
			pilot := ths.Pilots[i%sz]
			app1, app2 := ths.AppList[i+1], ths.AppList[appNum-i-1]
			result[pilot] = append(result[pilot], app1.AppName)
			result[pilot] = append(result[pilot], app2.AppName)
		}
		minPilot := ths.Pilots[0]
		result[minPilot] = append(result[minPilot], ths.AppList[0].AppName)
	}

	b, err := json.Marshal(result)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	return string(b), nil
}
