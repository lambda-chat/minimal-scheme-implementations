package main

import (
	"fmt"
	"sync"
)

func generateInt(next func(int) int, x int) <-chan int {
	ch := make(chan int, 0) // block
	go func() {
		defer close(ch)
		for {
			ch <- x
			x = next(x)
		}
	}()
	return ch
}

func GeneratorMain1() {
	ch := generateInt(func(x int) int { return x + 1 }, 0)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}

func clientProcess(yield chan<- int, send <-chan int) string {
	yield <- 1
	_, _ = fmt.Printf("[client] recieved %d\n", <-send)
	yield <- 2
	_, _ = fmt.Printf("[client] recieved %d\n", <-send)
	yield <- 3
	_, _ = fmt.Printf("[client] recieved %d\n", <-send)
	return "client process done."
}

func serverProcess(yield chan<- int, send <-chan int) string {
	_, _ = fmt.Printf("[server] recieved %d\n", <-send)
	yield <- 11
	_, _ = fmt.Printf("[server] recieved %d\n", <-send)
	yield <- 22
	_, _ = fmt.Printf("[server] recieved %d\n", <-send)
	yield <- 33
	return "server process done."
}

func GeneratorMain2() {
	client := make(chan int, 0)
	server := make(chan int, 0)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		end := clientProcess(server, client)
		_, _ = fmt.Println(end)
		wg.Done()
	}()
	go func() {
		end := serverProcess(client, server)
		_, _ = fmt.Println(end)
		wg.Done()
	}()
	wg.Wait()
}
