package main

import (
	"fmt"
	"os"

	"mynewt.apache.org/newtmgr/nmxact/nmqtt"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
	"mynewt.apache.org/newtmgr/nmxact/xact"
)

// "mynewt.apache.org/newtmgr/nmxact/nmserial"
// "mynewt.apache.org/newtmgr/nmxact/mqtt"

func main() {
	// Initialize the serial transport.

	cfg := nmqtt.NewXportCfg()
	cfg.Id = "12345"
	cfg.Broker = "tcp://localhost:1883"
	cfg.Password = ""
	cfg.User = ""
	cfg.DeviceId = 0

	xport := nmqtt.NewMqttXport(*cfg)
	if err := xport.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting serial transport: %s\n",
			err.Error())
		os.Exit(1)
	}
	defer xport.Stop()

	// Create and open a session for connected Mynewt device.
	sc := sesn.NewSesnCfg()
	sc.MgmtProto = sesn.MGMT_PROTO_NMP

	s, err := xport.BuildSesn(sc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating serial session: %s\n",
			err.Error())
		os.Exit(1)
	}

	if err := s.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting serial session: %s\n",
			err.Error())
		os.Exit(1)
	}
	defer s.Close()

	// Send an echo command to the device.
	//c := xact.NewEchoCmd()
	//c.Payload = "hello"
	c := xact.NewImageStateReadCmd()
	res, err := c.Run(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing echo command: %s\n",
			err.Error())
		os.Exit(1)
	}

	if res.Status() != 0 {
		fmt.Printf("Device responded negatively to echo command; status=%d\n",
			res.Status())
	}

	eres := res.(*xact.ImageStateReadResult)
	fmt.Printf("Device send back: %s\n", eres.Rsp.Images[0].Hash)
}
