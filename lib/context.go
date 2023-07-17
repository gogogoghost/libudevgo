package udev

import (
	"errors"
	"unsafe"
)

type UDevContext struct {
	ptr unsafe.Pointer
}

func NewContext() (obj *UDevContext, err error) {
	ctx := udev_new()
	if ctx == nil {
		return nil, errors.New("fail to create context")
	}
	return &UDevContext{
		ptr: ctx,
	}, nil
}

func (obj *UDevContext) NewEnumerator() (*UDevEnumerator, error) {
	enumerate := udev_enumerate_new(obj.ptr)
	if enumerate == nil {
		return nil, errors.New("fail to create enumerate")
	}
	return &UDevEnumerator{
		ctx: obj,
		ptr: enumerate,
	}, nil
}

func (obj *UDevContext) NewMonitor(monitorType MonitorType) (mon *UDevMonitor, err error) {
	ptr := udev_monitor_new_from_netlink(obj.ptr, string(monitorType))
	if ptr == nil {
		return nil, errors.New("fail to create monitor")
	}
	return &UDevMonitor{
		ctx: obj,
		ptr: ptr,
		fd:  -1,
	}, nil
}
