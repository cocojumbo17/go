package main

import (
	"bytes"
	"sort"
	"strconv"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	prev_out := make(chan interface{}, 20)
	wg := sync.WaitGroup{}
	for _, j := range jobs {
		var ch_input chan interface{} = prev_out
		var ch_output = make(chan interface{}, 20)
		wg.Add(1)
		go func(jo job, in, out chan interface{}, w *sync.WaitGroup) {
			defer w.Done()
			jo(in, out)
			close(out)
		}(j, ch_input, ch_output, &wg)
		prev_out = ch_output
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	for data := range in {
		str := strconv.FormatInt(int64(data.(int)), 10)
		var first, second string
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			first = DataSignerCrc32(str)
		}()
		go func() {
			defer wg.Done()
			second = DataSignerCrc32(DataSignerMd5(str))
		}()
		wg.Wait()
		out <- first + "~" + second
	}
}

func MultiHash(in, out chan interface{}) {
	for data := range in {
		var m map[int]string = make(map[int]string)
		wgr := sync.WaitGroup{}
		mux := sync.Mutex{}
		wgr.Add(6)
		for i := 0; i < 6; i++ {
			go func(index int, data string, res map[int]string, wg *sync.WaitGroup, mx *sync.Mutex) {
				defer wg.Done()
				str := strconv.FormatInt(int64(index), 10)
				s := DataSignerCrc32(str + data)
				mx.Lock()
				res[index] = s
				mx.Unlock()
			}(i, data.(string), m, &wgr, &mux)
		}
		wgr.Wait()
		s := ""
		for i := 0; i < 6; i++ {
			s += m[i]
		}
		out <- s
	}
}

func CombineResults(in, out chan interface{}) {
	var ss []string
	for s := range in {
		ss = append(ss, s.(string))
	}
	sort.Strings(ss)
	var res bytes.Buffer
	for _, str := range ss {
		res.WriteString(str)
		res.WriteString("_")
	}
	res.Truncate(res.Len() - 1)
	out <- res.String()
}
