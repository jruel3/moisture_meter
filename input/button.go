package input

import (
	"machine"
	"time"
)

// MomentaryTrigger sends an event for a pin trigger.
func MomentaryTrigger(pin machine.Pin, eventChan chan struct{}) {
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	lastState := true

	for {
		currentState := pin.Get()

		if lastState && !currentState {
			time.Sleep(20 * time.Millisecond)
			if !pin.Get() {
				println("<<< Button Press >>>")
				select {
				case eventChan <- struct{}{}:
				default:
					println("Default Case")
				}
			}
		}

		lastState = currentState
		time.Sleep(20 * time.Millisecond)
	}
}
