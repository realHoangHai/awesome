package log_test

import (
	"github.com/realHoangHai/awesome/pkg/log"
	"os"
	"time"
)

func Example() {
	log.Info("Hello, micro!")
	log.Infof("Hello, %s!", "jack")
	log.Warn("This is a warning")
	log.Debug("hello micro debug")
	log.Error("this is an error")
	//log.Panic("this is a panic")
	log.Fields("name", "micro", "site", "VN").Info("this is a log with fields")
}

func ExampleInit_withOptions() {
	log.Init(
		log.WithFields("name", "my service"),
		log.WithFormat(log.FormatJSON),
		log.WithLevel(log.LevelDebug),
		log.WithTimeFormat(time.RFC1123),
		log.WithWriter(os.Stdout),
	)
	log.Info("hello")
}
