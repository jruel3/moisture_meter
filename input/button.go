package input

import (
	"machine"
	"time"
)

// MomentaryTrigger sends a true bool event for a pin trigger.
func MomentaryTrigger(pin machine.Pin, eventChan chan bool) {
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	lastState := true

	for {
		currentState := pin.Get()

		if lastState && !currentState {
			time.Sleep(20 * time.Millisecond)
			if !pin.Get() {
				println("<<< Button Press >>>")
				select {
				case eventChan <- true:
				default:
					println("Chan Full")
				}
			}
		}

		lastState = currentState
		time.Sleep(20 * time.Millisecond)
	}
}
