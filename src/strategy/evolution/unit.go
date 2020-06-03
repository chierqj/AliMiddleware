package evolution

import (
	"fmt"
	"math"

	data "src/data_manager"

	"github.com/golang/glog"
)

const (
	MB_EACH_NODE = 0.01
)

// 编码结点
type DecodeNode struct {
	apps       []*data.App // 包含的app
	loadMemory float64     // 加载内存
	sidecar    float64     // 连接sidecar数
	serviceNum int         // 连接服务数
	appNum     int         // app数目
}
type Unit struct {
	decode        []*DecodeNode // 编码
	fitness       float64       // 适应度值
	loadMemory    float64       // 加载内存
	serviceMemory float64       // 服务总内存
	stdMemory     float64       // 加载内存标准差
	stdConnection float64       // 连接数标准差
	splitIndex    []int         // 分割的位置
	totalApps     []*data.App   // 完整的applist
}

func (ths *DecodeNode) MsgStr() string {
	return fmt.Sprintf("[连接: %.0f, 加载内存: %.3f, 服务数目: %d, apps: %d]\n",
		ths.sidecar,
		ths.loadMemory,
		ths.serviceNum,
		ths.appNum,
	)
}

func (ths *Unit) Print() {
	glog.Infoln("-----------------------------------------------------------------")
	glog.Infof("@ Fitness: %.3f\n", ths.fitness)
	glog.Infof("@ 加载内存/服务内存: %.3f\n", ths.loadMemory/ths.serviceMemory)
	glog.Infof("@ 加载内存: %.3f\n", ths.loadMemory)
	glog.Infof("@ 内存标准差: %.3f\n", ths.stdMemory)
	glog.Infof("@ 连接标准差: %.3f\n", ths.stdConnection)
	for i, node := range ths.decode {
		glog.Infof("* %d: %s", i, node.MsgStr())
	}
}

func (ths *Unit) CreateDecodeBySplitIndex() {
	pre := 0
	for _, idx := range ths.splitIndex {
		decodeNode := &DecodeNode{}
		for i := pre; i < idx; i++ {
			decodeNode.apps = append(decodeNode.apps, ths.totalApps[i])
		}
		pre = idx
		ths.decode = append(ths.decode, decodeNode)
	}
}

func (ths *DecodeNode) CalMemoryAndSidecar() {
	// vis := mapset.NewSet()
	// var vis [1000000]bool
	vis := make([]bool, len(ths.apps), 100)
	ths.loadMemory = 0
	ths.sidecar = 0
	ths.serviceNum = 0
	for _, app := range ths.apps {
		cnt := 0
		for _, svc := range app.ServiceList {
			if !vis[svc.ID] {
				vis[svc.ID] = true
				ths.loadMemory += float64(svc.Count) * MB_EACH_NODE
				cnt++
			}
		}
		ths.serviceNum += cnt
		ths.sidecar += float64(app.SidecarCount)
	}
	ths.appNum = len(ths.apps)
}

func (ths *Unit) CalculateFitness() {
	avgMem, avgConn := 0.0, 0.0
	for _, node := range ths.decode {
		node.CalMemoryAndSidecar()
		avgMem += node.loadMemory
		avgConn += node.sidecar
	}
	length := float64(len(ths.decode))
	ths.loadMemory = avgMem
	avgMem /= length
	avgConn /= length
	stdMem, stdConn := 0.0, 0.0
	for _, node := range ths.decode {
		stdMem += (node.loadMemory - avgMem) * (node.loadMemory - avgMem)
		stdConn += (node.sidecar - avgConn) * (node.sidecar - avgConn)
	}
	stdMem, stdConn = math.Sqrt(stdMem), math.Sqrt(stdConn)
	ths.stdMemory, ths.stdConnection = stdMem, stdConn
	ths.fitness = (ths.loadMemory / ths.serviceMemory) * (ths.stdMemory + ths.stdConnection)
}
