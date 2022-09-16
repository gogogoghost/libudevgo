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
	enumerate := Udev_enumerate_new.Call(ctx.ptr).Pointer()
	if enumerate == nil {
		return nil, errors.New("fail to create enumerate")
	}
	return &UDevEnumerator{
		ctx: ctx,
		ptr: enumerate,
	}, nil
}

func (obj *UDevEnumerator) List() []*UDevice {
	Udev_enumerate_scan_devices.Call(obj.ptr).Int32()
	entry := Udev_enumerate_get_list_entry.Call(obj.ptr).Pointer()
	var devList []*UDevice
	for entry != nil {
		//获取device
		name := Udev_list_entry_get_name.Call(entry).Pointer()
		dev := Udev_device_new_from_syspath.Call(
			&obj.ctx.ptr,
			name,
		).Pointer()
		//获取props
		propEntry := Udev_device_get_properties_list_entry.Call(
			dev,
		).Pointer()
		env := make(map[string]string)
		for propEntry != nil {
			key := Udev_list_entry_get_name.Call(propEntry).String()
			value := Udev_list_entry_get_value.Call(propEntry).String()
			env[key] = value
			propEntry = Udev_list_entry_get_next.Call(propEntry).Pointer()
		}
		entry = Udev_list_entry_get_next.Call(entry).Pointer()
		devList = append(devList, &UDevice{
			SubSystem: env["SUBSYSTEM"],
			Env:       env,
		})
	}
	return devList
}
