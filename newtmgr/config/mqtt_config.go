package config

import (
	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/nmxact/nmqtt"
	"strconv"
	"strings"
)

func ParseMqttConnString(cs string) (*nmqtt.MqttXPortCfg, error) {
	sc := nmqtt.NewXportCfg()

	parts := strings.Split(cs, ",")

	sc.Id = parts[0]

	sc.Broker = parts[1]
	sc.User = parts[2]
	sc.Password = parts[3]
	sc.DeviceId, _ = strconv.Atoi(parts[4])
	i, _ := strconv.Atoi(parts[5])
	sc.Qos = int8(i)

	return sc, nil
}

func BuildMqttXport(sc *nmqtt.MqttXPortCfg) (*nmqtt.MqttXPort, error) {
	sx := nmqtt.NewMqttXport(*sc)
	if err := sx.Start(); err != nil {
		return nil, util.ChildNewtError(err)
	}

	return sx, nil
}
