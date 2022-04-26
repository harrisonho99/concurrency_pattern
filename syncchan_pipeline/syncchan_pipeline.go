package syncchan_pipeline

import "fmt"

type Done chan bool

type Pipe chan int

func UsePipe() {
	generarted := make(Pipe)
	squared := make(Pipe)
	done := make(Done)
	num := 10
	go Count(generarted, num)
	go Square(generarted, squared)
	go Print(squared, done)
	for range done {
	}
}

func Count(out chan<- int, n int) {
	for i := 0; i <= n; i++ {
		out <- i
	}
	//close pipe when done the work
	close(out)
}

//read from in, write to out
func Square(in <-chan int, out chan<- int) {
	for val := range in {
		out <- val * val
	}
	close(out)
}

func Print(in <-chan int, done Done) {
	for val := range in {
		fmt.Println(val)
	}
	close(done)
}
