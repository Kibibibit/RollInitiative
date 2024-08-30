package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt)
	go func() {
		<-closeSignal
		fmt.Println("Interrupt")
		os.Exit(1)
	}()

}
