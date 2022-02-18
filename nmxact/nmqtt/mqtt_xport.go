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

	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
)

type MqttXPortCfg struct {
	RxTopic   string
	TxTopic   string
	DeviceId  int
	Broker    string
	Password  string
	User      string
	Id        string
	Cleansess bool
	Qos       int8
	// payload   string
	// action    string
	Store string
	// started bool
}

type MqttXPort struct {
	MqttCfg MqttXPortCfg
	started bool
}

func NewXportCfg() *MqttXPortCfg {
	return &MqttXPortCfg{
		RxTopic: "server_rx",
		TxTopic: "server_tx",
		Qos:     1,
		Store:   ":memory:",
	}
}
func NewMqttXport(cfg MqttXPortCfg) *MqttXPort {
	return &MqttXPort{
		MqttCfg: cfg,
	}
}

func (mx *MqttXPort) BuildSesn(cfg sesn.SesnCfg) (sesn.Sesn, error) {
	return NewMqttSesn(mx, cfg)
}

func (ux *MqttXPort) Start() error {
	if ux.started {
		return nmxutil.NewXportError("Mqtt xport started twice")
	}
	ux.started = true
	return nil
}

func (ux *MqttXPort) Stop() error {
	if !ux.started {
		return nmxutil.NewXportError("MqttXPort xport stopped twice")
	}
	ux.started = false
	return nil
}

func (ux *MqttXPort) Tx(bytes []byte) error {
	return fmt.Errorf("unsupported")
}
