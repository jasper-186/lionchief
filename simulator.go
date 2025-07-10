package lionchief

import (
	"fmt"
	"log"
	"math"
	"time"

	"tinygo.org/x/bluetooth"
)

type TrainSimulator struct {
	engine *TrainEngine
}

func NewSimulator(trainAddress bluetooth.Address) (*TrainSimulator, error) {
	train, err := NewEngineDefaultBluetoothAdapter(trainAddress)
	if err != nil {
		return nil, err
	}

	simulator := TrainSimulator{
		engine: train,
	}

	return &simulator, nil
}

func (a *TrainSimulator) Disconnect() error {
	return a.engine.Disconnect()
}

func (a *TrainSimulator) AdjustSpeedTo(speed int) error {
	if speed < 0 || 31 < speed {
		return fmt.Errorf("speed must be between 0 and 31")
	}

	initialSpeed := a.engine.GetSpeed()
	var increment int
	if initialSpeed == speed {
		return fmt.Errorf("train is already at speed %d", speed)
	} else if initialSpeed > speed {
		increment = -1
	} else {
		increment = 1
	}

	currentSpeed := initialSpeed
	for currentSpeed != speed {
		newSpeed := currentSpeed + increment
		newVolume := int(math.Ceil(float64(newSpeed) / 3))
		err := a.engine.SetRunningVolume(SOUNDTYPE_ENGINE, newVolume)
		if err != nil {
			return err
		}

		err = a.engine.SetSpeed(newSpeed)
		if err != nil {
			return err
		}

		currentSpeed = newSpeed
	}
	return nil
}

func (a *TrainSimulator) BeginTrainService() error {
	err := a.engine.SetBell(true)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.engine.SetBell(false)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	err = a.AdjustSpeedTo(3)
	if err != nil {
		return err
	}
	return nil
}

func (a *TrainSimulator) EndTrainService() error {
	err := a.engine.SetHorn(true)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.engine.SetHorn(false)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.AdjustSpeedTo(0)
	if err != nil {
		return err
	}
	return nil
}

func (a *TrainSimulator) ReverseTrainService() error {
	err := a.engine.SetHorn(true)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.engine.SetHorn(false)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	err = a.engine.SetHorn(true)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.engine.SetHorn(false)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	originalSpeed := a.engine.GetSpeed()
	err = a.AdjustSpeedTo(0)
	if err != nil {
		return err
	}

	err = a.engine.SetReverse(!a.engine.GetReverse())
	if err != nil {
		return err
	}

	err = a.AdjustSpeedTo(originalSpeed)
	if err != nil {
		return err
	}

	return nil
}

func (a *TrainSimulator) SoundHorn(length int) {
	a.engine.SetHorn(true)
	time.Sleep(time.Second * time.Duration(length))
	a.engine.SetHorn(false)
}

func (a *TrainSimulator) SoundBell(length int) {
	a.engine.SetBell(true)
	time.Sleep(time.Second * time.Duration(length))
	a.engine.SetBell(false)
}

func (a *TrainSimulator) Speak() error {
	return a.engine.Speak()
}

func (a *TrainSimulator) Lights(enabled bool) error {
	log.Println("Lights")
	err := a.engine.SetLight(enabled)
	log.Println("Lights-Done")
	return err
}

func (a *TrainEngine) GetCurrentState() *TrainState {
	return a.state
}

func (a *TrainSimulator) ToggleLights() error {
	log.Println("ToggleLights")
	return a.engine.SetLight(!a.engine.GetLight())
}
