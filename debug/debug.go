package debug // import "retort.dev/debug"

import (
	"fmt"
	"os"
	"sync"
	// "github.com/davecgh/go-spew/spew"
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

func Spew(message ...interface{}) {
	return
	// start := time.Now()

	// go func() {
	debugMutex.Lock()
	debugFile.WriteString(debugLineBreak)

	debugFile.WriteString(fmt.Sprint(message...))
	// spew.Fdump(debugFile, message...)

	// debugFile.WriteString(fmt.Sprintf("debug print time %s", time.Since(start)))
	debugMutex.Unlock()
	// }()

	// fmt.Println(message...)
}
