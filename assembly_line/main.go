package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Employee struct {
	id int
}

type Item1 struct{}

type Item2 struct{}

type Item3 struct{}

type Item interface {
	// Process 這是一個耗時操作
	Process()
}

func (Item1) Process() { time.Sleep(120 * time.Millisecond) }
func (Item2) Process() { time.Sleep(220 * time.Millisecond) }
func (Item3) Process() { time.Sleep(350 * time.Millisecond) }

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	const (
		workers    = 5
		perTypeN   = 10
		totalTasks = perTypeN * 3
	)

	// 準備物品：各 10 件
	items := make([]Item, 0, totalTasks)
	for i := 0; i < perTypeN; i++ {
		items = append(items, Item1{}, Item2{}, Item3{})
	}

	// 打亂順序
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

	jobs := make(chan Item)

	// 每位員工處理數量統計
	counts := make([]int, workers)

	// 等待所有任務完成
	var doneWg sync.WaitGroup
	doneWg.Add(totalTasks)

	// 啟動 5 個員工
	var workerWg sync.WaitGroup
	workerWg.Add(workers)
	for wid := 0; wid < workers; wid++ {
		go func(id int) {
			defer workerWg.Done()
			for it := range jobs {
				// 開始處理
				log.Printf("[Worker %d] 開始處理: %T", id, it)

				start := time.Now()
				it.Process()
				elapsed := time.Since(start)

				// 完成處理
				log.Printf("[Worker %d] 完成處理: %T，用時=%v", id, it, elapsed)

				counts[id]++
				doneWg.Done()
			}
		}(wid)
	}

	// 發派任務並計時
	totalStart := time.Now()
	go func() {
		for _, it := range items {
			jobs <- it
		}
		close(jobs)
	}()

	// 等所有任務完成
	doneWg.Wait()
	totalElapsed := time.Since(totalStart)

	// 等所有員工結束
	workerWg.Wait()

	// 輸出統計
	fmt.Println("========== 統計 ==========")
	fmt.Printf("總處理時間: %v\n", totalElapsed)
	for i, c := range counts {
		fmt.Printf("員工 %d 處理件數: %d\n", i, c)
	}
}
