package razer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/godbus/dbus"
)

// Device represents a single Razer device
type Device struct {
	Name       string
	keys       Keys
	dbusObject dbus.BusObject
}

// Devices returns the available Razer devices
func Devices() ([]Device, error) {
	d := []Device{}
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		return d, err
	}

	var s []string
	err = dbusConn.Object("org.razer", dbus.ObjectPath("/org/razer")).
		Call("razer.devices.getDevices", 0).Store(&s)
	if err != nil {
		return d, err
	}

	for _, vs := range s {
		d = append(d, Device{
			Name: vs,
			dbusObject: dbusConn.Object("org.razer",
				dbus.ObjectPath(fmt.Sprintf("/org/razer/device/%s", vs))),
		})
	}

	return d, nil
}

func (d *Device) String() string {
	rows, columns := d.MatrixDimensions()

	return fmt.Sprintf(
		"%s (type %s, serial: %s)\n\t"+
			"Dimensions: %dx%d\n\t"+
			"Brightness: %.2f%%\n\t"+
			"Firmware: %s\n\t"+
			"GameMode: %s",
		d.Name, d.Type(), d.Serial(),
		rows, columns, d.Brightness(),
		d.Firmware(), strconv.FormatBool(d.GameMode()))
}

// Type returns the device type (e.g. "keyboard")
func (d *Device) Type() string {
	var s string
	err := dbusCall(d.dbusObject.Call("razer.device.misc.getDeviceType", 0)).Store(&s)
	if err != nil {
		log.Printf("reading device type failed: %s", err)
	}

	return s
}

// Serial returns the serial no of the device
func (d *Device) Serial() string {
	var s string
	err := dbusCall(d.dbusObject.Call("razer.device.misc.getSerial", 0)).Store(&s)
	if err != nil {
		log.Printf("reading device serial failed: %s", err)
	}

	return s
}

// Firmware returns the firmware version of the device
func (d *Device) Firmware() string {
	var s string
	err := dbusCall(d.dbusObject.Call("razer.device.misc.getFirmware", 0)).Store(&s)
	if err != nil {
		log.Printf("reading firmware version failed: %s", err)
	}

	return s
}

// MatrixDimensions returns the matrix dimensions of the device (rows & columns)
func (d *Device) MatrixDimensions() (int, int) {
	var i []int
	err := dbusCall(d.dbusObject.Call("razer.device.misc.getMatrixDimensions", 0)).Store(&i)
	if err != nil {
		log.Printf("reading matrixDimensions failed: %s", err)
	}

	return i[0], i[1]
}

// GameMode returns whether the device is currently in "game-mode"
func (d *Device) GameMode() bool {
	var b bool
	err := dbusCall(d.dbusObject.Call("razer.device.led.gamemode.getGameMode", 0)).Store(&b)
	if err != nil {
		log.Printf("reading gameMode failed: %s", err)
	}

	return b
}

// HasDedicatedMacroKeys returns whether the device has dedicated macro keys
func (d *Device) HasDedicatedMacroKeys() bool {
	var b bool
	err := dbusCall(d.dbusObject.Call("razer.device.misc.hasDedicatedMacroKeys", 0)).Store(&b)
	if err != nil {
		log.Printf("reading hasDedicatedMacroKeys failed: %s", err)
	}

	return b
}

// Brightness returns the device's current brightness
func (d *Device) Brightness() float64 {
	var b float64
	err := dbusCall(d.dbusObject.Call("razer.device.lighting.brightness.getBrightness", 0)).Store(&b)
	if err != nil {
		log.Printf("reading brightness failed: %s", err)
	}

	return b
}

// SetBrightness sets the brightness (between 0 and 100, in percent)
func (d *Device) SetBrightness(b float64) {
	dbusCall(d.dbusObject.Call("razer.device.lighting.brightness.setBrightness", 0, b))
}

// activateCustomSettings activates the custom color settings
func (d *Device) activateCustomSettings() {
	dbusCall(d.dbusObject.Call("razer.device.lighting.chroma.setCustom", 0))
}

func dbusCall(call *dbus.Call) *dbus.Call {
	if call.Err != nil {
		log.Printf("dbus call failed: %s", call.Err)
	}

	return call
}
