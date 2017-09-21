package lift

import (
	"errors"
	"fmt"
	"time"
)

// Option can configure the settings when spawning a new LiftSim.
type Option func(*Lift) error

// SimNumFloors sets the number of floors in simulated lift.
func SimNumFloors(floors int) Option {
	return func(l *Lift) error {
		if floors < 2 || floors > 9 {
			return errors.New("number of floors must be between 2 and 9")
		}
		l.simConfig.NumFloors = floors
		return nil
	}
}

// SimTravelTimeBetweenFloors sets the travel time between each floor.
func SimTravelTimeBetweenFloors(d time.Duration) Option {
	return func(l *Lift) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.simConfig.TravelTimeBetweenFloors = int(ms)
		return nil
	}
}

// SimTravelTimePassingFloors sets the time spent within the sensors range when passing a floor.
func SimTravelTimePassingFloors(d time.Duration) Option {
	return func(l *Lift) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.simConfig.TravelTimePassingFloors = int(ms)
		return nil
	}
}

// SimBtnDepressedTime sets the duration a button will be considered pressed after the actual
// key event happens.
func SimBtnDepressedTime(d time.Duration) Option {
	return func(l *Lift) error {
		ms := d.Nanoseconds() * 1000 * 1000
		l.simConfig.BtnDepressedTime = int(ms)
		return nil
	}
}

// SimComPort sets what port the simulator should listen on. This must be unique per
// instance of simulators. Must be in the range 1024 - 65535;
func SimComPort(port int) Option {
	return func(l *Lift) error {
		if port < 1024 || port > 65535 {
			return fmt.Errorf("illegal port: %d", port)
		}
		l.simConfig.ComPort = port
		return nil
	}
}
