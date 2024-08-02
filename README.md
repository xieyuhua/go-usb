# GO usb relay

GoLang implementation of usb-relay for control USB HID relay module.

A driver for control usb relay module implement with Golang.
We access devices through karalabe/hid

## Support
|    OS     |  Is supported |
|:---------:|:-------------:|
| MacOS     |  Yes          |
| Windows   |  Yes          |
| GNU/Linux |  Yes ( No test yet ) |

## Install
```sh
# go build
```

## Usage
```sh
# Trun On channel 1 relay.
$ go-usb -n 1 -o

# Set SN
$ go-usb -sn 12345
```
