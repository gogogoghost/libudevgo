package udev

import (
	"errors"
	"unsafe"
)

type UDevContext struct {
	ptr unsafe.Pointer
}

func NewContext() (obj *UDevContext, err error) {
	ctx := Udev_new.Call().Pointer()
	if ctx == nil {
		return nil, errors.New("fail to create context")
	}
	return &UDevContext{
		ptr: ctx,
	}, nil
}
