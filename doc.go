// Package lift provide an united interface for connection to either
// physical lifts or simulated lifts for use in the course TTK4145 at
// The Norwegian University of Science and Technology. This package is
// based on work done by Martin Korsgaard, Morten Fyhn and github.com/klasbo.
//
// This package have not been tested with physical hardware, but the simulator part
// are tested on Ubuntu 17.04. To be able to compile you must have the Comedi drivers
// which can be downloaded by:
//	wget http://www.comedi.org/download/comedilib-0.10.2.tar.gz
//	tar -xvzf comedilib-0.10.2.tar.gz
//	cd comedilib-0.10.2 && ./configure && make && sudo make install
//
// Furthermore is the simulator by github.com/klasbo statically embedded in the project using the
// http://github.com/gobuffalo/packr tool. This mean that to generate standalone binaries first do
//	go get -u github.com/gobuffalo/packr/... 	# Installation, only needed once
// 	packr 						# Generate static files that contain binaries and config.
package lift
