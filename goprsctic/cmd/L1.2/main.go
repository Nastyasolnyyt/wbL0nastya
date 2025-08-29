package main

import "fmt"

func main() {
	m := [5]int{2, 4, 6, 8, 10}
	ch := make(chan int, len(m))
	for i := 0; i < 5; i++ {
		go func(n int) {
			ch <- n * n
		}(m[i])
	}

	for i := 0; i < len(m); i++ {
		fmt.Println(<-ch)
	}
}
