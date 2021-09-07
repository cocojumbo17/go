package main

import (
	"fmt"
	"runtime"
	"strings"
)

const routinsNum = 5
const interations = 7

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat(" ", in), "█",
		strings.Repeat(" ", routinsNum-in), "th",
		in, "iter", j, strings.Repeat("█", j))
}

func doWork(in int) {
	for i := 0; i < interations; i++ {
		fmt.Printf(formatWork(in, i))
		if i%2 == 1 {
			runtime.Gosched()
		}
	}
}
func one() {
	for i := 0; i < routinsNum; i++ {
		go doWork(i)
	}
	fmt.Scanln()
}
func two() {
	ch1 := make(chan int, 1)
	go func(in chan int) {
		fmt.Println("GO: before read from channel: ")
		val := <-in
		fmt.Println("GO: read from channel: ", val)
		fmt.Println("GO: after read from channel")
	}(ch1)
	fmt.Println("MAIN: before write to channel")
	ch1 <- 42
	fmt.Println("MAIN: after write to channel")
	ch1 <- 43
	fmt.Scanln()
}

func three() {
	ch1 := make(chan int, 0)
	go func(out chan<- int) {
		fmt.Println("GO: start generator")
		for i := 0; i <= 10; i++ {
			fmt.Println("GO:before")
			out <- i
			fmt.Println("GO:after ", i)
		}
		close(out)
		fmt.Println("GO: finish generator")
	}(ch1)
	fmt.Println("MAIN: start iteration")
	for k := range ch1 {
		fmt.Println("MAIN: gen ", k)
	}
	fmt.Println("MAIN: finish iteration")
}
func four() {
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	ch2 <- 3
LOOP:
	for {
		select {
		case val := <-ch1:
			fmt.Println("ch1 val ", val)
		case val := <-ch2:
			fmt.Println("ch2 val ", val)
		default:
			fmt.Println("default")
			break LOOP
		}
	}
}

func five() {
	cancelCh := make(chan struct{})
	dataCh := make(chan int)
	go func(cCh chan struct{}, dCh chan int) {
		fmt.Println("GO: begin")
		defer fmt.Println("GO: end")
		val := 0
		for {
			select {
			case <-cCh:
				fmt.Println("GO: cancel")
				return
			case dCh <- val:
				fmt.Println("GO: val is written: ", val)
				val++
			}
		}
	}(cancelCh, dataCh)
	fmt.Println("MAIN: begin")
	for v := range dataCh {
		fmt.Println("MAIN: val is read: ", v)
		if v > 3 {
			fmt.Println("MAIN: before cancel")
			cancelCh <- struct{}{}
			fmt.Println("MAIN: after cancel")
			break
		}
	}
	fmt.Println("MAIN: end")
	fmt.Scanln()
}

func main() {
	//one()
	//two()
	//three()
	//four()
	five()
}
