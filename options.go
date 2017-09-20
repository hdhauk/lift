package lift

import (
	"errors"
	"io"
)

// Option can configure the settings when spawning a new LiftSim.
type Option func(*LiftSim) error

// Floors sets the number of floors, and where the lift start.
func Floors(start, floors int) Option {
	return func(ls *LiftSim) error {
		if floors < 2 {
			return errors.New("cannot have less than 2 floors")
		}
		if start < 0 || start > floors-1 {
			return errors.New("invalid start floor")
		}
		ls.floors = floors
		ls.floor = start
		ls.shaftHeight = float32(floors * floorHight)
		ls.liftHeight = float32(start * floorHight)
		return nil
	}
}

// Speed sets the speed of the lift in cm/second.
func Speed(speed float32) Option {
	return func(ls *LiftSim) error {
		ls.speed = speed
		return nil
	}
}

// WriteTo sets the output of the text based GUI.
func WriteTo(w io.Writer) Option {
	return func(ls *LiftSim) error {
		ls.output = w
		return nil
	}
}
