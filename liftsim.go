package lift

import (
	"errors"
	"io"
	"sync"
	"time"
)

const (
	floorHight       = 300 // [cm] distance between floors.
	inFloorThreshold = 20  // [cm] distance above and below a floor that count as in floor.
	endSensorOffset  = 20  // [cm] distance from first and last floor to end sensors.
)

// Lift defines what a lift can do.
type Lift interface {
	SetMotorDirection(direction int)
	OrderButtonLight(btn int, floor int, on bool)
	FloorIndicator(floor int)
	DoorLight(on bool)
	StopLight(on bool)

	OrderButton(button int, floor int) bool
	FloorSensor() (inFloor bool, floor int)
	StopButton() bool
	Obstruction() bool
}

// LiftSim represent one simulated lift.
type LiftSim struct {
	mu sync.RWMutex
	// Cart state
	floor      int     // 0-indexed => Ground floor = Floor 0
	liftHeight float32 // [cm] current hight of lift
	direction  int     // -1=down, 0=stopped, 1=up
	belowFirst bool
	aboveLast  bool

	// Light state
	stopLight bool
	doorLight bool

	// Lift specifications
	floors      int     // number of floors
	shaftHeight float32 // [cm]
	speed       float32 // [cm/s] movement speed of lift

	// Communication
	errors chan<- error
	name   string
	output io.Writer
}

// NewLiftSim returns a new simulated lift.
func NewLiftSim(err chan<- error, options ...Option) (*LiftSim, error) {
	sim := LiftSim{
		errors:      err,
		floors:      4,
		shaftHeight: floorHight * 4,
		speed:       50,
	}

	for _, opt := range options {
		if err := opt(&sim); err != nil {
			return nil, err
		}
	}

	return &sim, nil
}

// Run starts the lift.
func (ls *LiftSim) Run() {
	go func() {
		var scaling = 100
		ticker := time.NewTicker(time.Duration(1000/scaling) * time.Millisecond)
		for range ticker.C {
			ls.mu.Lock()
			ls.liftHeight = ls.liftHeight + float32(ls.direction)*ls.speed/float32(scaling)
			ls.mu.Unlock()

			ls.mu.RLock()
			if ls.liftHeight > ls.shaftHeight+10.0 {
				ls.errors <- errors.New("lift hit end sensor above top floor")
			} else if ls.liftHeight < -10.0 {
				ls.errors <- errors.New("lift hit end sensor below ground floor")
			}
			ls.mu.RUnlock()

			//ls.writeState()
		}
	}()
}

// SetMotorDirection sets the direction of travel for the lift. It accept only a integer with the convention
//  -1 = down, 0 = stop, 1 = up. If any other integer is supplied it will do nothing.
func (ls *LiftSim) SetMotorDirection(direction int) {
	invalid := direction != 0 && direction != -1 && direction != 1
	if invalid {
		return
	}

	ls.mu.Lock()
	ls.direction = direction
	ls.mu.Unlock()
}

// OrderButtonLight sets the lift for the specified button.
func (ls *LiftSim) OrderButtonLight(btn int, floor int, on bool) {
	panic("not implemented")
}

// FloorIndicator sets the floor indicator to the specified floor.
func (ls *LiftSim) FloorIndicator(floor int) {
	panic("not implemented")
}

// DoorLight toggles the door open light.
func (ls *LiftSim) DoorLight(on bool) {
	panic("not implemented")
}

// StopLight toggles the stop light.
func (ls *LiftSim) StopLight(on bool) {
	panic("not implemented")
}

// OrderButton polls if anyone is pressing the order button specified.
func (ls *LiftSim) OrderButton(button int, floor int) bool {
	panic("not implemented")
}

// FloorSensor polls the floor sensors. If the lift isn't in any floor at the moment inFloor returns false.
func (ls *LiftSim) FloorSensor() (inFloor bool, floor int) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	distanceFromFloor := int64(ls.liftHeight) % floorHight
	sensorTriggered := distanceFromFloor < inFloorThreshold || distanceFromFloor > floorHight-inFloorThreshold
	if !sensorTriggered {
		return false, 0
	}
	return true, int(ls.liftHeight / floorHight)

}

// StopButton polls of anyone is pressing the stop button.
func (ls *LiftSim) StopButton() bool {
	panic("not implemented")
}

// Obstruction polls if the obstruction switch is enabled.
func (ls *LiftSim) Obstruction() bool {
	panic("not implemented")
}
