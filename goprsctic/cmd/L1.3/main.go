package main

import (
	"fmt"
	"sync"
)

func main() {
	var n int
	fmt.Scan(&n)

	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}

	for i := 1; i <= 20; i++ {
		ch <- i
	}

	close(ch)

	wg.Wait()
}

func worker(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range ch {
		res := i * i
		fmt.Println("Worker", id, res)
	}

}
