package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	

	fmt.Println("--- Способ 1: sync.Mutex ---")
	var counterMutex int
	var mu sync.Mutex
	var wg1 sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mu.Lock()
			counterMutex++
			mu.Unlock()
		}()
	}
	wg1.Wait()
	fmt.Println("Результат с Mutex:", counterMutex)

	
	fmt.Println("\n--- Способ 2: sync/atomic ---")
	var counterAtomic int64
	var wg2 sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			atomic.AddInt64(&counterAtomic, 1)
		}()
	}
	wg2.Wait()
	fmt.Println("Результат с Atomic:", counterAtomic)
}