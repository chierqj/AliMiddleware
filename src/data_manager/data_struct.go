package data_manager

import "flag"

const (
	MB_EACH_NODE = 0.01
)

var (
	filePath = flag.String("data_path", "/root/input/data.json", "input data path.")
)

type LoadData struct {
	Apps         map[string]int            `json:"apps"`
	Dependencies map[string]map[string]int `json:"dependencies"`
}

type Service struct {
	ID    int
	Name  string
	Count int
}
type App struct {
	AppName          string    // 应用名称
	SidecarCount     int       // Sidecar个数
	ServiceTolMemory int       // 服务总内存
	ServiceList      []Service // 服务依赖列表
}

type Scenario struct {
	Params      LoadData            // 读入的所有数据
	AppList     []*App              // 所有的app
	Pilots      []string            // pilots
	MapApp      map[string]*App     // MapApp
	PilotsSize  int                 // pilots size
	AppListSize int                 // app_list size
	TotalMemory float64             // 服务总内存
	Result      map[string][]string // 返回结果
}

var scenario *Scenario

func GetInstance() *Scenario {
	if scenario == nil {
		scenario = &Scenario{}
		return scenario
	}
	return scenario
}
