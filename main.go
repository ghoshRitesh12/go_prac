package main

import (
	"fmt"
	"time"
)

func Main2() {
	c1 := make(chan string)
	c2 := make(chan string)
	
	go func() { 
		for {
			c1 <- "every 500ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() { 
		for {
			c2 <- "every 2s"
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		select {
			case msg1 := <- c1:
				fmt.Println(msg1)
			case msg2 := <- c2:
				fmt.Println(msg2)
		}
	}
}

func worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- fib(n)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n - 1) + fib(n - 2)
}

func main() {
	bufferSize := 50
	jobs := make(chan int, bufferSize)
	results := make(chan int, bufferSize)

	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)

	for i := 0; i < bufferSize; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < bufferSize; i++ {
		fmt.Println(<-results)
	}
}
