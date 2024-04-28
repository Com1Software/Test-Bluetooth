package main

import (
	"os"
	"strconv"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var DeviceAddress string

func connectAddress() string {
	if len(os.Args) < 2 {
		println("usage: discover [address]")
		os.Exit(1)
	}
	address := os.Args[1]
	return address
}

func wait() {
	time.Sleep(3 * time.Second)
}

func done() {
	println("Done.")

	time.Sleep(1 * time.Hour)
}

func main() {
	wait()

	println("enabling")
	must("enable BLE stack", adapter.Enable())
	ch := make(chan bluetooth.ScanResult, 1)

	println("scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		println("found device:", result.Address.String(), result.RSSI, result.LocalName())
		if result.Address.String() == connectAddress() {
			adapter.StopScan()
			ch <- result
		}
	})

	var device bluetooth.Device
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())
			return
		}

		println("connected to ", result.Address.String())
	}

	println("discovering services/characteristics")
	srvcs, err := device.DiscoverServices(nil)

	must("discover services", err)
	buf := make([]byte, 255)
	for _, srvc := range srvcs {
		println("- service", srvc.UUID().String())
		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			println(err)
		}
		for _, char := range chars {
			println("-- characteristic", char.UUID().String())
			mtu, err := char.GetMTU()
			if err != nil {
				println("    mtu: error:", err.Error())
			} else {
				println("    mtu:", mtu)
			}
			n, err := char.Read(buf)
			if err != nil {
				println("    ", err.Error())
			} else {
				println("    data bytes", strconv.Itoa(n))
				println("    value =", string(buf[:n]))
			}
		}
	}

	err = device.Disconnect()
	if err != nil {
		println(err)
	}

	done()
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
