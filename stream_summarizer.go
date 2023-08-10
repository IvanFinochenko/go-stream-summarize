package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {

	N := 300
	M := 5
	K := 1000

	tickerGenerator := time.NewTicker(time.Duration(N) * time.Millisecond)
	tickerPublisher := time.NewTicker(time.Duration(K) * time.Millisecond)

	batchChannel := make(chan []int)
	summarizeChannel := make(chan int)

	go generator(tickerGenerator.C, batchChannel)

	for i := 0; i < M; i++ {
		go worker(batchChannel, summarizeChannel, i)
	}

	summarize(summarizeChannel, tickerPublisher.C)
}

func generator(tickerGenerate <-chan time.Time, batchChannel chan<- []int) {
	for {
		t := <-tickerGenerate
		fmt.Println("Tick at", t)
		arr := generateRandomSlice(10)
		batchChannel <- arr
	}
}

func generateRandomSlice(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(10)
	}
	return arr
}

func worker(ch <-chan []int, sum chan<- int, i int) {
	for {
		arr := <-ch
		fmt.Printf("Array %d in worker %d\n", arr, i)
		sort.Ints(arr)
		max3 := arr[7:]
		for _, x := range max3 {
			sum <- x
		}
	}
}

func summarize(sumChannel <-chan int, tickerPublisher <-chan time.Time) {
	sum := 0
	for {
		select {
		case x := <-sumChannel:
			sum = sum + x
		case <-tickerPublisher:
			fmt.Printf("Accumulator = %d\n", sum)
		}
	}
}
