package sge

import (
	"runtime"
)

var glThread chan func() = make(chan func(), 100)
var GL chan<- func() = glThread

func runGL() {
	runtime.LockOSThread()
	for {
		(<-glThread)()
	}
}

func init() {
	go runGL()
}
