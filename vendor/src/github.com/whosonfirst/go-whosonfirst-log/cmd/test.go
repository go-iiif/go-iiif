package main

import (
       "fmt"
       "io"
       "os"
       "github.com/whosonfirst/go-whosonfirst-log"
)

func woo (log log.WOFLog){

     fmt.Println("woo")
     log.Status("wof %v", log)
}

func main() {

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger()
	logger.AddLogger(writer, "debug")

	_, err := logger.AddLogger(writer, "debug")

	if err != nil {
	   panic(err)
	}

	logger.Info("Writing all your logs to %s", "wub wub wub")
	logger.Debug("Hello world")

	logger.Status("STATUS")

	mock := &log.MockLogger{}
	mock.Debug("hello %s", "world")

	fmt.Println("logger")
	woo(logger)

	fmt.Println("mock")
	woo(mock)

	fmt.Println("DONE")

	l := log.SimpleWOFLogger("foo")
	l.Status("FOO")
	l.Info("Info")
}
