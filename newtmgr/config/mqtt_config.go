package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/nmxact/nmqtt"
	"strconv"
	"strings"
)

//ParseMqttConnString
//
//connstring = "DeviceID, BrokerAddr, UserName, Password, QoS, ClientID, MtuIn, MtuOut" (MtuIn and MtuOut is optional)
//
//For example:
//
//"path/to/topic/1,tcp://localhost:1883,,,0,0"
//
//"path/to/topic/2,tcp://localhost:1883,User,Password,0,1"
//
//"path/to/topic/3,tcp://localhost:1883,User,Password,0,1,1024"
//
//"path/to/topic/4,tcp://localhost:1883,User,Password,0,1,255,1024"
func ParseMqttConnString(cs string) (*nmqtt.MqttXPortCfg, error) {
	sc := nmqtt.NewXportCfg()

	parts := strings.Split(cs, ",")
	strlen := len(parts)
	if strlen < 6 || strlen > 8 {
		return nil, errors.New("wrong args format")
	}
	sc.Id = parts[0]
	sc.Broker = parts[1]
	sc.User = parts[2]
	sc.Password = parts[3]
	qos, e := strconv.Atoi(parts[4])
	if e != nil {
		return nil, e
	}
	if qos != 0 && qos != 1 && qos != 2 {
		return nil, e
	}
	sc.Qos = int8(qos)
	sc.DeviceId, e = strconv.Atoi(parts[5])
	if e != nil {
		return nil, e
	}
	if strlen > 6 {
		sc.MtuIn, e = strconv.Atoi(parts[6])
		sc.MtuOut = sc.MtuIn
		if e != nil {
			return nil, e
		}
	}
	if strlen == 8 {
		sc.MtuOut, e = strconv.Atoi(parts[7])
		if e != nil {
			return nil, e
		}
	}
	log.Infof("mqtt id '%s'; Broker '%s'; user '%s'; qos '%d'; devId '%d'; MtuIn '%d'; MtuOut '%d'",
		sc.Id, sc.Broker, sc.User, sc.Qos, sc.DeviceId, sc.MtuIn, sc.MtuOut)
	return sc, nil
}

func BuildMqttXport(sc *nmqtt.MqttXPortCfg) (*nmqtt.MqttXPort, error) {
	sx := nmqtt.NewMqttXport(*sc)
	if err := sx.Start(); err != nil {
		return nil, util.ChildNewtError(err)
	}

	return sx, nil
}
