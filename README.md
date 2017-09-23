# lift
[![GoDoc](https://godoc.org/github.com/hdhauk/lift?status.svg)](https://godoc.org/github.com/hdhauk/lift)
[![Go Report Card](https://goreportcard.com/badge/github.com/hdhauk/lift)](https://goreportcard.com/report/github.com/hdhauk/lift)

Package lift provide an united interface for connection to either physical lifts or simulated lifts for use in the course TTK4145 at The Norwegian University of Science and Technology. This package is based on work done by Martin Korsgaard, Morten Fyhn and [klasbo](http://github.com/klasbo).

# Installation
```
go get -u github.com/hdhauk/lift
```

This package have not been tested with physical hardware, but the simulator part are tested on Ubuntu 17.04. To be able to compile you must have the Comedi drivers which can be downloaded by:
```
wget http://www.comedi.org/download/comedilib-0.10.2.tar.gz
tar -xvzf comedilib-0.10.2.tar.gz
cd comedilib-0.10.2 && ./configure && make && sudo make install
```

Furthermore is the simulator by github.com/klasbo statically embedded in the project using the [gobuffalo/packr](http://github.com/gobuffalo/packr) tool. This mean that to generate standalone binaries first do
```
go get -u github.com/gobuffalo/packr/... 	# Installation, only needed once
packr 						# Generate static files that contain binaries and config.
```

## Example
```GO
package main

import "github.com/hdhauk/lift"

func main{
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
```

## Documentation
Please see: [GoDoc](https://godoc.org/github.com/hdhauk/lift)
