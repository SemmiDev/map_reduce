package main

import (
	"fmt"
	"sync"
)

func inputData(out[3]chan <- string) {
	data := [][]string{
		{"izzah", "izzah","izzah","sammi","sammi", "sammi"},
		{"sammi", "sammi","sammi","izzah","izzah", "izzah"},
		{"izzah", "izzah","izzah","sammi","sammi", "sammi"},
	}

	for i := range data {
		go func(ch chan <- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], data[i])
	}
}

func mapper(in <- chan string, out chan <- map[string]int) {
	count := map[string]int{}
	for word := range in {
		count[word] = count[word] + 1
	}
	out <- count
	close(out)
}

func reducer(in <-chan int, out chan<- float32) {
	sum, count := 0, 0
	for n := range in {
		sum += n
		count++
	}
	out <- float32(sum) / float32(count)
	close(out)
}

func shuffler(in []<-chan map[string]int, out [2]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["sammi"]
				if ok {
					out[0] <- nc
				}
				vc, ok := m["izzah"]
				if ok {
					out[1] <- vc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
	}()
}
func outputWriter(in []<-chan float32) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	name := []string{"sammi", "izzah"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for avg := range c {
				fmt.Printf("Average number of %ss per input text: %f\n", name[n], avg)
			}
			wg.Done()
		}(i, in[i])
	}
	wg.Wait()
}

func main() {
	size := 10

	// -------------------------------------------
	s1 := make(chan string, size)
	s2 := make(chan string, size)
	s3 := make(chan string, size)

	text := [3]chan <- string{s1, s2, s3}
	go inputData(text)

	// -------------------------------------------
	map1 := make(chan map[string]int, size)
	map2 := make(chan map[string]int, size)
	map3 := make(chan map[string]int, size)

	go mapper(s1, map1)
	go mapper(s2, map2)
	go mapper(s3, map3)

	// -------------------------------------------
	reduce1 := make(chan int, size)
	reduce2 := make(chan int, size)

	go shuffler([]<-chan map[string]int{map1, map2, map3}, [2]chan<- int{reduce1, reduce2})
	// -------------------------------------------

	avg1 := make(chan float32, size)
	avg2 := make(chan float32, size)

	go reducer(reduce1, avg1)
	go reducer(reduce2, avg2)

	// -------------------------------------------
	outputWriter([]<-chan float32{avg1, avg2})
}