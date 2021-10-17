package main

import (
	"fmt"
	handle "tgcontextbot/internal/handling"
	start "tgcontextbot/internal/startup"
)

func main() {
	err, bot := start.InitializeBot()

	if err != nil {
		fmt.Println("Unable to start up, terminating.")
		handle.HandleError(err)
		return
	}

	fmt.Println("Started up")

	Err := start.ServeBot(bot)
	if Err != nil {
		handle.HandleError(Err)
		return
	}

}
