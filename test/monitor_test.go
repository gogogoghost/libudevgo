package test

import (
	"fmt"
	"testing"

	"github.com/gogogoghost/udev"
)

func TestEnumerator(t *testing.T) {
	udev.Init()
	ctx, err := udev.NewContext()
	if err != nil {
		panic(err)
	}
	enumerator, err := ctx.NewEnumerator()
	if err != nil {
		panic(err)
	}
	err = enumerator.AddFilter("tty")
	if err != nil {
		panic(err)
	}
	for _, dev := range enumerator.List() {
		if dev.SubSystem == "tty" {
			fmt.Println(dev.Env["DEVNAME"])
		}
	}
}

func TestMonitor(t *testing.T) {
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
	channel, err := monitor.StartMonitor()
	if err != nil {
		panic(err)
	}
	for {
		res := <-channel
		fmt.Println(res)
	}
}
