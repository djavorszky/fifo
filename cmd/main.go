package cmd

import (
	"log"
	"os"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "[fifo] ", log.Lshortfile|log.LstdFlags)

	logger.Print("Yay we're running")

}
