package lift

import (
	"errors"
	"fmt"
	"time"
)

// SimOption are option functions for simulated lifts.
type SimOption func(*liftSim) error

// NumFloors sets the number of floors in simulated lift.
func NumFloors(floors int) SimOption {
	return func(l *liftSim) error {
		if floors < 2 || floors > 9 {
			return errors.New("number of floors must be between 2 and 9")
		}
		l.NumFloors = floors
		return nil
	}
}

// TravelTimeBetweenFloors sets the travel time between each floor.
func TravelTimeBetweenFloors(d time.Duration) SimOption {
	return func(l *liftSim) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.TravelTimeBetweenFloors = int(ms)
		return nil
	}
}

// TravelTimePassingFloors sets the time spent within the sensors range when passing a floor.
func TravelTimePassingFloors(d time.Duration) SimOption {
	return func(l *liftSim) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.TravelTimePassingFloors = int(ms)
		return nil
	}
}

// BtnDepressedTime sets the duration a button will be considered pressed after the actual
// key event happens.
func BtnDepressedTime(d time.Duration) SimOption {
	return func(l *liftSim) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.BtnDepressedTime = int(ms)
		return nil
	}
}

// ComPort sets what port the simulator should listen on. This must be unique per
// instance of simulators. Must be in the range 1024 - 65535.
func ComPort(port int) SimOption {
	return func(l *liftSim) error {
		if port < 1024 || port > 65535 {
			return fmt.Errorf("illegal port: %d", port)
		}
		l.ComPort = port
		return nil
	}
}
