package main

import (
	"machine"
	"time"

	"moisture_meter/input"
	"moisture_meter/output"
)

var (
	pollFreq  = (10 * time.Second)
	eventChan = make(chan bool, 2)

	sensors = input.SensorCluster{
		Sensors: []input.Sensor{
			input.Sensor{
				Name:       "Moisture Meter 1",
				MachinePin: machine.I2C0,
				PinConfig: machine.I2CConfig{
					Frequency: 100 * machine.KHz,
				},
				Address:    0x36,
				WriteBuf:   []byte{0x0F, 0x10},
				MaxRange:   1025,
				MinRange:   320,
				AlertValue: 25,
			},
		},
		ManualPoll: false,
	}
	pushButton = machine.Pin(13)
	led        = output.LED{Pin: machine.Pin(4)}
)

func main() {
	sensors.Init()
	if err := sensors.CheckErr("Encountered an error initializing sensor cluster."); err {
		return
	}
	led.Init()

	ticker := time.NewTicker(pollFreq)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			eventChan <- false
		}
	}()

	go input.MomentaryTrigger(pushButton, eventChan)

	for manPoll := range eventChan {
		sensors.ManPoll(manPoll)
		if err := sensors.CheckErr("Error polling sensors: "); err {
			return
		}
		var alert = false
		for _, v := range sensors.Sensors {
			if v.NormValue < v.AlertValue {
				alert = true
			}
		}
		if alert {
			led.On()
		} else {
			led.Off()
		}
		resp, err := sensors.JSON()
		if err != nil {
			println(err)
		}
		println(string(resp))

		time.Sleep(200 * time.Millisecond)
	}
}
