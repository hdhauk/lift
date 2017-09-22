package lift

// Lifter is the interface for what the lifts can do. It offers few guarantees
// and will generally simply do nothing if used wrong.
// The following conventions are used:
//	Direction:	-1 = down, 0 = stop, 1 = up
// 	Buttons:	0 = hall up, 1 = hall down, 2 = inside lift
type Lifter interface {
	Init() error

	SetMotorDirection(direction int)
	OrderButtonLight(button int, floor int, on bool)
	FloorIndicator(floor int)
	DoorLight(on bool)
	StopLight(on bool)

	OrderButton(button int, floor int) (pressed bool)
	FloorSensor() (inFloor bool, floor int)
	StopButton() (pressed bool)
	Obstruction() (active bool)
}

// New returns a new physical lift. It is assumed to only be one physical lift.
func New() (Lifter, error) {
	return &lift{}
}

type lift struct{}

func (l *lift) Init() error {
	return ioInit()
}

func (l *lift) SetMotorDirection(direction int) {
	switch direction {
	case 0:
		ioWriteAnalog(motor, 0)
	case 1:
		ioClearBit(motorDirDown)
		ioWriteAnalog(motor, 2800)
	case -1:
		ioSetBit(motorDirDown)
		ioWriteAnalog(motor, 2800)
	}
}

func (l *lift) OrderButtonLight(button int, floor int, on bool) {
	if on {
		ioSetBit(lampChannelMatrix[floor][button])
		return
	}
	ioClearBit(lampChannelMatrix[floor][button])
}

func (l *lift) FloorIndicator(floor int) {
	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(floorLED1)
	} else {
		ioClearBit(floorLED1)
	}

	if floor&0x01 > 0 {
		ioSetBit(floorLED2)
	} else {
		ioClearBit(floorLED2)
	}
}

func (l *lift) DoorLight(on bool) {
	if on {
		ioSetBit(doorOpenLED)
		return
	}
	ioClearBit(doorOpenLED)

}

func (l *lift) StopLight(on bool) {
	if on {
		ioSetBit(stopLED)
	}
	ioClearBit(stopLED)
}

func (l *lift) OrderButton(button int, floor int) bool {
	if ioReadBit(buttonChannelMatrix[floor][button]) {
		return true
	}
	return false
}

func (l *lift) FloorSensor() (inFloor bool, floor int) {
	if ioReadBit(sensorFloor1) {
		return true, 0
	} else if ioReadBit(sensorFloor2) {
		return true, 1
	} else if ioReadBit(sensorFloor3) {
		return true, 2
	} else if ioReadBit(sensorFloor4) {
		return true, 3
	} else {
		return false, -1
	}
}

func (l *lift) StopButton() bool {
	return ioReadBit(stopBtn)
}

func (l *lift) Obstruction() bool {
	return ioReadBit(obstruct)
}
