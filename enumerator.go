package udev

import (
	"C"
	"unsafe"
)
import (
	"fmt"

	ffi "github.com/gogogoghost/libffigo"
)

type UDevEnumerator struct {
	ctx *UDevContext
	ptr unsafe.Pointer
}

func (obj *UDevEnumerator) AddFilter(subSystem string) error {
	var subSystemPtr unsafe.Pointer
	if len(subSystem) > 0 {
		subSystemPtr := C.CString(subSystem)
		defer ffi.FreePtr(unsafe.Pointer(subSystemPtr))
	}
	res := udev_enumerate_add_match_subsystem(
		obj.ptr,
		subSystemPtr,
	)
	if res != 0 {
		return fmt.Errorf("add filter return:%d", res)
	}
	return nil
}

func (obj *UDevEnumerator) List() []*UDevice {
	udev_enumerate_scan_devices(obj.ptr)
	entry := udev_enumerate_get_list_entry(obj.ptr)
	var devList []*UDevice
	for entry != nil {
		//获取device
		name := udev_list_entry_get_name(entry)
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
