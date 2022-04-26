package merge_channel

import (
	"fmt"
	"reflect"
	"sync"
)

func UseMergeNChannel() {
	ch1 := Produce(1, 2, 3, 4)
	ch2 := Produce(5, 6, 7, 8)
	ch3 := Produce(9, 10, 11, 12)

	// Print(MergeTwo(ch1, ch2))
	// Print(MergeN(ch1, ch2))

	Print(MergeRelfect(ch1, ch2, ch3))
}

// Worker
func Produce(sequence ...int) <-chan int {
	ch := make(chan int)

	go func() {
		defer func() { close(ch) }()

		for _, v := range sequence {
			ch <- v
		}
	}()

	return ch
}

// Multiplexing 2 chan
func MergeTwo(ch1, ch2 <-chan int) <-chan int {

	ch := make(chan int)
	go func() {
		defer func() { close(ch) }()

		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					fmt.Println("ch1", v)
					ch1 = nil
					continue
				}
				ch <- v
			case v, ok := <-ch2:
				if !ok {
					fmt.Println("ch2", v)
					ch2 = nil
					continue
				}
				ch <- v
			}
		}
	}()

	return ch
}

func Print(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

// Perfomance not relly good
// Easy to read
func MergeN(chs ...<-chan int) <-chan int {
	out := make(chan int)

	go func() {
		wg := sync.WaitGroup{}
		defer func() { close(out) }()
		wg.Add(len(chs))
		for _, ch := range chs {
			go func(ch <-chan int) {
				defer wg.Done()
				for v := range ch {
					out <- v
				}
			}(ch)
		}

		wg.Wait()
	}()
	return out
}

//High perfomance
//Harder to read
func MergeRelfect(chs ...<-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer func() { close(out) }()
		cases := []reflect.SelectCase{}

		for _, ch := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			// closed chan
			if !ok {
				if i != len(cases) {
					cases = append(cases[:i], cases[i+1:]...)
				} else {
					cases = cases[i : len(cases)-1]
				}
			} else {
				out <- int(v.Int())
			}
		}
	}()

	return out
}
