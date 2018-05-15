package razer

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

// Effects known by the keyboard hardware itself
const (
	EffectNone = iota
	EffectWave
	EffectSpectrum
	EffectBreath
	EffectStarlight
	EffectRipple
)

// Effect type safety
type Effect int

// StringToEffect converts a string to an Effect
func StringToEffect(s string) Effect {
	switch strings.ToLower(s) {
	case "wave":
		return EffectWave
	case "spectrum":
		return EffectSpectrum
	case "breath":
		return EffectBreath
	case "starlight":
		return EffectStarlight
	case "ripple":
		return EffectRipple
	}

	return EffectNone
}

// effectToDBusMethod returns the name of the DBus method matching the Effect
func effectToDBusMethod(e Effect) string {
	switch e {
	case EffectWave:
		return "Wave"
	case EffectSpectrum:
		return "Spectrum"
	case EffectBreath:
		return "BreathRandom"
	case EffectStarlight:
		return "StarlightRandom"
	case EffectRipple:
		return "RippleRandomColour"
	}

	return "None"
}

// SetEffect activates an Effect
func (d *Device) SetEffect(effect Effect) {
	var call *dbus.Call

	switch effect {
	case EffectWave:
		call = d.dbusObject.Call("razer.device.lighting.chroma.setWave", 0, 1)
	case EffectStarlight:
		call = d.dbusObject.Call("razer.device.lighting.chroma.setStarlightRandom", 0, 100)
	case EffectRipple:
		call = d.dbusObject.Call("razer.device.lighting.custom.setRippleRandomColour", 0, 0.0)
	case EffectSpectrum:
		fallthrough
	case EffectBreath:
		call = d.dbusObject.Call(fmt.Sprintf("razer.device.lighting.chroma.set%s", effectToDBusMethod(effect)), 0)
	}

	dbusCall(call)
}
