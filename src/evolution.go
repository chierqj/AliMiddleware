package main

import (
	"math/rand"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/golang/glog"
)

const (
	EPOCH_TIME      = 10 // 迭代次数
	POPULATION_SIZE = 6  // 种群大小
)

// 进化算法
type Evolution struct {
	Population [POPULATION_SIZE]*Unit // 种群
}

func (ths *Evolution) Print() {
	glog.Infoln("============== Evolution ==============")
	glog.Infof("@ 种群大小: %d\n", POPULATION_SIZE)
	glog.Infof("@ 迭代次数: %d\n", EPOCH_TIME)
	for i, it := range ths.Population {
		if i >= 6 {
			break
		}
		it.Print()
	}
}

// 入口
func (ths *Evolution) Run(player *Player) {
	ths.InitializePopulation(player)
	ths.Print()
	for epoch := 0; epoch < EPOCH_TIME; epoch++ {

	}
}

func (ths *Evolution) InitializePopulation(player *Player) {
	rand.Seed(time.Now().UnixNano())
	// 生成POPULATION_SIZE个个体
	for p := 0; p < POPULATION_SIZE; p++ {
		// 生成pilotSize - 1个不重复的随机数用来初始化分配apps
		var index []int
		cnt, vis := 0, mapset.NewSet()
		for cnt < player.PilotsSize-1 {
			x := rand.Intn(player.AppListSize)
			// x := (player.AppListSize / player.PilotsSize) * (cnt + 1)
			if !vis.Contains(x) {
				index = append(index, x)
				cnt++
				vis.Add(x)
			}
		}
		// 排序，每个数字作为分割点
		sort.Ints(index)
		// 打乱applist，保证每个个体尽量随机
		rand.Shuffle(player.AppListSize, func(i, j int) {
			player.AppList[i], player.AppList[j] = player.AppList[j], player.AppList[i]
		})
		index = append(index, player.AppListSize)

		// 初始化每个个体unit
		unit := &Unit{serviceMemory: player.TotalMemory}
		pre := 0
		for _, idx := range index {
			decodeNode := &DecodeNode{}
			for i := pre; i < idx; i++ {
				decodeNode.apps = append(decodeNode.apps, player.AppList[i])
			}
			pre = idx
			unit.decode = append(unit.decode, decodeNode)
		}
		unit.CalculateFitness()
		ths.Population[p] = unit
	}
}

func (ths *Evolution) Cross() {

}

func (ths *Evolution) Variation() {
}

func (ths *Evolution) Select() {

}
