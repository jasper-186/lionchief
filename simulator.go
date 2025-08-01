package lionchief

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"tinygo.org/x/bluetooth"
)

type TrainSimulator struct {
	address bluetooth.Address
	engine  *TrainEngine
}

func NewSimulator(trainAddress bluetooth.Address) (*TrainSimulator, error) {
	train, err := NewEngineDefaultBluetoothAdapter(trainAddress)
	if err != nil {
		return nil, err
	}

	simulator := TrainSimulator{
		address: trainAddress,
		engine:  train,
	}

	return &simulator, nil
}

func (a *TrainSimulator) Disconnect() error {
	return a.engine.Disconnect()
}

func (a *TrainSimulator) Reconnect() error {
	train, err := NewEngineDefaultBluetoothAdapter(a.address)
	if err != nil {
		return err
	}
	a.engine = train
	return nil
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
		err := a.engine.SetEngineVolume(newVolume)
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

func (a *TrainSimulator) SoundHorn(length int) error {
	err := a.engine.SetHorn(true)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(length))
	err = a.engine.SetHorn(false)
	if err != nil {
		return err
	}
	return nil
}

func (a *TrainSimulator) SoundBell(length int) error {
	err := a.engine.SetBell(true)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(length))
	err = a.engine.SetBell(false)
	if err != nil {
		return err
	}
	return nil
}

func (a *TrainSimulator) Speak() error {
	validPhrases := []int{SPEECHPHRASE_CALL_ME_PENNSYLVANIA_FLYER, SPEECHPHRASE_FASTEST_FREIGHT_YOU_CAN_HIRE, SPEECHPHRASE_HEY_THERE_WHAT_ARE_YOU_WAITING_FOR, SPEECHPHRASE_I_MAKE_STEAM_FROM_WATER_AND_FIRE, SPEECHPHRASE_PENNSYLVANIA_FLYER_IS_READY_TO_ROLL, SPEECHPHRASE_IM_FEELING_A_LITTLE_SQUEAKY_GIVE_ME_A_LITTLE_OIL}
	phrase := validPhrases[rand.Intn(len(validPhrases))]
	return a.engine.SpeakPhrase(SpeechPhrase(phrase))
}

func (a *TrainSimulator) SpeakPhrase(phrase SpeechPhrase) error {
	return a.engine.SpeakPhrase(phrase)
}

func (a *TrainSimulator) SpeakSpeel() error {

	for i := 4; i < 7; i++ {
		a.engine.SpeakPhrase(SpeechPhrase(i))
		time.Sleep(3 * time.Second)
	}
	return nil
}

func (a *TrainSimulator) Lights(enabled bool) error {
	log.Println("Lights")
	err := a.engine.SetLight(enabled)
	log.Println("Lights-Done")
	return err
}

func (a *TrainSimulator) GetCurrentState() *TrainState {
	return a.engine.state
}

func (a *TrainSimulator) ToggleLights() error {
	log.Println("ToggleLights")
	return a.engine.SetLight(!a.engine.GetLight())
}

func (a *TrainSimulator) SetMainVolume(volume int) error {
	return a.engine.SetMainVolume(volume)
}

func (a *TrainSimulator) SetBellVolume(volume int) error {
	return a.engine.SetBellVolume(volume)
}

func (a *TrainSimulator) SetEngineVolume(volume int) error {
	return a.engine.SetEngineVolume(volume)
}

func (a *TrainSimulator) SetHornVolume(volume int) error {
	return a.engine.SetHornVolume(volume)
}

func (a *TrainSimulator) SetSpeechVolume(volume int) error {
	return a.engine.SetSpeechVolume(volume)
}

func (a *TrainSimulator) SetBellPitch(pitch SoundPitch) error {
	return a.engine.SetBellPitch(pitch)
}

func (a *TrainSimulator) SetEnginePitch(pitch SoundPitch) error {
	return a.engine.SetEnginePitch(pitch)
}

func (a *TrainSimulator) SetHornPitch(pitch SoundPitch) error {
	return a.engine.SetHornPitch(pitch)
}

func (a *TrainSimulator) SetSpeechPitch(pitch SoundPitch) error {
	return a.engine.SetSpeechPitch(pitch)
}
