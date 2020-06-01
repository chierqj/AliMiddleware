package main

import (
	"encoding/json"
	"sort"

	"github.com/golang/glog"
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
	ths.LoadData()
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
