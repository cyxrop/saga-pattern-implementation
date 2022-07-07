package metrics

import (
	"runtime"
	"strconv"
)

type Goroutines struct{}

func (g *Goroutines) String() string {
	return strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
}
