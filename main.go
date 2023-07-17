package main

import (
	"fmt"

	udev "github.com/gogogoghost/libudevgo/lib"
)

func main() {
	udev.SetLibName("libudev.so.1")
	udev.Init()
	ctx, err := udev.NewContext()
	if err != nil {
		panic(err)
	}
	monitor, err := ctx.NewMonitor(udev.UDEV)
	if err != nil {
		panic(err)
	}
	monitor.AddFilter("tty", "")
	monitor.AddFilter("block", "")
	monitorChan, err := monitor.StartMonitor()
	if err != nil {
		panic(err)
	}
	for {
		evt := <-monitorChan
		fmt.Printf("%s %s\n", evt.Action, evt.Device.Env["DEVNAME"])
	}
	// enumerator, err := ctx.NewEnumerator()
	// if err != nil {
	// 	panic(err)
	// }
	// err = enumerator.AddFilter("tty")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, dev := range enumerator.List() {
	// 	if dev.Env["DEVNAME"] == "/dev/ttyACM0" {
	// 		fmt.Println(dev.Env)
	// 	}
	// }
}
