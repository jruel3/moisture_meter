// Package input provides functionality for fetching input data
package input

import (
	"bytes"
	"encoding/json"
	"fmt"
	"machine"
	"time"
)

// Type Sensor delcares a sensor object
type Sensor struct {
	Name       string
	MachinePin *machine.I2C      `json:"-"`
	PinConfig  machine.I2CConfig `json:"-"`
	Address    uint16            `json:"-"`
	WriteBuf   []byte            `json:"-"`
	MaxRange   float32
	MinRange   float32
	RawValue   uint16
	NormValue  float32
	Err        error
}

// Type SensorCluster represents a slice of Sensors
type SensorCluster []Sensor

// Check a cluster for errors and print them to console.
// Returns a bool (true if an error is found).
func (sc SensorCluster) CheckErr(message string) bool {
	errFound := false
	for _, s := range sc {
		if s.Err != nil {
			if !errFound {
				println(message)
			}
			println(fmt.Sprintf("{%s: Raw_Value: %v, Error: %v}\n", s.Name, s.RawValue, s.Err))
			errFound = true
		}
	}
	return errFound
}

// Initialize I2C pins for a sensor
func (s *Sensor) init() {
	if err := s.MachinePin.Configure(s.PinConfig); err != nil {
		s.Err = err
	}
}

// Initialize I2C pins in a cluster
func (sc SensorCluster) Init() {
	for i := range sc {
		sc[i].init()
	}
}

// Export SensorCluster as a JSON
// Custom Marshaller required given TinyGo constraints.
func (sc SensorCluster) JSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteByte('{')

	for i, s := range sc {
		b, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		buf.Write(b)

		if i < len(sc)-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// Normalize the raw value to a percentage between the MinRange and MaxRange
func (s *Sensor) normalize() {
	s.NormValue = (float32(s.RawValue) - s.MinRange) * (100) / (s.MaxRange - s.MinRange)
}

// Polls a Sensor and sets the Raw and Normalized Values
func (s *Sensor) poll() {
	s.resetVals()
	if err := s.MachinePin.Tx(s.Address, s.WriteBuf, nil); err != nil {
		s.Err = err
		return
	}
	time.Sleep(5 * time.Millisecond)
	readBuf := []byte{0, 0}
	if err := s.MachinePin.Tx(s.Address, nil, readBuf); err != nil {
		s.Err = err
		return
	}
	s.RawValue = (uint16(readBuf[0]) << 8) | uint16(readBuf[1])
	s.normalize()
}

// Polls the Sensors in the SensorCluster and sets their values
func (sc SensorCluster) Poll() {
	for i := range sc {
		sc[i].poll()
		time.Sleep(5 * time.Millisecond)
	}
}

// Reregister the sensor's address on the I2C bus.
// IMPORTANT: Only run this with a single sensor connected.
func (s Sensor) RegisterAddress(newAddr uint16) error {
	// Base: 0x00 (Status), Function: 0x11 (EEPROM I2C Address Register)
	// Followed by the new 7-bit address shifted left by 1 byte
	cmd := []byte{0x00, 0x11, byte(newAddr)}

	err := s.MachinePin.Tx(s.Address, cmd, nil)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Millisecond)

	// Trigger a Software Reset to apply changes
	// Base: 0x00 (Status), Function: 0x7F (Reset)
	resetCmd := []byte{0x00, 0x7F, 0xFF}
	_ = s.MachinePin.Tx(s.Address, resetCmd, nil)
	s.Address = newAddr
	return nil
}

// Resets value and err fields in a Sensor
func (s *Sensor) resetVals() {
	s.RawValue = 0
	s.NormValue = 0.0
	s.Err = nil
}

// Reset value and err fields in Sensors in a cluster
func (sc SensorCluster) ResetVals() {
	for i := range sc {
		sc[i].resetVals()
	}
}
