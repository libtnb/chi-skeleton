package main

import "runtime/debug"

func main() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(128 << 20)

	app, err := initApp()
	if err != nil {
		panic(err)
	}

	app.Run()
}
