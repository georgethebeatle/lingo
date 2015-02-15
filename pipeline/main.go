package main

import (
	"fmt"
	"sync"
)

func numbers(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
		}
	}()

	return out
}

func sq(in <-chan int, colour string) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			fmt.Printf("%s handling '%d'\n", colour, num)
			out <- num * num
		}
	}()

	return out
}

func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	drain := func(channel <-chan int) {
		defer wg.Done()
		for num := range channel {
			out <- num
		}
	}

	// drain all channels
	wg.Add(len(channels))
	for _, channel := range channels {
		go drain(channel)
	}

	// wait for all routines to finish and close out
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func main() {
	num_chan := numbers(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	red := sq(num_chan, "red")
	green := sq(num_chan, "green")
	blue := sq(num_chan, "blue")

	results := merge(red, green, blue)

	for num := range results {
		fmt.Printf("=> %d\n", num)
	}
}
