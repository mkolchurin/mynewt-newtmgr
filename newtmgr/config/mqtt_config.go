package config

import (
	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/nmxact/nmqtt"
	"strconv"
	"strings"
)

// connstring = "DeviceID, BrokerAddr, UserName, Password, QoS, ClientID".
// For example:
// connstring="Client1,tcp://localhost:1883,,,0,0" or
// connstring="Client1,tcp://localhost:1883,User,Password,0,1"
func ParseMqttConnString(cs string) (*nmqtt.MqttXPortCfg, error) {
	sc := nmqtt.NewXportCfg()

	parts := strings.Split(cs, ",")

	sc.Id = parts[0]
	sc.Broker = parts[1]
	sc.User = parts[2]
	sc.Password = parts[3]
	qos, _ := strconv.Atoi(parts[4])
	sc.Qos = int8(qos)
	sc.DeviceId, _ = strconv.Atoi(parts[5])

	return sc, nil
}

func BuildMqttXport(sc *nmqtt.MqttXPortCfg) (*nmqtt.MqttXPort, error) {
	sx := nmqtt.NewMqttXport(*sc)
	if err := sx.Start(); err != nil {
		return nil, util.ChildNewtError(err)
	}

	return sx, nil
}
