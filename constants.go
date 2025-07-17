package lionchief

import "tinygo.org/x/bluetooth"

type CommandType int

var (
	ReadWriteService, _                = bluetooth.ParseUUID("e20a39f4-73f5-4bc4-a12f-17d1ad07a961")
	WriteCharacteristic, _             = bluetooth.ParseUUID("08590f7e-db05-467e-8757-72f6faeb13d4")
	ReadCharateristic, _               = bluetooth.ParseUUID("08590f7e-db05-467e-8757-72f6faeb14d3")
	SystemId, _                        = bluetooth.ParseUUID("00002a23-0000-1000-8000-00805f9b34fb")
	DeviceName, _                      = bluetooth.ParseUUID("00002a00-0000-1000-8000-00805f9b34fb") // 'LC-0-1-0429-754D'
	ModelNumber, _                     = bluetooth.ParseUUID("00002a24-0000-1000-8000-00805f9b34fb")
	SerialNumber, _                    = bluetooth.ParseUUID("00002a25-0000-1000-8000-00805f9b34fb")
	FirmwareRevision, _                = bluetooth.ParseUUID("00002a26-0000-1000-8000-00805f9b34fb")
	HardwareRevision, _                = bluetooth.ParseUUID("00002a27-0000-1000-8000-00805f9b34fb")
	SoftwareRevision, _                = bluetooth.ParseUUID("00002a28-0000-1000-8000-00805f9b34fb")
	ManufacturerName, _                = bluetooth.ParseUUID("00002a29-0000-1000-8000-00805f9b34fb")
	RegulatoryCertificationDataList, _ = bluetooth.ParseUUID("00002a2a-0000-1000-8000-00805f9b34fb")
	PnpId, _                           = bluetooth.ParseUUID("00002a50-0000-1000-8000-00805f9b34fb")
)

const (
	// Low Level command ids
	COMMANDTYPE_SOUND_RUNNING = 68
	COMMANDTYPE_SPEED         = 69
	COMMANDTYPE_REVERSE       = 70
	COMMANDTYPE_BELL          = 71
	COMMANDTYPE_HORN          = 72
	COMMANDTYPE_DISCONNECT    = 75
	COMMANDTYPE_SOUND_MAIN    = 76
	COMMANDTYPE_SPEAK         = 77
	COMMANDTYPE_LIGHTS        = 81
)

type SoundType int

const (
	SOUNDTYPE_BELL   = 2
	SOUNDTYPE_ENGINE = 4
	SOUNDTYPE_HORN   = 1
	SOUNDTYPE_SPEECH = 3
)

type SoundPitch int

const (
	SOUNDPITCH_LOWEST  int = 254
	SOUNDPITCH_LOW     int = 255
	SOUNDPITCH_NORMAL  int = 0
	SOUNDPITCH_HIGH    int = 1
	SOUNDPITCH_HIGHEST int = 2
)
