/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package nmqtt

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/runtimeco/go-coap"
	"mynewt.apache.org/newtmgr/nmxact/mgmt"
	"mynewt.apache.org/newtmgr/nmxact/nmcoap"
	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/omp"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
)

const MAX_PACKET_SIZE = 256

type MqttSesn struct {
	cfg  sesn.SesnCfg
	mx   *MqttXPort
	conn MQTT.Client
	txvr *mgmt.Transceiver
}

func NewMqttSesn(mx *MqttXPort, cfg sesn.SesnCfg) (*MqttSesn, error) {
	s := &MqttSesn{
		cfg: cfg,
		mx:  mx,
	}

	txvr, err := mgmt.NewTransceiver(cfg.TxFilter, cfg.RxFilter, true,
		cfg.MgmtProto, 3)
	if err != nil {
		return nil, err
	}
	s.txvr = txvr

	return s, nil
}

func (s *MqttSesn) Open() error {
	if s.conn != nil {
		return nmxutil.NewSesnAlreadyOpenError(
			"Attempt to open an already-open UDP session")
	}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(s.mx.MqttCfg.Broker)
	opts.SetClientID(s.mx.MqttCfg.Id + strconv.Itoa(rand.Intn(5000)))
	opts.SetUsername(s.mx.MqttCfg.User)
	opts.SetPassword(s.mx.MqttCfg.Password)
	opts.SetCleanSession(s.mx.MqttCfg.Cleansess)
	if s.mx.MqttCfg.Store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(s.mx.MqttCfg.Store))
	}

	s.conn = MQTT.NewClient(opts)
	if token := (s.conn).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	idStr := ""
	idStr = strconv.Itoa(s.mx.MqttCfg.DeviceId) + "/"
	if token := (s.conn).Subscribe((s.mx.MqttCfg.Id)+"/update/"+idStr+(s.mx.MqttCfg.RxTopic),
		byte(s.mx.MqttCfg.Qos), func(client MQTT.Client, message MQTT.Message) {
			if client.IsConnected() == true {
				s.txvr.DispatchNmpRsp(message.Payload())
			} else {
				fmt.Println("Connection is lost")
			}
		}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	log.Infof("Mqtt connected")
	return nil
}

func (s *MqttSesn) Close() error {
	if s.conn == nil {
		return nmxutil.NewSesnClosedError(
			"Attempt to close an unopened UDP session")
	}
	(s.conn).Disconnect(250)
	s.mx.started = false
	s.conn = nil
	s.mx = nil
	return nil
}

func (s *MqttSesn) IsOpen() bool {
	return s.conn != nil
}

func (s *MqttSesn) MtuIn() int {
	return MAX_PACKET_SIZE -
		omp.OMP_MSG_OVERHEAD -
		nmp.NMP_HDR_SIZE
}

func (s *MqttSesn) MtuOut() int {
	return MAX_PACKET_SIZE -
		omp.OMP_MSG_OVERHEAD -
		nmp.NMP_HDR_SIZE
}

func (s *MqttSesn) TxRxMgmt(m *nmp.NmpMsg,
	timeout time.Duration) (nmp.NmpRsp, error) {

	if !s.IsOpen() {
		return nil, fmt.Errorf("attempt to transmit over closed MQTT session")
	}

	txRaw := func(b []byte) error {
		idStr := strconv.Itoa(s.mx.MqttCfg.DeviceId) + "/"
		token := (s.conn).Publish((s.mx.MqttCfg.Id)+"/update/"+idStr+(s.mx.MqttCfg.TxTopic), byte(s.mx.MqttCfg.Qos), false, b)
		token.Wait()
		err := token.Error()
		if err != nil {
			log.Infof("token pub '%s'", err.Error())
		}
		return err
	}

	return s.txvr.TxRxMgmt(txRaw, m, s.MtuOut(), timeout)
}

func (s *MqttSesn) TxRxMgmtAsync(m *nmp.NmpMsg,
	timeout time.Duration, ch chan nmp.NmpRsp, errc chan error) error {
	rsp, err := s.TxRxMgmt(m, timeout)
	if err != nil {
		errc <- err
	} else {
		ch <- rsp
	}
	return nil
}

func (s *MqttSesn) AbortRx(seq uint8) error {
	s.txvr.ErrorAll(fmt.Errorf("rx aborted"))
	s.conn.Disconnect(100)
	return nil
}

func (s *MqttSesn) TxCoap(m coap.Message) error {
	txRaw := func(b []byte) error {
		idStr := strconv.Itoa(s.mx.MqttCfg.DeviceId) + "/"
		token := (s.conn).Publish((s.mx.MqttCfg.Id)+"/update/"+idStr+(s.mx.MqttCfg.TxTopic), byte(s.mx.MqttCfg.Qos), false, b)
		token.Wait()
		err := token.Error()
		if err != nil {
			log.Infof("token pub '%s'", err.Error())
		}
		return err
	}

	return s.txvr.TxCoap(txRaw, m, s.MtuOut())
}

func (s *MqttSesn) MgmtProto() sesn.MgmtProto {
	return s.cfg.MgmtProto
}

func (s *MqttSesn) ListenCoap(mc nmcoap.MsgCriteria) (*nmcoap.Listener, error) {
	return s.txvr.ListenCoap(mc)
}

func (s *MqttSesn) StopListenCoap(mc nmcoap.MsgCriteria) {
	s.txvr.StopListenCoap(mc)
}

func (s *MqttSesn) CoapIsTcp() bool {
	return false
}

func (s *MqttSesn) RxAccept() (sesn.Sesn, *sesn.SesnCfg, error) {
	return nil, nil, fmt.Errorf("op not implemented yet")
}

func (s *MqttSesn) RxCoap(opt sesn.TxOptions) (coap.Message, error) {
	return nil, fmt.Errorf("op not implemented yet")
}

func (s *MqttSesn) Filters() (nmcoap.TxMsgFilter, nmcoap.RxMsgFilter) {
	return s.txvr.Filters()
}

func (s *MqttSesn) SetFilters(txFilter nmcoap.TxMsgFilter,
	rxFilter nmcoap.RxMsgFilter) {

	s.txvr.SetFilters(txFilter, rxFilter)
}
