package main

import "fmt"

type Broadcaster interface {
	Register() <-chan int
	Send() chan<- int
	Close()
}

type broadcaster struct {
	sender    chan int
	receivers []chan int
}

func NewBroadCaster() Broadcaster {
	receivers := make([]chan int, 0, 10)
	broadcaster := &broadcaster{
		sender:    make(chan int),
		receivers: receivers,
	}
	go broadcast(broadcaster)
	return broadcaster
}

func (b *broadcaster) Register() <-chan int {
	receiver := make(chan int)
	b.receivers = append(b.receivers, receiver)
	return receiver
}

func (b *broadcaster) Send() chan<- int {
	return b.sender
}

func (b *broadcaster) Close() {
	close(b.sender)
}

func broadcast(b *broadcaster) {
	for msg := range b.sender {
		for _, receiver := range b.receivers {
			println("Re-Sending ", msg)
			receiver <- msg
		}
	}
	println("Sender is closed")
	for _, receiver := range b.receivers {
		println("Close receiver")
		close(receiver)
	}
}

func receive(one, two <-chan int) {
	for {
		if one == nil && two == nil {
			return
		}
		select {
		case msg, ok := <-one:
			if ok {
				fmt.Println("Channel one received: ", msg)
			} else {
				fmt.Println("Channel one is closed")
				one = nil
			}
		case msg, ok := <-two:
			if ok {
				fmt.Println("Channel two received: ", msg)
			} else {
				fmt.Println("Channel one is closed")
				two = nil
			}
		}
	}
}

func main() {
	b := NewBroadCaster()
	one := b.Register()
	two := b.Register()

	go func() {
		defer b.Close()

		println("Sending 1")
		b.Send() <- 1
		println("Sending 2")
		b.Send() <- 2
		println("Sending 3")
		b.Send() <- 3
	}()

	receive(one, two)

}
