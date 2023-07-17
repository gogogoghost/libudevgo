## libudev binding for go

[![GoDoc](https://godoc.org/github.com/gogogoghost/udev?status.svg)](https://godoc.org/github.com/gogogoghost/udev)

loading libudev.so with dynamic loading at runtime by [libffigo](https://github.com/gogogoghost/libffigo)

### usage

get this module

```sh
go get github.com/gogogoghost/udev
```

init and create context

```go
udev.init()
ctx, err := udev.NewContext()
if err != nil {
    panic(err)
}
```

enumerate device

```go
enumerator, err := udev.NewEnumerator(ctx)
if err != nil {
    panic(err)
}
for _, dev := range enumerator.List() {
    //read device
}
```
monitor event

```go
monitor, err := udev.NewMonitor(ctx, udev.UDEV)
if err != nil {
    panic(err)
}
monitor.AddFilter("tty", "")
channel, err := monitor.StartMonitor()
if err != nil {
    panic(err)
}
for {
    evt := <-channel
    //read event and device
}
```