package data_manager

import "github.com/golang/glog"

/*
 * @ title:				Ready
 * @ description: 初始化Scenario参数，检查程序是否准备启动相应请求
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			bool: true准备好相应/api/ready接口, false未准备好
 */
func (ths *Scenario) Ready() bool {
	ths.Params = LoadData{
		Apps:         make(map[string]int),
		Dependencies: make(map[string]map[string]int),
	}
	ths.AppList = ths.AppList[0:0]
	ths.Pilots = ths.Pilots[0:0]
	ths.MapApp = make(map[string]*App)
	return true
}

/*
 * @ title:				P1
 * @ description: 响应 /api/p1_start 接口，返回pilots的分配方案
 * @ author: 			FatChier
 * @ param:				pilots ["p1","p2",...]string 需要分配的 pilots
 * @ return:			string: {"p1":["app1","app2",...],"p2":[]...} 分配方案对应的json字符串
 * @ 							error: error or nil
 */
func (ths *Scenario) HandleP1(pilots []string) {
	glog.Infof("P1 start: %v\n", pilots)
	ths.Pilots = pilots
	ths.PilotsSize = len(pilots)
	ths.LoadData()
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
func (ths *Scenario) HandleP2(params LoadData) {
	glog.Infof("P2 start: %v\n", ths.Pilots)
	ths.UpdateParams(params)
}
