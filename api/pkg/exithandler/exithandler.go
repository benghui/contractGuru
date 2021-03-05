package exithandler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Exit accepts a callback function which invokes when program exits unexpectedly or is terminated by the user.
func Exit(cb func()) {
	sigs := make(chan os.Signal, 1)
	terminate := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Exit reason: ", sig)
		terminate <- true
	}()

	<-terminate
	cb()
	log.Print("Exiting program")
}
