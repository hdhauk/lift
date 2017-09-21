package lift

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os/exec"

	"github.com/gobuffalo/packr"
)

const (
	floorHight       = 300 // [cm] distance between floors.
	inFloorThreshold = 20  // [cm] distance above and below a floor that count as in floor.
	endSensorOffset  = 20  // [cm] distance from first and last floor to end sensors.
)

// Lifter defines what a lift can do.
type Lifter interface {
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

// Lift ...
type Lift struct {
	sim bool
}

// NewLift returns a new simulated lift.
func NewLift(err chan<- error, options ...Option) (*Lift, error) {
	lift := Lift{}

	for _, opt := range options {
		if err := opt(&lift); err != nil {
			return nil, err
		}
	}

	return &lift, nil
}

// Run starts the lift.
func (ls *Lift) Run() {
	panic("not implemented")
}

// SetMotorDirection sets the direction of travel for the lift. It accept only a integer with the convention
//  -1 = down, 0 = stop, 1 = up. If any other integer is supplied it will do nothing.
func (ls *Lift) SetMotorDirection(direction int) {
	panic("not implemented")

}

// OrderButtonLight sets the lift for the specified button.
func (ls *Lift) OrderButtonLight(btn int, floor int, on bool) {
	panic("not implemented")
}

// FloorIndicator sets the floor indicator to the specified floor.
func (ls *Lift) FloorIndicator(floor int) {
	panic("not implemented")
}

// DoorLight toggles the door open light.
func (ls *Lift) DoorLight(on bool) {
	panic("not implemented")
}

// StopLight toggles the stop light.
func (ls *Lift) StopLight(on bool) {
	panic("not implemented")
}

// OrderButton polls if anyone is pressing the order button specified.
func (ls *Lift) OrderButton(button int, floor int) bool {
	panic("not implemented")
}

// FloorSensor polls the floor sensors. If the lift isn't in any floor at the moment inFloor returns false.
func (ls *Lift) FloorSensor() (inFloor bool, floor int) {
	panic("not implemented")
}

// StopButton polls of anyone is pressing the stop button.
func (ls *Lift) StopButton() bool {
	panic("not implemented")
}

// Obstruction polls if the obstruction switch is enabled.
func (ls *Lift) Obstruction() bool {
	panic("not implemented")
}

type simConfig struct {
	TravelTimeBetweenFloors int
	TravelTimePassingFloors int
	BtnDepressedTime        int

	NumFloors int

	ComPort int
}

func spawnSimulator(config simConfig) error {
	// Unbox executable binary.
	box := packr.NewBox("./bin")
	binData := box.Bytes("sim_server")

	// Create temporary file.
	tmpExec, err := ioutil.TempFile("", "sim_server")
	if err != nil {
		return fmt.Errorf("unable to create temp executable: %v", err)
	}

	// Make temporary executable.
	if err := tmpExec.Chmod(0777); err != nil {
		return fmt.Errorf("unable to let temp executable actually be executable: %v", err)
	}
	defer tmpExec.Close()

	// Write sim_server binary data to temporary file.
	if _, err := tmpExec.Write(binData); err != nil {
		return fmt.Errorf("unable to write binary data to temp file: %d", err)
	}

	// Make temporary config file.
	tmpConfig, err := ioutil.TempFile("", "sim_config")
	if err != nil {
		return fmt.Errorf("unable to create temp config: %v", err)
	}
	defer tmpConfig.Close()

	// Fill config file
	b, err := createConfigBytes(config)
	if err != nil {
		return err
	}
	tmpConfig.Write(b)

	// Close files.
	tmpExec.Close()
	tmpConfig.Close()

	// Open gnome-terminal and run simulator inside it.
	c := exec.Command("gnome-terminal", "--working-directory=/tmp", "-e", tmpExec.Name()+"	"+tmpConfig.Name())
	if _, err := c.Output(); err != nil {
		return fmt.Errorf("failed to run simulator in gnome-terminal: %v", err)
	}

	return nil
}

func createConfigBytes(config simConfig) ([]byte, error) {
	box := packr.NewBox("./bin")
	configTemplate := box.String("config.con")
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	buf := bytes.NewBuffer(nil)

	if err := tmpl.Execute(buf, config); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.Bytes(), nil
}
