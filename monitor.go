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

	"github.com/gogogoghost/libffigo/ffi"
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
	name := C.CString(string(monitorType))
	defer ffi.FreePtr(unsafe.Pointer(name))
	ptr := Udev_monitor_new_from_netlink.Call(ctx.ptr, name).Pointer()
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
	res := Udev_monitor_filter_add_match_subsystem_devtype.Call(
		self.ptr,
		subSystemPtr,
		devTypePtr,
	).Int32()
	if res != 0 {
		return fmt.Errorf("add filter return:%d", res)
	}
	return nil
}

func (self *UDevMonitor) StartMonitor() (chan UEvent, error) {
	res := Udev_monitor_enable_receiving.Call(self.ptr).Int32()
	if res != 0 {
		return nil, fmt.Errorf("enable receiving return:%d", res)
	}
	fd := Udev_monitor_get_fd.Call(self.ptr).Int32()
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
		device := Udev_monitor_receive_device.Call(self.ptr).Pointer()
		if device == nil {
			continue
		}
		action := Udev_device_get_action.Call(device).String()
		env := make(map[string]string)
		propEntry := Udev_device_get_properties_list_entry.Call(device).Pointer()
		for propEntry != nil {
			key := Udev_list_entry_get_name.Call(propEntry).String()
			value := Udev_list_entry_get_value.Call(propEntry).String()
			env[key] = value
			propEntry = Udev_list_entry_get_next.Call(propEntry).Pointer()
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
