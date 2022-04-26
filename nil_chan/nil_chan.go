package nil_chan

import (
	"fmt"
	"log"
)

func NilChan() {
	var ch1 = Produce(5, 3, 5, 6, 12, 453)
	var ch2 = Produce(100, 423, 54, 32, 5343, 343)
	c := merge(ch1, ch2)
	Printch(c)
}

func Produce(list ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer func() {
			log.Println("Produce worker close")
			close(ch)
		}()
		for _, v := range list {
			ch <- v
		}
	}()
	return ch
}

func merge(ch1, ch2 <-chan int) <-chan int {
	c := make(chan int)

	go func() {
		defer func() {
			log.Println("merge worker close")
			close(c)
		}()

		for ch1 != nil || ch2 != nil {
			select {
			case i, ok := <-ch1:
				if !ok {
					ch1 = nil
					fmt.Println("case 1")
					continue
				}
				c <- i

			case i, ok := <-ch2:
				if !ok {
					ch2 = nil
					fmt.Println("case 2")
					continue

				}
				c <- i

			}
		}
	}()

	return c
}

func Printch(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
