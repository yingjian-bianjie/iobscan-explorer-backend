package main

import (
	"github.com/bianjieai/iobscan-explorer-backend/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)
	cmd.Execute()
}
