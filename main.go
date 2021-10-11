package main
import (
	"fmt"
	start "tgcontextbot/internal/startup"
)

func main() {
	err, _ :=  start.InitializeBot()

	if err != nil  {
		fmt.Println("Unable to start up, terminating.")
		return
	}

	fmt.Println("Started up")
	return

}
