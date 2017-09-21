package lift

// Option can configure the settings when spawning a new LiftSim.
type Option func(*Lift) error

// SimFloors sets the number of floors in simulated lift.
func SimFloors(start, floors int) Option {
	return func(ls *Lift) error {

		return nil
	}
}

// Speed sets the speed of the lift in cm/second.
func Speed(speed float32) Option {
	return func(ls *Lift) error {

		return nil
	}
}
