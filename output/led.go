package output

import (
	"errors"
	"machine"
)

// Type LED declares a LED object
type LED struct {
	Pin        machine.Pin
	configured bool // Defaults to "false"
}

// Type MultiLED declares a map of LEDs
type MultiLED map[string]*LED

var (
	// LEDNotInitilaized is returned when an action is requested on an uninitialized LED
	LEDNotInitialized = errors.New("requested LED pin has not been initalized")

	outputPinConfig = machine.PinConfig{Mode: machine.PinOutput}
)

// Init initializes an LED
func (l *LED) Init() {
	l.Pin.Configure(outputPinConfig)
	l.configured = true
	l.Pin.Low()
}

// Init initializes LEDs in a MultiLED object
func (m *MultiLED) Init() {
	for _, v := range *m {
		v.Init()
	}
}

// On sets an LED pin to high
func (l LED) On() error {
	if err := l.valid(); err != nil {
		return err
	}
	l.Pin.High()
	return nil
}

// Off sets an LED pin to low
func (l LED) Off() error {
	if err := l.valid(); err != nil {
		return err
	}
	l.Pin.Low()
	return nil
}

// Off sets all LED pins to low in a MultiLED
func (m MultiLED) Off() error {
	for _, v := range m {
		if err := v.Off(); err != nil {
			return err
		}
	}
	return nil
}

// Toggle swaps the state of the LED pin
func (l LED) Toggle() error {
	if err := l.valid(); err != nil {
		return err
	}
	l.Pin.Set(!l.Pin.Get())
	return nil
}

func (l LED) valid() error {
	if !l.configured {
		return LEDNotInitialized
	}
	return nil
}
