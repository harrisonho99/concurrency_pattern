package cancel_one_goroutine

import (
	"fmt"
	"os"
	"time"
)

func LaunchOrAbort() {

	// tick pattern
	// time.NewTicker to instantiate ticker
	ticker := time.NewTicker(time.Second)
	quit := make(chan bool)
	fmt.Println("Commencing countdown. Press return to abort.")
	abort := InputCancelationCommand(quit)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-ticker.C:
		case <-abort:
			fmt.Println("Aborted")
			return
		}
	}

	//stop ticker
	//prevent goroutine leak
	ticker.Stop()

	fmt.Println("rocket is launched")
}

func InputCancelationCommand(quit chan bool) <-chan struct{} {
	abort := make(chan struct{})
	go func(quit chan bool) {
		var enter byte = 10
		var c_return byte = 13
		buffer := make([]byte, 1)
		os.Stdin.Read(buffer) // read 1 byte
		if buffer[0] == enter || buffer[0] == c_return {
			abort <- struct{}{}
		}
	}(quit)
	return abort
}
