package main

import (
	"fmt"
	chat "tgcontextbot/internal/chatStorer"
	start "tgcontextbot/internal/startup"
)





func main() {
	chat.BuildStorer()
	err, bot :=  start.InitializeBot()

	if err != nil  {
		fmt.Println("Unable to start up, terminating.")
		return
	}

	fmt.Println("Started up")

	start.ServeBot(bot)

}
