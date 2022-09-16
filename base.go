package udev

import (
	"github.com/gogogoghost/libffigo/ffi"
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
	Udev_new                                        *ffi.Cif
	Udev_enumerate_new                              *ffi.Cif
	Udev_enumerate_scan_devices                     *ffi.Cif
	Udev_enumerate_get_list_entry                   *ffi.Cif
	Udev_device_new_from_syspath                    *ffi.Cif
	Udev_list_entry_get_name                        *ffi.Cif
	Udev_device_get_properties_list_entry           *ffi.Cif
	Udev_list_entry_get_value                       *ffi.Cif
	Udev_list_entry_get_next                        *ffi.Cif
	Udev_monitor_new_from_netlink                   *ffi.Cif
	Udev_monitor_enable_receiving                   *ffi.Cif
	Udev_monitor_filter_add_match_subsystem_devtype *ffi.Cif
	Udev_monitor_get_fd                             *ffi.Cif
	Udev_monitor_receive_device                     *ffi.Cif
	Udev_device_get_action                          *ffi.Cif
)

func Init() {
	lib, err := ffi.Open("libudev.so", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	Udev_new = lib.SymMust("udev_new", ffi.PTR)
	Udev_enumerate_new = lib.SymMust("udev_enumerate_new", ffi.PTR, ffi.PTR)
	Udev_enumerate_scan_devices = lib.SymMust(
		"udev_enumerate_scan_devices", ffi.SINT32, ffi.PTR)
	Udev_enumerate_get_list_entry = lib.SymMust(
		"udev_enumerate_get_list_entry", ffi.PTR, ffi.PTR)
	Udev_device_new_from_syspath = lib.SymMust(
		"udev_device_new_from_syspath",
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	Udev_list_entry_get_name = lib.SymMust(
		"udev_list_entry_get_name",
		ffi.PTR,
		ffi.PTR,
	)
	Udev_device_get_properties_list_entry = lib.SymMust(
		"udev_device_get_properties_list_entry",
		ffi.PTR,
		ffi.PTR,
	)
	Udev_list_entry_get_value = lib.SymMust(
		"udev_list_entry_get_value",
		ffi.PTR,
		ffi.PTR,
	)
	Udev_list_entry_get_next = lib.SymMust(
		"udev_list_entry_get_next",
		ffi.PTR,
		ffi.PTR,
	)
	Udev_monitor_new_from_netlink = lib.SymMust(
		"udev_monitor_new_from_netlink",
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	Udev_monitor_enable_receiving = lib.SymMust(
		"udev_monitor_enable_receiving",
		ffi.SINT32,
		ffi.PTR,
	)
	Udev_monitor_filter_add_match_subsystem_devtype = lib.SymMust(
		"udev_monitor_filter_add_match_subsystem_devtype",
		ffi.SINT32,
		ffi.PTR,
		ffi.PTR,
		ffi.PTR,
	)
	Udev_monitor_get_fd = lib.SymMust(
		"udev_monitor_get_fd",
		ffi.SINT32,
		ffi.PTR,
	)
	Udev_monitor_receive_device = lib.SymMust(
		"udev_monitor_receive_device",
		ffi.PTR,
		ffi.PTR,
	)
	Udev_device_get_action = lib.SymMust(
		"udev_device_get_action",
		ffi.PTR,
		ffi.PTR,
	)
}
