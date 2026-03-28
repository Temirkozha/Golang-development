package main

import (
	"fmt"
	"sync"
)

func main() {
	
	fmt.Println("--- Task 1: RWMutex ---")
	mapMutex := make(map[string]int)
	var mu sync.RWMutex
	var wg1 sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg1.Add(1)
		go func(v int) {
			defer wg1.Done()
			mu.Lock()
			mapMutex["key"] = v
			mu.Unlock()
		}(i)
	}
	wg1.Wait()
	fmt.Println("Final value (Mutex):", mapMutex["key"])

	
	fmt.Println("\n--- Task 1: sync.Map ---")
	var mapSync sync.Map
	var wg2 sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func(v int) {
			defer wg2.Done()
			mapSync.Store("key", v)
		}(i)
	}
	wg2.Wait()
	val, _ := mapSync.Load("key")
	fmt.Println("Final value (sync.Map):", val)
}