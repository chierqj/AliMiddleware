package evolution

// package main

// import (
// 	"math/rand"
// 	"sort"
// 	"time"

// 	mapset "github.com/deckarep/golang-set"
// 	"github.com/golang/glog"
// )

// const (
// 	EPOCH_TIME      = 10  // 迭代次数
// 	POPULATION_SIZE = 500 // 种群大小
// 	CROSSOVER_PROB  = 0.8 // 交叉概率
// 	MUTATION_PROB   = 0.2 // 变异概率
// )

// // 进化算法
// type Evolution struct {
// 	Population []*Unit // 种群
// 	Progeny    []*Unit // 每一代迭代中的子代种群
// 	BestUnit   *Unit   // 每一代最优个体
// }

// func (ths *Evolution) Print() {
// 	glog.Infoln("============== Evolution ==============")
// 	glog.Infof("@ 种群大小: %d\n", len(ths.Population))
// 	glog.Infof("@ 迭代次数: %d\n", EPOCH_TIME)
// 	for i, it := range ths.Population {
// 		if i >= 6 {
// 			break
// 		}
// 		it.Print()
// 	}
// }

// // 入口
// func (ths *Evolution) Run() {
// 	ths.InitializePopulation()
// 	ths.Print()
// 	for epoch := 0; epoch < EPOCH_TIME; epoch++ {
// 		ths.CrossOver()
// 		ths.Mutation()
// 		ths.Selection()
// 		ths.BestUnit.Print()
// 	}
// }

// func (ths *Evolution) InitializePopulation() {
// 	rand.Seed(time.Now().UnixNano())
// 	// 生成POPULATION_SIZE个个体
// 	for p := 0; p < POPULATION_SIZE; p++ {
// 		// 生成pilotSize - 1个不重复的随机数用来初始化分配apps
// 		var index []int
// 		cnt, vis := 0, mapset.NewSet()
// 		for cnt < ths.PilotsSize-1 {
// 			x := rand.Intn(player.AppListSize)
// 			// x := (player.AppListSize / player.PilotsSize) * (cnt + 1)
// 			if !vis.Contains(x) {
// 				index = append(index, x)
// 				cnt++
// 				vis.Add(x)
// 			}
// 		}
// 		// 排序，每个数字作为分割点
// 		sort.Ints(index)
// 		// // 打乱applist，保证每个个体尽量随机
// 		rand.Shuffle(player.AppListSize, func(i, j int) {
// 			player.AppList[i], player.AppList[j] = player.AppList[j], player.AppList[i]
// 		})
// 		index = append(index, player.AppListSize)

// 		// // 初始化每个个体unit
// 		unit := &Unit{serviceMemory: player.TotalMemory, splitIndex: index, totalApps: player.AppList}
// 		unit.CreateDecodeBySplitIndex()
// 		unit.CalculateFitness()
// 		ths.Population = append(ths.Population, unit)
// 		if p%10 == 0 {
// 			glog.Infoln(p)
// 		}
// 	}
// }

// func (ths *Evolution) doCrossOver(unit1, unit2 *Unit) (*Unit, *Unit) {
// 	// p1 := rand.Intn(ths.UnitDecodeLen)
// 	// p2 := rand.Intn(ths.UnitDecodeLen)
// 	// if p1 == p2 {
// 	// 	return nil, nil
// 	// }
// 	// if p1 > p2 {
// 	// 	p1, p2 = p2, p1
// 	// }
// 	// return unit1, unit2
// }
// func (ths *Evolution) CrossOver() {
// 	rand.Seed(time.Now().UnixNano())
// 	ths.Progeny = ths.Progeny[0:0]
// 	for i, unit := range ths.Population {
// 		rd := rand.Float32()
// 		if rd > CROSSOVER_PROB {
// 			continue
// 		}
// 		j := rand.Intn(POPULATION_SIZE)
// 		if j == i {
// 			continue
// 		}
// 		newUnit1, newUnit2 := ths.doCrossOver(unit, ths.Population[j])
// 		if newUnit1 == nil || newUnit2 == nil {
// 			continue
// 		}
// 		ths.Progeny = append(ths.Progeny, newUnit1)
// 		ths.Progeny = append(ths.Progeny, newUnit2)
// 	}
// }

// func (ths *Evolution) doMutation(unit *Unit) {

// }
// func (ths *Evolution) Mutation() {
// 	rand.Seed(time.Now().UnixNano())
// 	for _, unit := range ths.Progeny {
// 		rd := rand.Float32()
// 		if rd >= MUTATION_PROB {
// 			continue
// 		}
// 		ths.doMutation(unit)
// 	}
// }

// func (ths *Evolution) Selection() {
// 	ths.Progeny = append(ths.Progeny, ths.Population...)
// 	ths.BestUnit = ths.Population[0]
// }
