package udev

import (
	"C"
	"unsafe"
)
import (
	"errors"
	"fmt"
)

type UDevEnumerator struct {
	ctx *UDevContext
	ptr unsafe.Pointer
}

func NewEnumerator(ctx *UDevContext) (obj *UDevEnumerator, err error) {
	enumerate := udev_enumerate_new(ctx.ptr)
	if enumerate == nil {
		return nil, errors.New("fail to create enumerate")
	}
	return &UDevEnumerator{
		ctx: ctx,
		ptr: enumerate,
	}, nil
}

func (obj *UDevEnumerator) List() []*UDevice {
	udev_enumerate_scan_devices(obj.ptr)
	entry := udev_enumerate_get_list_entry(obj.ptr)
	var devList []*UDevice
	for entry != nil {
		//获取device
		name := udev_list_entry_get_name(entry)
		fmt.Println(name)
		dev := udev_device_new_from_syspath(
			obj.ctx.ptr,
			name,
		)
		//获取props
		propEntry := udev_device_get_properties_list_entry(
			dev,
		)
		env := make(map[string]string)
		for propEntry != nil {
			key := udev_list_entry_get_name(propEntry)
			value := udev_list_entry_get_value(propEntry)
			env[key] = value
			propEntry = udev_list_entry_get_next(propEntry)
		}
		entry = udev_list_entry_get_next(entry)
		devList = append(devList, &UDevice{
			SubSystem: env["SUBSYSTEM"],
			Env:       env,
		})
	}
	return devList
}
