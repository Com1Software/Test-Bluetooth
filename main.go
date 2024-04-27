package main

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	must("enable BLE stack", adapter.Enable())
	test := 1
	switch {

	case test == 1:
		println("scanning...")
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if len(device.LocalName()) > 0 {
				println("found device:", device.Address.String(), device.RSSI, device.LocalName())
			}
		})

		must("start scan", err)

	case test == 0:
		println("scanning...")
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())
		})
		must("start scan", err)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
