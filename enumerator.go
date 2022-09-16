package udev

import (
	"C"
	"unsafe"
)
import "errors"

type UDevEnumerator struct {
	ctx *UDevContext
	ptr unsafe.Pointer
}

func NewEnumerator(ctx *UDevContext) (obj *UDevEnumerator, err error) {
	enumerate := Udev_enumerate_new(ctx.ptr)
	if enumerate == nil {
		return nil, errors.New("fail to create enumerate")
	}
	return &UDevEnumerator{
		ctx: ctx,
		ptr: enumerate,
	}, nil
}

func (obj *UDevEnumerator) List() []*UDevice {
	Udev_enumerate_scan_devices(obj.ptr)
	entry := Udev_enumerate_get_list_entry(obj.ptr)
	var devList []*UDevice
	for entry != nil {
		//获取device
		name := Udev_list_entry_get_name(entry)
		dev := Udev_device_new_from_syspath(
			obj.ctx.ptr,
			name,
		)
		//获取props
		propEntry := Udev_device_get_properties_list_entry(
			dev,
		)
		env := make(map[string]string)
		for propEntry != nil {
			key := Udev_list_entry_get_name(propEntry)
			value := Udev_list_entry_get_value(propEntry)
			env[key] = value
			propEntry = Udev_list_entry_get_next(propEntry)
		}
		entry = Udev_list_entry_get_next(entry)
		devList = append(devList, &UDevice{
			SubSystem: env["SUBSYSTEM"],
			Env:       env,
		})
	}
	return devList
}
