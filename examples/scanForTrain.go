package main

import (
	"github.com/jasper-186/lionchief"
	"tinygo.org/x/bluetooth"
)

func scanForTrain() {

	adapter := bluetooth.DefaultAdapter

	// After you determine the name of your train set it here
	trainName := "LC-0-1-0429-754D"

	// Enable BLE interface.
	err := adapter.Enable()
	if err != nil {
		panic("failed to " + "enable BLE stack" + ": " + err.Error())
	} else {
		println("enable BLE stack")
	}

	// Start scanning.
	var trainAddress bluetooth.Address
	println("scanning...")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		println("found device:", device.Address.String(), device.RSSI, device.LocalName())

		if device.LocalName() == trainName {
			trainAddress = device.Address
			adapter.StopScan()
		}
	})

	if err != nil {
		panic("failed to " + "start scan" + ": " + err.Error())
	} else {
		println("start scan")
	}

	// Get a simulator
	simulator, err := lionchief.NewSimulator(trainAddress)

	if err != nil {
		panic("failed to " + "begin simulation" + ": " + err.Error())
	} else {
		println("begin simulation")
	}

	if simulator == nil {
		panic("Simulator is nil")
	}
	simulator.SoundBell(2)

}
