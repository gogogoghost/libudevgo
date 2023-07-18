package udev

import (
	//#include "helper.h"
	//
	"C"
	"unsafe"
)
import (
	"fmt"

	ffi "github.com/gogogoghost/libffigo"
)

type MonitorType string

const (
	KERNEL MonitorType = "kernel"
	UDEV   MonitorType = "udev"
)

type UDevMonitor struct {
	ctx *UDevContext
	ptr unsafe.Pointer
	fd  int
}

func (obj *UDevMonitor) AddFilter(subSystem string, devType string) error {
	var subSystemPtr unsafe.Pointer
	var devTypePtr unsafe.Pointer
	if len(subSystem) > 0 {
		subSystemPtr = unsafe.Pointer(C.CString(subSystem))
		defer ffi.FreePtr(subSystemPtr)
	}
	if len(devType) > 0 {
		devTypePtr = unsafe.Pointer(C.CString(subSystem))
		defer ffi.FreePtr(devTypePtr)
	}
	res := udev_monitor_filter_add_match_subsystem_devtype(
		obj.ptr,
		subSystemPtr,
		devTypePtr,
	)
	if res != 0 {
		return fmt.Errorf("add filter return:%d", res)
	}
	return nil
}

func (obj *UDevMonitor) StartMonitor() (chan UEvent, error) {
	res := udev_monitor_enable_receiving(obj.ptr)
	if res != 0 {
		return nil, fmt.Errorf("enable receiving return:%d", res)
	}
	fd := udev_monitor_get_fd(obj.ptr)
	if fd < 0 {
		return nil, fmt.Errorf("fail to get fd:%d", fd)
	}
	obj.fd = int(fd)
	channel := make(chan UEvent)
	// 开始后台poll
	go func() {
		for {
			res := obj.poll(channel)
			if res == 4 {
				fmt.Println("receive EINTR. retry")
				continue
			}
			fmt.Printf("Unexcept poll exit:%d\n", res)
			break
		}
	}()
	return channel, nil
}

func (obj *UDevMonitor) poll(channel chan UEvent) int {
	for obj.fd > 0 {
		res := int(C.poll_fd(C.int(obj.fd), -1))
		// 返回-1 结束
		if res < 0 {
			return int(C.get_errno())
		}
		// 为0 读取数据
		device := udev_monitor_receive_device(obj.ptr)
		if device == nil {
			continue
		}
		// action := Udev_device_get_action(device)
		env := make(map[string]string)
		propEntry := udev_device_get_properties_list_entry(device)
		for propEntry != nil {
			key := udev_list_entry_get_name(propEntry)
			value := udev_list_entry_get_value(propEntry)
			env[key] = value
			propEntry = udev_list_entry_get_next(propEntry)
		}
		channel <- UEvent{
			Action: env["ACTION"],
			Device: UDevice{
				SubSystem: env["SUBSYSTEM"],
				Env:       env,
			},
		}
	}
	return 0
}
