package random_walk

import (
	data "src/data_manager"
)

type RandomWalk struct {
}

func (ths *RandomWalk) ExecuteP1() map[string][]string {
	return ths.arrange()
}

func (ths *RandomWalk) ExecuteP2() map[string][]string {
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
func (ths *RandomWalk) arrange() map[string][]string {
	sce := data.GetInstance()
	result := make(map[string][]string)

	sz := len(sce.Pilots)
	for i, app := range sce.AppList {
		pilot := sce.Pilots[i%sz]
		result[pilot] = append(result[pilot], app.AppName)
	}
	return result
}
