package lift_test

import "github.com/hdhauk/lift"

// Basic usage.
func Example() {
	// The order to which options or how many options you supply doesn't matter.
	// Here we configure the simulator to have 6 floors, and to listen on port 9999.
	// If neither is supplied the defaults will be used (4 floors and a random port).
	sim, err := lift.NewSim(lift.NumFloors(6), lift.ComPort(9999))
	if err != nil {
		panic(err)
	}

	if err := sim.Init(); err != nil {
		panic(err)
	}

	// The simulator can then be controlled by using the methods defined by Lifter.
	// Example:
	sim.SetMotorDirection(1)
}
