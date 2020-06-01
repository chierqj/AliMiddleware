package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	filePath = flag.String("data_path", "/root/input/data.json", "input data path.")
)

type LoadData struct {
	Apps         map[string]int            `json:"apps"`
	Dependencies map[string]map[string]int `json:"dependencies"`
}

type Service struct {
	Name  string
	Count int
}
type App struct {
	AppName      string    // 应用名称
	SidecarCount int       // Sidecar个数
	ServiceList  []Service // 服务依赖列表
}

type Player struct {
	Params  LoadData        // 读入的所有数据
	AppList []*App          // 所有的app
	Pilots  []string        // pilots
	MapApp  map[string]*App // MapApp
}

/*
 * @ title:				LogMsg
 * @ description: 构造App结构体各参数输出字符串
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			string: logstr
 */
func (ths *App) LogMsg() string {
	logStr := fmt.Sprintf("AppName: %s, SidecarCount: %d, ServiceList: %d",
		ths.AppName,
		ths.SidecarCount,
		len(ths.ServiceList),
	)
	return logStr
}

/*
 * @ title:				Ready
 * @ description: 初始化Player参数，检查程序是否准备启动相应请求
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			bool: true准备好相应/api/ready接口, false未准备好
 */
func (ths *Player) Ready() bool {
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
 * @ title:				Print
 * @ description: Player类打印应用
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			null
 */
func (ths *Player) Print() {
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
func (ths *Player) LoadData() {
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
 * @ description: 根据Player成员变量Param，将map转换到AppList
 * @ author: 			FatChier
 * @ param:				null
 * @ return:			null
 */
func (ths *Player) CreateAppList() {
	ths.AppList = ths.AppList[0:0]
	for appName, sidecar := range ths.Params.Apps {
		var srvList []Service
		for srv, count := range ths.Params.Dependencies[appName] {
			srvList = append(srvList, Service{srv, count})
		}
		app := &App{
			AppName:      appName,
			SidecarCount: sidecar,
			ServiceList:  srvList,
		}
		ths.AppList = append(ths.AppList, app)
	}
}

/*
 * @ title:				UpdateParams
 * @ description: 访问 /api/p2_start 时，更新Player.Param参数，同时更新AppList
 * @ author: 			FatChier
 * @ param:				param: LoadData {
 *	 									Apps map[string]int            `json:"apps"`
 *										Dependencies map[string]map[string]int `json:"dependencies"`
 *								}
 * @ return:			null
 */
func (ths *Player) UpdateParams(params LoadData) {
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
