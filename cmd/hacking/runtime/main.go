package main

import (
	"fmt"
	"runtime"
)

func call() {
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	pc = pc[:n]
	fn := runtime.FuncForPC(pc[0])
	file, line := fn.FileLine(pc[0])
	fmt.Println(fn.Name(), file, line)

	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		fmt.Printf("- more:%v | %s | %s | %s | %d \n", more, frame.Function, frame.Func.Name(), frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func First() {
	Second()
}

func Second() {
	Third()
}

func Third() {
	for c := 0; c < 5; c++ {
		fmt.Println(runtime.Caller(c))
	}
}

func stackExample() {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	fmt.Printf("\n%s", stackSlice[:s])
}

func main() {
	call()
	stackExample()
	First()
}
