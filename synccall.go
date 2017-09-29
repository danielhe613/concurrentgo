/*
This package is used to automatically backup network devices' configuration by ssh.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const WorkingTime = 1 * time.Second
const EchoTimeout = 2 * time.Second

func main() {
	result, err := Echo("input", EchoTimeout)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	time.Sleep(2 * time.Second)
}

func Echo(in string, timeout time.Duration) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	out := make(chan string, 1) //bufferred channel, the size setting is as per the result size.

	go doEcho(ctx, out, in)

	select {
	case result := <-out:
		//Processing result if needed…
		result += "_echo"

		cancel() //to notify doEcho goroutine to exit
		fmt.Println("Echo exits normally!")
		return result, nil

	case <-ctx.Done(): //timeout
		//Just exits
		fmt.Println("Echo exits with timeout!")
		return "", errors.New("Timeout")
	}
}

func doEcho(ctx context.Context, out chan<- string, in string) {
	defer close(out)

	//Solution1:
	select {
	case <-ctx.Done(): //canceled or timeout     similar to if stopFlag, non-blocking check
		return
	default:
	}
	// blocked work…
	time.Sleep(WorkingTime)
	out <- "result=" + in //out is bufferred, so don't need to block here when timeout occurres first.

	//block-waiting for the context done (A cancel context or timeout)。
	<-ctx.Done()

	//Solution2:
	// for {
	// 	select {
	// 	case <-ctx.Done(): //canceled or timeout 类似 if stopFlag
	// 		return
	// 	default:
	// 		//do partial work each time
	// 		if complete {
	// 			out <- result //out is bufferred, so don't need to block here when timeout occurres first.
	// 			<-ctx.Done()  //Waiting for the context done (A cancel context or timeout)。其实还是类似stopflag                                           return
	// 		}
	// 	}
	// }

	fmt.Println("doEcho exits!")
}
