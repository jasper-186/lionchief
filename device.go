package lionchief

// Train Commands in HEX
// Set horn volume/pitch: 44 01 <00-0f> <fe-02>
// Set bell volume/pitch: 44 02 <00-0f> <fe-02>
// Set speech volume/pitch: 44 03 <00-0f> <fe-02>
// Set engine volume/pitch: 44 04 <00-0f> <fe-02>
// Set speed : 45 <00-1f>
// Forward : 46 01
// Reverse : 46 02
// Bell start: 47 01
// Bell stop : 47 00
// Horn start: 48 01
// Horn stop : 48 00
// Disconnect: 4b 0 0
// Set overall volume: 4c <00-07>
// Speech : 4d XX 00
// Set lights off: 51 00
// Set lights on: 51 01

import (
	"errors"
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
)

type TrainState struct {
	Speed        int
	Reverse      bool
	Light        bool
	Volume       int
	VolumeHorn   int
	VolumeEngine int
	VolumeBell   int
	VolumeSpeech int
}

type TrainEngine struct {
	device              *bluetooth.Device
	writeService        *bluetooth.DeviceService
	writeCharacteristic *bluetooth.DeviceCharacteristic
	state               *TrainState
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func calculateChecksum(cmdBuffer []byte) byte {
	sumValue := 0
	for _, value := range cmdBuffer {
		sumValue = (int(value) + sumValue)
	}
	return byte(uint(sumValue))
}

func NewEngineDefaultBluetoothAdapter(trainAddress bluetooth.Address) (*TrainEngine, error) {
	return NewEngine(trainAddress, bluetooth.DefaultAdapter)
}

func NewEngine(trainAddress bluetooth.Address, adapter *bluetooth.Adapter) (*TrainEngine, error) {

	connectionParams := bluetooth.ConnectionParams{}
	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		log.Println("Connection Handler")
		// if trainAddress.MAC.String() == device.Address.String() && !connected {
		// 	log.Printf("Reconnecting to '%s' \n", device.Address.MAC.String())
		// 	adapter.Connect(device.Address, connectionParams)
		// }
	})

	log.Println("Enabling Adapter")
	adapter.Enable()

	device, err := adapter.Connect(trainAddress, connectionParams)
	must(fmt.Sprintf("Connecting to '%v'\n", trainAddress.MAC.String()), err)

	devicesServices, err := device.DiscoverServices([]bluetooth.UUID{ReadWriteService})
	must("Discovering Services", err)

	log.Printf("Found '%v' services", len(devicesServices))
	if len(devicesServices) < 1 {
		return nil, errors.New("failed to find read/write service")
	}

	log.Println("Discovering Characteristics")
	characteristics, err := devicesServices[0].DiscoverCharacteristics([]bluetooth.UUID{WriteCharacteristic})
	if err != nil {
		return nil, err
	}

	log.Printf("Found '%v' characteristics", len(characteristics))
	if len(characteristics) < 1 {
		return nil, errors.New("write characteristic not found")
	}

	train := TrainEngine{
		device:              &device,
		writeService:        &devicesServices[0],
		writeCharacteristic: &characteristics[0],
		state: &TrainState{
			Speed:        0,
			Reverse:      false,
			Light:        true,
			Volume:       1,
			VolumeHorn:   1,
			VolumeEngine: 0,
			VolumeBell:   1,
			VolumeSpeech: 1,
			//VolumeChuff:  1,
		},
	}

	// Make sure the train is in the default state (specifically Volumes) before we return it
	err = train.ResetState()
	if err != nil {
		return nil, err
	}
	return &train, nil
}

func (a *TrainEngine) Disconnect() error {
	return a.device.Disconnect()
}

func (a *TrainEngine) ResetState() error {
	(*a).state.Speed = 0
	err := a.SetSpeed(0)
	if err != nil {
		return err
	}

	(*a).state.Reverse = false
	err = a.SetReverse(false)
	if err != nil {
		return err
	}

	(*a).state.Light = true
	err = a.SetLight(true)
	if err != nil {
		return err
	}

	(*a).state.Volume = 7
	err = a.SetMainVolume(7)
	if err != nil {
		return err
	}

	(*a).state.VolumeHorn = 7
	err = a.SetRunningVolume(SOUNDTYPE_HORN, 7)
	if err != nil {
		return err
	}

	(*a).state.VolumeEngine = 7
	err = a.SetRunningVolume(SOUNDTYPE_ENGINE, 7)
	if err != nil {
		return err
	}

	(*a).state.VolumeBell = 7
	err = a.SetRunningVolume(SOUNDTYPE_BELL, 7)
	if err != nil {
		return err
	}

	(*a).state.VolumeSpeech = 7
	err = a.SetRunningVolume(SOUNDTYPE_SPEECH, 7)
	if err != nil {
		return err
	}

	return nil
}

func (a *TrainEngine) sendCommand(cmdByteArray []byte) error {
	log.Println("sendCommand")
	checksumedCmd := make([]byte, len(cmdByteArray)+2)
	checksumedCmd[0] = 0
	// Copy the values but offset them by 1
	for i, v := range cmdByteArray {
		checksumedCmd[i+1] = v
	}

	checksumedCmd[len(cmdByteArray)+1] = calculateChecksum(cmdByteArray)
	written, err := a.writeCharacteristic.WriteWithoutResponse(checksumedCmd)

	if err != nil {
		return err
	}

	if written != len(checksumedCmd) {
		return fmt.Errorf("writing command only wrote '%v' bytes of '%v'", written, len(checksumedCmd))
	}
	log.Println("sendCommand-Done")
	return nil
}

func (a *TrainEngine) SetMainVolume(volume int) error {
	log.Println("SetMainVolume")
	defer log.Println("SetMainVolume-Done")
	min := int(0)
	max := int(7)
	if volume > max || volume < min {
		return fmt.Errorf("invalid volume, must be between '%d' and '%d' (inclusive)", min, max)
	}

	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_SOUND_MAIN)
	cmdArray[1] = byte(volume)
	err := a.sendCommand(cmdArray)
	if err == nil {
		(*a).state.Volume = volume
	}
	return err
}

func (a *TrainEngine) SetRunningVolume(soundtype SoundType, volume int) error {
	log.Printf("SetRunningVolume-%v\n", soundtype)
	defer log.Printf("SetRunningVolume-%v-Done\n", soundtype)
	min := int(0)
	max := int(13)
	if volume > max || volume < min {
		return fmt.Errorf("invalid volume, must be between '%d' and '%d' (inclusive)", min, max)
	}

	cmdArray := make([]byte, 3)
	cmdArray[0] = byte(COMMANDTYPE_SOUND_RUNNING)
	cmdArray[1] = byte(soundtype)
	cmdArray[2] = byte(volume)
	err := a.sendCommand(cmdArray)
	if err == nil {
		switch soundtype {
		case SOUNDTYPE_BELL:
			(*a).state.VolumeBell = volume
		case SOUNDTYPE_ENGINE:
			(*a).state.VolumeEngine = volume
		case SOUNDTYPE_HORN:
			(*a).state.VolumeHorn = volume
		case SOUNDTYPE_SPEECH:
			(*a).state.VolumeSpeech = volume
		}
	}
	return err
}

func (a *TrainEngine) SetRunningPitch(soundtype SoundType, pitch SoundPitch) error {
	log.Println("SetRunningPitch")
	defer log.Println("SetRunningPitch-Done")
	cmdArray := make([]byte, 4)
	cmdArray[0] = byte(COMMANDTYPE_SOUND_RUNNING)
	cmdArray[1] = byte(soundtype)
	cmdArray[2] = byte(14)
	cmdArray[3] = byte(pitch)
	err := a.sendCommand(cmdArray)
	return err
}

func (a *TrainEngine) SetSpeed(speed int) error {
	log.Println("SetSpeed")
	defer log.Println("SetSpeed-Done")
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_SPEED)
	cmdArray[1] = byte(speed)
	err := a.sendCommand(cmdArray)
	(*a).state.Speed = speed
	return err
}

func (a *TrainEngine) GetSpeed() int {
	return (*a).state.Speed
}

func (a *TrainEngine) SetHorn(enabled bool) error {
	log.Println("SetHorn")
	defer log.Println("SetHorn-Done")
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_HORN)
	var soundHorn int
	if enabled {
		soundHorn = 1
	} else {
		soundHorn = 0
	}

	cmdArray[1] = byte(soundHorn)
	err := a.sendCommand(cmdArray)
	return err
}

func (a *TrainEngine) SetReverse(enabled bool) error {
	log.Println("SetReverse")
	defer log.Println("SetReverse-Done")
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_REVERSE)
	var soundHorn int
	if enabled {
		soundHorn = 1
	} else {
		soundHorn = 0
	}

	cmdArray[1] = byte(soundHorn)
	err := a.sendCommand(cmdArray)
	(*a).state.Reverse = enabled
	return err
}

func (a *TrainEngine) GetReverse() bool {
	return (*a).state.Reverse
}

func (a *TrainEngine) SetBell(enabled bool) error {
	log.Println("SetBell")
	defer log.Println("SetBell-Done")
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_BELL)
	var soundHorn int
	if enabled {
		soundHorn = 1
	} else {
		soundHorn = 0
	}

	cmdArray[1] = byte(soundHorn)
	err := a.sendCommand(cmdArray)
	return err
}

func (a *TrainEngine) SetLight(enabled bool) error {
	log.Println("SetLight")
	defer log.Println("SetLight-Done")
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_LIGHTS)
	var soundHorn int
	if enabled {
		soundHorn = 1
	} else {
		soundHorn = 0
	}

	cmdArray[1] = byte(soundHorn)
	err := a.sendCommand(cmdArray)
	(*a).state.Light = enabled
	return err
}

func (a *TrainEngine) GetLight() bool {
	return (*a).state.Light
}

func (a *TrainEngine) Speak() error {
	cmdArray := make([]byte, 2)
	cmdArray[0] = byte(COMMANDTYPE_SPEAK)
	cmdArray[1] = byte(0)
	err := a.sendCommand(cmdArray)
	return err
}

func (a *TrainEngine) SendCustomCommand(cmd []byte) error {
	return a.sendCommand(cmd)
}
