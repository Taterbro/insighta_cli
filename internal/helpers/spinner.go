package helpers

import (
	"fmt"
	"time"
)

func StartSpinner(stop chan bool) {
	chars := []rune{'|', '/', '-', '\\'}
	i := 0

	for {
		select {
		case <-stop:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\rLoading... %c", chars[i%len(chars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
