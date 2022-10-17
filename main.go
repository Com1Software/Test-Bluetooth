package main

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	must("enable BLE stack", adapter.Enable())
	test:=1
	
	switch{
	
		case test ==1:
		adv := adapter.DefaultAdvertisement()
	         must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		         LocalName: "Go Bluetooth",
	                 }))
	must("start adv", adv.Start())

	println("advertising...")
	address, _ := adapter.Address()
	for {
		println("Go Bluetooth /", address.MAC.String())
		time.Sleep(time.Second)
	}
		case test==0:
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
