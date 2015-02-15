package main

import (
	"fmt"
	"os"
	"sync"
)

var statistics = make(map[string]int)

func numbers(done <-chan struct{}, nums ...int) <-chan int {
	statistics["total"] = len(nums)
	out := make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums {
			select {
			case out <- num:
			case <-done:
				return
			}
		}
	}()

	return out
}

func worker_sq(done <-chan struct{}, in <-chan int, colour string) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			select {
			case out <- num * num:
				statistics[colour] += 1
				fmt.Printf("%s handling '%d'\n", colour, num)
			case <-done:
				return
			}
		}
	}()

	return out
}

func merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	drain := func(channel <-chan int) {
		defer wg.Done()
		for num := range channel {
			select {
			case out <- num:
			case <-done:
				return
			}
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
	// channel for interruption of the pipeline
	done := make(chan struct{})

	num_chan := numbers(done, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	red := worker_sq(done, num_chan, "red")
	green := worker_sq(done, num_chan, "green")
	blue := worker_sq(done, num_chan, "blue")

	results := merge(done, red, green, blue)

	if len(os.Args) > 1 && os.Args[1] == "--interrupt" {
		// interrupt the pipeline execution at random time
		go func() {
			fmt.Println("Ahaaaaa!")
			close(done)
		}()
	}

	for num := range results {
		fmt.Printf("=> %d\n", num)
	}

	fmt.Println("---")
	fmt.Printf("Total numbers: %d\n", statistics["total"])
	fmt.Printf("Red: %d\n", statistics["red"])
	fmt.Printf("Green: %d\n", statistics["green"])
	fmt.Printf("Blue: %d\n", statistics["blue"])

}
