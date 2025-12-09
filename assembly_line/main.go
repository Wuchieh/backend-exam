package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Employee struct {
	ID           int
	ProcessCount int
	mu           sync.Mutex
}

func (e *Employee) Work(items <-chan Item, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range items {
		fmt.Printf("員工 %d 開始處理 %s\n", e.ID, item.Name())
		item.Process()
		fmt.Printf("員工 %d 完成處理 %s\n", e.ID, item.Name())
		e.mu.Lock()
		e.ProcessCount++
		e.mu.Unlock()
	}
}

type Item1 struct{}

func (i Item1) Process() {
	time.Sleep(100 * time.Millisecond)
}

func (i Item1) Name() string {
	return "物品1"
}

type Item2 struct{}

func (i Item2) Process() {
	time.Sleep(200 * time.Millisecond)
}

func (i Item2) Name() string {
	return "物品2"
}

type Item3 struct{}

func (i Item3) Process() {
	time.Sleep(150 * time.Millisecond)
}

func (i Item3) Name() string {
	return "物品3"
}

type Item interface {
	// Process 這是一個耗時操作
	Process()
	Name() string
}

func main() {
	// 建立三種物品各十件
	items := make([]Item, 0, 30)
	for i := 0; i < 10; i++ {
		items = append(items, Item1{})
	}
	for i := 0; i < 10; i++ {
		items = append(items, Item2{})
	}
	for i := 0; i < 10; i++ {
		items = append(items, Item3{})
	}

	// 隨機打亂物品順序
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	// 建立工作佇列
	itemChan := make(chan Item, len(items))
	for _, item := range items {
		itemChan <- item
	}
	close(itemChan)

	// 建立五個員工
	employees := make([]*Employee, 5)
	for i := 0; i < 5; i++ {
		employees[i] = &Employee{ID: i + 1}
	}

	// 記錄開始時間
	startTime := time.Now()
	fmt.Printf("流水線開始時間: %s\n", startTime.Format("15:04:05.000"))

	// 啟動員工工作
	var wg sync.WaitGroup
	for _, emp := range employees {
		wg.Add(1)
		go emp.Work(itemChan, &wg)
	}

	// 等待所有工作完成
	wg.Wait()

	// 記錄結束時間
	endTime := time.Now()
	fmt.Printf("流水線結束時間: %s\n", endTime.Format("15:04:05.000"))

	// 統計總處理時間
	totalTime := endTime.Sub(startTime)
	fmt.Printf("\n總處理時間: %v\n", totalTime)

	// 統計每個員工處理的物品數量
	fmt.Println("\n員工處理統計:")
	for _, emp := range employees {
		fmt.Printf("員工 %d 處理了 %d 件物品\n", emp.ID, emp.ProcessCount)
	}
}
