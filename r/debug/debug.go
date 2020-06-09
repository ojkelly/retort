package debug // import "retort.dev/r/debug"

import (
	"fmt"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

const debugLogPath = "debug.log"
const debugLineBreak = "\n--------------------------------------------------------------------------------\n"

var debugMutex = &sync.Mutex{}
var debugFile *os.File

func init() {
	f, err := os.OpenFile(debugLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		panic(f)
	}
	debugFile = f
	// defer f.Close()
}

func Log(message ...interface{}) {
	debugMutex.Lock()
	debugFile.WriteString(debugLineBreak)

	debugFile.WriteString(fmt.Sprint(message...))

	debugMutex.Unlock()
}

func Spew(message ...interface{}) {
	debugMutex.Lock()
	debugFile.WriteString(debugLineBreak)

	spew.Fdump(debugFile, message...)
	debugMutex.Unlock()
}
