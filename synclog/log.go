package synclog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var LogWriter io.Writer

var lock sync.Mutex

func log(format string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if LogWriter == nil {
		return
	}
	fmt.Fprintf(LogWriter, format, args...)
}

func LogToFile(path string) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	LogWriter = f
	return nil
}
