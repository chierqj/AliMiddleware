package data_manager

import (
	"encoding/json"
	"fmt"
	"os"

	mapset "github.com/deckarep/golang-set"
)

/*
 * @ title:				LogMsg
 * @ description: 构造App结构体各参数输出字符串
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			string: logstr
 */
func (ths *App) LogMsg() string {
	logStr := fmt.Sprintf("app: %s, sidecar: %d, svcs: %d, svclen: %d",
		ths.AppName,
		ths.SidecarCount,
		ths.ServiceTolMemory,
		len(ths.ServiceList),
	)
	return logStr
}

/*
 * @ title:				Print
 * @ description: Scenario类打印应用
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			null
 */
func (ths *Scenario) Print() {
	for idx, app := range ths.AppList {
		fmt.Printf("@ %d: [%s]\n", idx, app.LogMsg())
	}
}

/*
 * @ title:				LoadData
 * @ description: 访问 /api/p1_start 接口时，从文件data.json加载app以及相应依赖数据
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			null
 */
func (ths *Scenario) LoadData() {
	f, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
	}
	if err = json.NewDecoder(f).Decode(&ths.Params); err != nil {
		fmt.Println(err)
	}
	ths.CreateAppList()
}

/*
 * @ title:				CreateAppList
 * @ description: 根据Scenario成员变量Param，将map转换到AppList
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			null
 */
func (ths *Scenario) CreateAppList() {
	ths.AppList = ths.AppList[0:0]
	ths.TotalMemory = 0
	vis := mapset.NewSet()
	idx := 0
	for appName, sidecar := range ths.Params.Apps {
		var svcList []Service
		sum := 0
		for svc, count := range ths.Params.Dependencies[appName] {
			svcList = append(svcList, Service{ID: idx, Name: svc, Count: count})
			idx++
			sum += count
			if !vis.Contains(svc) {
				vis.Add(svc)
				ths.TotalMemory += float64(count) * MB_EACH_NODE
			}
		}
		app := &App{
			AppName:          appName,
			SidecarCount:     sidecar,
			ServiceTolMemory: sum,
			ServiceList:      svcList,
		}
		ths.AppList = append(ths.AppList, app)
	}
	ths.AppListSize = len(ths.AppList)
}

/*
 * @ title:				UpdateParams
 * @ description: 访问 /api/p2_start 时，更新Scenario.Param参数，同时更新AppList
 * @ author: 			FatChier
 * @ param:				param: LoadData {
 *	 									Apps map[string]int            `json:"apps"`
 *										Dependencies map[string]map[string]int `json:"dependencies"`
 *								}
 * @ return:			null
 */
func (ths *Scenario) UpdateParams(params LoadData) {
	for app, num := range params.Apps {
		ths.Params.Apps[app] = num
	}
	for app, svcs := range params.Dependencies {
		existed := ths.Params.Dependencies[app]
		if existed == nil {
			existed = map[string]int{}
			ths.Params.Dependencies[app] = existed
		}
		for svc, num := range svcs {
			existed[svc] = num
		}
	}
	ths.CreateAppList()
}
