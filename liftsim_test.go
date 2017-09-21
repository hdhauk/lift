package lift

import (
	"testing"
)

var testErrorCh = make(chan error)

// func TestNewLiftSim(t *testing.T) {
// 	t.Run("invalid_floors", func(t *testing.T) {
// 		numFloors, startFloor := 3, 3
// 		_, err := NewLift(testErrorCh, Floors(startFloor, numFloors))
// 		if err == nil {
// 			t.Errorf("Created lift with %d floors, and starting in floor %d (zero-indexed). Expected error, but got nil.", numFloors, startFloor)
// 		}
// 	})

// 	t.Run("check_correct_initial_floor", func(t *testing.T) {
// 		startFloor := 2
// 		sim, err := NewLiftSim(testErrorCh, Floors(startFloor, 3))
// 		if err != nil {
// 			t.Fatalf("%v", err)
// 		}

// 		inFloor, floor := sim.FloorSensor()
// 		if !inFloor {
// 			t.Error("Expected lift to be in a floor.")
// 		}
// 		if startFloor != floor {
// 			t.Errorf("Expected lift to be in floor %d, but found it in floor %d", startFloor, floor)
// 		}
// 	})

// }

// func TestTODO(t *testing.T) {
// 	sim, err := NewLiftSim(testErrorCh, Floors(2, 3), Speed(200))
// 	if err != nil {
// 		t.Fatalf("%v", err)
// 	}

// 	sim.direction = -1
// 	// sim.Run()

// 	// err = <-testErrorCh
// 	// t.Fatalf("%v", err)
// }

func Test_spawnSimulator(t *testing.T) {
	c := simConfig{
		TravelTimeBetweenFloors: 500,
		TravelTimePassingFloors: 200,
		BtnDepressedTime:        50,
		NumFloors:               6,
		ComPort:                 5000,
	}
	if err := spawnSimulator(c); err != nil {
		t.Fatal(err)
	}

}

// func Test_createConfigBytes(t *testing.T) {
// 	c := simConfig{
// 		TravelTimeBetweenFloors: 500,
// 		TravelTimePassingFloors: 200,
// 		BtnDepressedTime:        50,
// 		NumFloors:               4,
// 		ComPort:                 5000,
// 	}
// 	ret, err := createConfigBytes(c)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	println(string(ret))
// }
