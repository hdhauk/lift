package lift

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net"
	"os/exec"

	"github.com/gobuffalo/packr"
)

// Lift ...
type liftSim struct {
	simConfig
	txWithResp, txWithoutResp, rx chan []byte
}
type simConfig struct {
	TravelTimeBetweenFloors int
	TravelTimePassingFloors int
	BtnDepressedTime        int

	NumFloors int

	ComPort int
}

// NewSim returns a new simulated lift.
func NewSim(options ...SimOption) (Lifter, error) {
	sim := liftSim{
		simConfig: simConfig{
			TravelTimeBetweenFloors: 2000,
			TravelTimePassingFloors: 500,
			BtnDepressedTime:        200,
			NumFloors:               4,
			ComPort:                 1024 + rand.Intn(65535-1024), // Pick random port
		},
		txWithResp:    make(chan []byte),
		txWithoutResp: make(chan []byte),
		rx:            make(chan []byte),
	}

	for _, opt := range options {
		if err := opt(&sim); err != nil {
			return nil, err
		}
	}

	return &sim, nil
}

// Init initializes the lift for use. For a simulated lift this will spawn a new
// gnome-terminal window with the simulator in it.
func (ls *liftSim) Init() error {
	if err := spawnSimulator(ls.simConfig); err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		endpoint := fmt.Sprintf("localhost:%d", ls.ComPort)
		conn, err := net.Dial("tcp", endpoint)
		if err != nil {
			errCh <- fmt.Errorf("unable to connect to simulator: %v", err)
		}
		defer conn.Close()
		errCh <- nil

		for {
			select {
			case cmd := <-ls.txWithResp:
				conn.Write(cmd)
				resp := make([]byte, 4)
				conn.Read(resp)
				ls.rx <- resp
			case cmd := <-ls.txWithoutResp:
				conn.Write(cmd)
			}
		}
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

// SetMotorDirection sets the direction of travel for the lift. It accept only a integer with the convention
//  -1 = down, 0 = stop, 1 = up. If any other values will do nothing.
func (ls *liftSim) SetMotorDirection(direction int) {
	switch direction {
	case 1:
		ls.txWithoutResp <- []byte("\x01\x01\x00\x00")
	case 0:
		ls.txWithoutResp <- []byte("\x01\x00\x00\x00")
	case -1:
		ls.txWithoutResp <- []byte("\x01\xFF\x00\x00")
	default:
		return
	}
}

// OrderButtonLight sets the lift for the specified button.
func (ls *liftSim) OrderButtonLight(button int, floor int, on bool) {
	ls.txWithoutResp <- []byte{2, byte(button), byte(floor), bool2Byte(on)}
}

// FloorIndicator sets the floor indicator to the specified floor.
func (ls *liftSim) FloorIndicator(floor int) {
	if floor < 0 || floor > ls.NumFloors-1 {
		return
	}
	ls.txWithoutResp <- []byte{3, byte(floor), 0, 0}
}

// DoorLight toggles the door open light.
func (ls *liftSim) DoorLight(on bool) {
	ls.txWithoutResp <- []byte{4, bool2Byte(on), 0, 0}
}

// StopLight toggles the stop light.
func (ls *liftSim) StopLight(on bool) {
	ls.txWithoutResp <- []byte{5, bool2Byte(on), 0, 0}
}

// OrderButton polls if anyone is pressing the order button specified.
func (ls *liftSim) OrderButton(button int, floor int) bool {
	ls.txWithResp <- []byte{6, byte(button), byte(floor), 0}
	resp := <-ls.rx
	return resp[1] == 1
}

// FloorSensor polls the floor sensors. If the lift isn't in any floor at the moment inFloor returns false.
func (ls *liftSim) FloorSensor() (inFloor bool, floor int) {
	ls.txWithResp <- []byte{7, 0, 0, 0}
	resp := <-ls.rx
	return resp[1] != 0, int(resp[2])
}

// StopButton polls of anyone is pressing the stop button.
func (ls *liftSim) StopButton() bool {
	ls.txWithResp <- []byte{8, 0, 0, 0}
	resp := <-ls.rx
	return resp[1] == 1
}

// Obstruction polls if the obstruction switch is enabled.
func (ls *liftSim) Obstruction() bool {
	ls.txWithResp <- []byte{9, 0, 0, 0}
	resp := <-ls.rx
	return resp[1] == 1
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

func bool2Byte(b bool) byte {
	if b {
		return byte(1)
	}
	return byte(0)
}
