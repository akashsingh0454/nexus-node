package logger

import (
	"log"
	"os"
)

// Setup configures the global logger.
func Setup() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
