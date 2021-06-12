package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
)

// сюда писать код
func ExecutePipeline(hashSignJobs ...job) {

	in := make(chan interface{})
	for _, worker := range hashSignJobs {
		in = func(a job, in chan interface{}) chan interface{} {
			out := make(chan interface{})
			go func(a job, in chan interface{}, out chan interface{}) {
				a(in, out)
				close(out)
			}(a, in, out)
			return out
		}(worker, in)
	}
	<-in
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for dataRaw := range in {
		data, ok := dataRaw.(int)
		if !ok {
			fmt.Println("cant convert to int")
		} else {
			wg.Add(1)
			go func(data int, wg *sync.WaitGroup, mu *sync.Mutex) {
				defer wg.Done()
				wginner := &sync.WaitGroup{}
				var func1 string
				var func2 string
				wginner.Add(2)
				datanew := fmt.Sprint(data)
				go func(data int, wg *sync.WaitGroup, str *string) {
					defer wg.Done()
					*str = DataSignerCrc32(datanew)
					runtime.Gosched()
				}(data, wginner, &func1)
				go func(data int, wg *sync.WaitGroup, str *string, mu *sync.Mutex) {
					defer wg.Done()
					mu.Lock()
					temp := DataSignerMd5(datanew)
					mu.Unlock()
					*str = DataSignerCrc32(temp)
					runtime.Gosched()
				}(data, wginner, &func2, mu)
				wginner.Wait()
				out <- func1 + "~" + func2
			}(data, wg, mu)
		}
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for dataRaw := range in {
		data, ok := dataRaw.(string)
		if !ok {
			fmt.Println("cant convert to string")
		} else {
			wg.Add(1)
			go func(data string, wg *sync.WaitGroup) {
				defer wg.Done()
				var strings = map[int]string{}
				var result string
				wginner := &sync.WaitGroup{}
				mu := &sync.Mutex{}
				for i := 0; i < 6; i++ {
					wginner.Add(1)
					go func(i int, data string, strings map[int]string, wg *sync.WaitGroup, mu *sync.Mutex) {
						defer wg.Done()
						r := DataSignerCrc32(fmt.Sprint(i) + data)
						mu.Lock()
						strings[i] = r
						mu.Unlock()
					}(i, data, strings, wginner, mu)
				}
				wginner.Wait()
				mu.Lock()
				for i := 0; i < 6; i++ {
					result = result + strings[i]
				}
				mu.Unlock()
				out <- result
			}(data, wg)
		}
	}
	wg.Wait()
}

type SortBy []string

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i] < a[j] }

func CombineResults(in, out chan interface{}) {
	str := []string{}
	for value := range in {
		str = append(str, value.(string))
	}
	sort.Sort(SortBy(str))
	var result string
	for i, elem := range str {
		if i > 0 {
			result = result + "_" + elem
		} else {
			result = result + elem
		}
	}
	out <- result
}
