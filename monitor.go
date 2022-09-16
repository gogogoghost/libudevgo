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

	"github.com/gogogoghost/libffigo"
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

func (self *UDevMonitor) AddFilter(subSystem string, devType string) error {
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
		self.ptr,
		subSystemPtr,
		devTypePtr,
	)
	if res != 0 {
		return fmt.Errorf("add filter return:%d", res)
	}
	return nil
}

func (self *UDevMonitor) StartMonitor() (chan UEvent, error) {
	res := Udev_monitor_enable_receiving(self.ptr)
	if res != 0 {
		return nil, fmt.Errorf("enable receiving return:%d", res)
	}
	fd := Udev_monitor_get_fd(self.ptr)
	if fd < 0 {
		return nil, fmt.Errorf("fail to get fd:%d", fd)
	}
	self.fd = int(fd)
	channel := make(chan UEvent)
	// 开始后台poll
	go self.poll(channel)
	return channel, nil
}

func (self *UDevMonitor) poll(channel chan UEvent) {
	for self.fd > 0 {
		res := int(C.poll_fd(C.int(self.fd), -1))
		// 返回-1 结束
		if res < 0 {
			break
		}
		// 为0 读取数据
		device := Udev_monitor_receive_device(self.ptr)
		if device == nil {
			continue
		}
		action := Udev_device_get_action(device)
		env := make(map[string]string)
		propEntry := Udev_device_get_properties_list_entry(device)
		for propEntry != nil {
			key := Udev_list_entry_get_name(propEntry)
			value := Udev_list_entry_get_value(propEntry)
			env[key] = value
			propEntry = Udev_list_entry_get_next(propEntry)
		}
		channel <- UEvent{
			Action: action,
			Device: UDevice{
				SubSystem: env["SUBSYSTEM"],
				Env:       env,
			},
		}
	}
}
