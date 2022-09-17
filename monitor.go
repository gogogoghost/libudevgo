package udev

import (
	//#include "helper.h"
	//
	"C"
	"unsafe"
)
import (
	"errors"
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

func NewMonitor(ctx *UDevContext, monitorType MonitorType) (mon *UDevMonitor, err error) {
	ptr := Udev_monitor_new_from_netlink(ctx.ptr, string(monitorType))
	if ptr == nil {
		return nil, errors.New("fail to create monitor")
	}
	return &UDevMonitor{
		ctx: ctx,
		ptr: ptr,
		fd:  -1,
	}, nil
}

func (obj *UDevMonitor) AddFilter(subSystem string, devType string) error {
	var subSystemPtr unsafe.Pointer
	var devTypePtr unsafe.Pointer
	if len(subSystem) > 0 {
		subSystemPtr := C.CString(subSystem)
		defer ffi.FreePtr(unsafe.Pointer(subSystemPtr))
	}
	if len(devType) > 0 {
		devTypePtr := C.CString(subSystem)
		defer ffi.FreePtr(unsafe.Pointer(devTypePtr))
	}
	res := Udev_monitor_filter_add_match_subsystem_devtype(
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
	res := Udev_monitor_enable_receiving(obj.ptr)
	if res != 0 {
		return nil, fmt.Errorf("enable receiving return:%d", res)
	}
	fd := Udev_monitor_get_fd(obj.ptr)
	if fd < 0 {
		return nil, fmt.Errorf("fail to get fd:%d", fd)
	}
	obj.fd = int(fd)
	channel := make(chan UEvent)
	// 开始后台poll
	go obj.poll(channel)
	return channel, nil
}

func (obj *UDevMonitor) poll(channel chan UEvent) {
	for obj.fd > 0 {
		res := int(C.poll_fd(C.int(obj.fd), -1))
		// 返回-1 结束
		if res < 0 {
			break
		}
		// 为0 读取数据
		device := Udev_monitor_receive_device(obj.ptr)
		if device == nil {
			continue
		}
		// action := Udev_device_get_action(device)
		env := make(map[string]string)
		propEntry := Udev_device_get_properties_list_entry(device)
		for propEntry != nil {
			key := Udev_list_entry_get_name(propEntry)
			value := Udev_list_entry_get_value(propEntry)
			env[key] = value
			propEntry = Udev_list_entry_get_next(propEntry)
		}
		channel <- UEvent{
			Action: env["ACTION"],
			Device: UDevice{
				SubSystem: env["SUBSYSTEM"],
				Env:       env,
			},
		}
	}
}