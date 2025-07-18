package main

import (
	"github.com/jasper-186/lionchief"
	"tinygo.org/x/bluetooth"
)

func manuallySetTrain() {
	// After you determine the name of your train set it here
	trainMac := "44:A6:E5:41:AE:72"
	// Enable BLE interface.
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		panic("failed to " + "enable BLE stack" + ": " + err.Error())
	} else {
		println("enable BLE stack")
	}

	address, _ := bluetooth.ParseMAC(trainMac)
	trainAddress := bluetooth.Address{
		MACAddress: bluetooth.MACAddress{MAC: address},
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
