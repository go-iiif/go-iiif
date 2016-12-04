# go-whosonfirst-log

A simple Go package to implement level-dependent logging.

## Usage

### Simple

```
import (
       "io"
       log "github.com/whosonfirst/go-whosonfirst-log"
)

loglevel := "info"
verbose := true

logfile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

if err != nil {
	panic(err)
}

writer := io.MultiWriter(logfile)

if verbose {
	writer = io.MultiWriter(os.Stdout, logfile)
	loglevel = "debug"
}

logger := log.NewWOFLogger("[your-app] ")
logger.AddLogger(writer, loglevel)

logger.Info("Writing all your logs to %s", logfile)
logger.Debug("Hello world")
```

### Fancy


```
import (
       "io"
       "os"
       "github.com/whosonfirst/go-whosonfirst-log"
       "github.com/whosonfirst/go-slackcat-writer"
)

func main() {

	writer := io.MultiWriter(os.Stdout)
	slack, _ := slackcat.NewWriter("~/.slackcat.conf")

	logger := log.NewWOFLogger("[your-app] ")

	logger.AddLogger(writer, "debug")
	logger.AddLogger(slack, "status")

	logger.Info("Writing all your logs to %s", "wub wub wub")
	logger.Debug("Hello world")
	logger.Status("Hello Slack")
}
```

### Messages

You can pass regular old `fmt.Sprintf` formats and parameters to any of the logging methods.
 
### Log levels

The following log levels are supported, in this order. That means if you specify that the default log level is `warning` any messages logged to the `Status` or `Info` or `Debug` methods will be ignored.

* debug
* info
* status
* warning
* error
* fatal

There are correspoding (public) instance methods available for each of these log levels. Invoking `logger.Fatal` will record the message to be logged and then call `os.Exit(1)`.

## See also

* http://dave.cheney.net/2015/11/05/lets-talk-about-logging
* https://github.com/whosonfirst/go-slackcat-writer
