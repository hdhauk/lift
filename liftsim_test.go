package lift

import (
	"fmt"
	"testing"
	"time"
)

func Test_TxOnlyCommands(t *testing.T) {
	sim, err := NewSim(NumFloors(6))
	if err != nil {
		t.Fatal(err)
	}

	if err := sim.Init(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	sim.OrderButtonLight(0, 0, true)
	sim.OrderButtonLight(1, 1, true)
	sim.OrderButtonLight(2, 2, true)

	sim.FloorIndicator(5)
	sim.SetMotorDirection(1)
	sim.DoorLight(true)
	sim.StopLight(false)
	time.Sleep(1 * time.Second)

	sim.SetMotorDirection(0)
	sim.FloorIndicator(1)
	sim.DoorLight(false)
	sim.StopLight(true)
	time.Sleep(1 * time.Second)

	sim.SetMotorDirection(-1)
	sim.FloorIndicator(2)
	sim.DoorLight(true)
	sim.StopLight(false)
	time.Sleep(1 * time.Second)

	sim.SetMotorDirection(0)
	sim.FloorIndicator(7)
	time.Sleep(1 * time.Second)
}

func Test_TxRx(t *testing.T) {
	sim, err := NewSim(NumFloors(9))
	if err != nil {
		t.Fatal(err)
	}

	if err := sim.Init(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	fmt.Printf("%v\n", sim.OrderButton(2, 2))

}
