# LGTV

LGTV is a [Go](http://golang.org/) package for controlling LG Smart TVs. It uses [Gorilla WebSocket](https://github.com/gorilla/websocket) for the connection to the TV.

[![Build Status](https://travis-ci.org/mjanser/lgtv.svg?branch=master)](https://travis-ci.org/mjanser/lgtv)
[![GoDoc](https://godoc.org/github.com/mjanser/lgtv?status.svg)](https://godoc.org/github.com/mjanser/lgtv)

## Installation

```
go get github.com/mjanser/lgtv
```

## Examples

```
package main

import (
	"github.com/mjanser/lgtv"
)

func main() {
	ip := "192.168.1.2"
	key := ""

	tv, err := lgtv.NewDefaultClient(ip, key)
	err = tv.Connect()
	defer tv.Disconnect()

	volume, err := tv.GetVolume()
	channel, err := tv.GetCurrentChannel()

	channels, err := tv.GetChannels()
	inputs, err := tv.GetInputs()
	apps, err := tv.GetApps()

	err = tv.SetVolume(10)
	err = tv.SwitchInput("HDMI_1")
	err = tv.SetChannel("1")

	err = tv.SubscribeVolume(func (volume lgtv.Volume) {
	})
	err = tv.SubscribeApp(func (app lgtv.AppInfo) {
	})
	err = tv.SubscribeChannel(func (channel lgtv.Channel) {
	})
```
