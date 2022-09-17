package udev

import (
	"unsafe"

	ffi "github.com/gogogoghost/libffigo"
)

// 表示一个设备
type UDevice struct {
	SubSystem string
	Env       map[string]string
}

// 表示一个事件
type UEvent struct {
	Action string
	Device UDevice
}

var (
	udev_new                                        func() unsafe.Pointer
	udev_enumerate_new                              func(unsafe.Pointer) unsafe.Pointer
	udev_enumerate_scan_devices                     func(unsafe.Pointer) int
	udev_enumerate_get_list_entry                   func(unsafe.Pointer) unsafe.Pointer
	udev_device_new_from_syspath                    func(unsafe.Pointer, string) unsafe.Pointer
	udev_list_entry_get_name                        func(unsafe.Pointer) string
	udev_device_get_properties_list_entry           func(unsafe.Pointer) unsafe.Pointer
	udev_list_entry_get_value                       func(unsafe.Pointer) string
	udev_list_entry_get_next                        func(unsafe.Pointer) unsafe.Pointer
	udev_monitor_new_from_netlink                   func(unsafe.Pointer, string) unsafe.Pointer
	udev_monitor_enable_receiving                   func(unsafe.Pointer) int
	udev_monitor_filter_add_match_subsystem_devtype func(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) int
	udev_monitor_get_fd                             func(unsafe.Pointer) int
	udev_monitor_receive_device                     func(unsafe.Pointer) unsafe.Pointer
)

func Init() {
	lib, err := ffi.Open("libudev.so", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	lib.SymMust("udev_new", &udev_new, ffi.PTR)
	lib.SymMust("udev_enumerate_new", &udev_enumerate_new, ffi.PTR, ffi.PTR)
	lib.SymMust(
		"udev_enumerate_scan_devices", &udev_enumerate_scan_devices, ffi.SINT32, ffi.PTR)
	lib.SymMust(
		"udev_enumerate_get_list_entry", &udev_enumerate_get_list_entry, ffi.PTR, ffi.PTR)
	lib.SymMust(
		"udev_device_new_from_syspath",
		&udev_device_new_from_syspath,
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_list_entry_get_name",
		&udev_list_entry_get_name,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_device_get_properties_list_entry",
		&udev_device_get_properties_list_entry,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_list_entry_get_value",
		&udev_list_entry_get_value,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_list_entry_get_next",
		&udev_list_entry_get_next,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_monitor_new_from_netlink",
		&udev_monitor_new_from_netlink,
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_monitor_enable_receiving",
		&udev_monitor_enable_receiving,
		ffi.SINT32,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_monitor_filter_add_match_subsystem_devtype",
		&udev_monitor_filter_add_match_subsystem_devtype,
		ffi.SINT32,
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_monitor_get_fd",
		&udev_monitor_get_fd,
		ffi.SINT32,
		ffi.PTR,
	)
	lib.SymMust(
		"udev_monitor_receive_device",
		&udev_monitor_receive_device,
		ffi.PTR,
		ffi.PTR,
	)
}
