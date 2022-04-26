package quit_routine_by_send_value

import (
	"fmt"
	"time"
)

// generate number and quit on receive
func generator(start, growBy int) chan int {
	ch := make(chan int)
	go func() {
		start := 1
		for {
			select {
			case ch <- start:
				start += growBy
			case <-ch: //quit on received
				close(ch)
				return
			}
		}
	}()
	return ch
}

func UseGenerator() {

	numbers := generator(0, 2)

	for n := range numbers {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(n)

		//when we want to stop
		// send to channel to notify stop goroutine
		// prevent go routine leak
		if n == 10 {
			numbers <- 0
		}
	}

	fmt.Println("Program exited!")
}
