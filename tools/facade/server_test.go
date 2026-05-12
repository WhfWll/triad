// Package facade
// @Author bcy2007  2025/12/22 17:10
package facade

import (
	"fmt"
	"testing"
	"time"
)

func TestFacadeServer(t *testing.T) {
	NewFacadeServer("127.0.0.1", 7890, 30)
	go func() {
		for {
			fmt.Println("read...")
			messages, total := ReadMessage(1, 10)
			for _, message := range messages {
				fmt.Println(message)
			}
			fmt.Println(total)
			fmt.Println(errInfo)
			time.Sleep(3 * time.Second)
		}
	}()
	time.Sleep(40 * time.Second)
}
