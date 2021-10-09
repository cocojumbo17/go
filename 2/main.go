package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const routinsNum = 10
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

func longSQLQuery() chan bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- true
	}()
	return ch
}

func fnTimer() {
	timer := time.NewTimer(3 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("timer.C timeout happened")
	case <-time.After(10 * time.Second):
		fmt.Println("timer.After timeout happened")
	case res := <-longSQLQuery():
		fmt.Println("longSQLQuery finished with result: ", res)
		if !timer.Stop() {
			fmt.Println("Force stop of timer")
			<-timer.C
		}
	}
}
func fnTicker() {
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			ticker.Stop()
			break
		}
	}
	fmt.Println("total", i)

	c := time.Tick(time.Second)
	i = 0
	for tickTime := range c {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			break
		}
	}
}

func sayHello() {
	fmt.Println("Hello AfterFunck")
}

func afterF() {
	t := time.AfterFunc(time.Second*5, sayHello)
	fmt.Scanln()
	t.Stop()
	fmt.Scanln()
}

func worker(ctx context.Context, number int, out chan<- int) {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
	fmt.Println("Go:", number, "sleep", waitTime)
	select {
	case <-ctx.Done():
		fmt.Println("Go: worker", number, "canceled")
		return
	case <-time.After(waitTime):
		fmt.Println("Go: worker", number, "done")
		out <- number
	}
}

func fnCancel() {
	ctx, finish := context.WithCancel(context.Background())
	out := make(chan int, 1)
	for i := 0; i < 10; i++ {
		go worker(ctx, i, out)
	}
	foundBy := <-out
	fmt.Println("Main: result found by", foundBy)
	finish()

	time.Sleep(time.Second)
}

func fnTimeOut() {
	workTime := 50 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), workTime)
	out := make(chan int, 1)
	for i := 0; i < 10; i++ {
		go worker(ctx, i, out)
	}
	totalFound := 0
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case foundBy := <-out:
			totalFound++
			fmt.Println("Main: result found by", foundBy)
		}
	}
	fmt.Println("Main: total found", totalFound)
	time.Sleep(time.Second)
}

func goPull(num int, in chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range in {
		fmt.Println("Go: ", num, ": ", data)
		runtime.Gosched()
	}
	fmt.Println("Go: finish of ", num, " worker")
}

func pullWorker() {
	fmt.Println("Main: start")
	wg := sync.WaitGroup{}
	const goNum int = 5
	in := make(chan string, 3)
	for i := 0; i < goNum; i++ {
		wg.Add(1)
		go goPull(i, in, &wg)
	}
	mounths := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for _, name := range mounths {
		in <- name
	}
	fmt.Println("Main: before closing")
	close(in)
	fmt.Println("Main: after closing")
	wg.Wait()
	fmt.Println("Main: finish")
}

func goQuota(num int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{}
	defer wg.Done()
	for i := 0; i < interations; i++ {
		fmt.Println(formatWork(num, i))
		// if i%4 == 0 {
		// 	<-quotaCh
		// 	quotaCh <- struct{}{}
		// }
		runtime.Gosched()
	}
	fmt.Println("Go: finish of ", num, " worker")
	<-quotaCh
}

func pullQuota() {
	fmt.Println("Main: start")
	wg := sync.WaitGroup{}
	const quotaLimit = 4
	quotaCh := make(chan struct{}, quotaLimit)
	for i := 0; i < routinsNum; i++ {
		wg.Add(1)
		go goQuota(i, &wg, quotaCh)
	}
	wg.Wait()
	fmt.Println("Main: finish")
}

func gaisenBug() {
	var counters = map[int]int{}
	mx := sync.Mutex{}
	for i := 0; i < 5; i++ {
		go func(m map[int]int, th int, mux *sync.Mutex) {
			for j := 0; j < 5; j++ {
				mux.Lock()
				m[th*5+j]++
				mux.Unlock()
			}
			runtime.Gosched()
		}(counters, i, &mx)
	}
	fmt.Scanln()
	mx.Lock()
	fmt.Println("result:", counters)
	mx.Unlock()
}

func Inc() {
	var count int32 = 0
	for i := 0; i < 1000; i++ {
		go func() {
			//count++
			atomic.AddInt32(&count, 1)
		}()
	}
	time.Sleep(time.Microsecond * 10)
	fmt.Println(count)

}

func main() {
	//one()
	//two()
	//three()
	//four()
	//five()
	//fnTimer()
	//fnTicker()
	//afterF()
	//fnCancel()
	//fnTimeOut()
	//pullWorker()
	//pullQuota()
	//gaisenBug()
	Inc()
}
