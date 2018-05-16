package razer

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

// Effects known by the keyboard hardware itself
const (
	EffectNone = iota
	EffectStatic
	EffectReactive
	EffectWave
	EffectSpectrum
	EffectBreath
	EffectBreathDual
	EffectBreathRandom
	EffectStarlight
	EffectStarlightDual
	EffectStarlightRandom
	EffectRipple
	EffectRippleRandom
)

// Effect type safety
type Effect int

// StringToEffect converts a string to an Effect
func StringToEffect(s string) Effect {
	switch strings.ToLower(s) {
	case "static":
		return EffectStatic
	case "reactive":
		return EffectReactive
	case "wave":
		return EffectWave
	case "spectrum":
		return EffectSpectrum
	case "breath":
		return EffectBreath
	case "breathdual":
		return EffectBreathDual
	case "breathrandom":
		return EffectBreathRandom
	case "starlight":
		return EffectStarlight
	case "starlightdual":
		return EffectStarlightDual
	case "starlightrandom":
		return EffectStarlightRandom
	case "ripple":
		return EffectRipple
	case "ripplerandom":
		return EffectRippleRandom
	}

	return EffectNone
}

// effectToDBusMethod returns the name of the DBus method matching the Effect
func effectToDBusMethod(e Effect) string {
	switch e {
	case EffectStatic:
		return "Static"
	case EffectReactive:
		return "Reactive"
	case EffectWave:
		return "Wave"
	case EffectSpectrum:
		return "Spectrum"
	case EffectBreath:
		return "BreathSingle"
	case EffectBreathDual:
		return "BreathDual"
	case EffectBreathRandom:
		return "BreathRandom"
	case EffectStarlight:
		return "StarlightSingle"
	case EffectStarlightDual:
		return "StarlightDual"
	case EffectStarlightRandom:
		return "StarlightRandom"
	case EffectRipple:
		return "Ripple"
	case EffectRippleRandom:
		return "RippleRandomColour"
	}

	return "None"
}

func defaultEffectArgs(e Effect) []interface{} {
	var a []interface{}

	switch e {
	case EffectStatic:
		a = append(a, 0)
		a = append(a, 255)
		a = append(a, 0)
	case EffectReactive:
		a = append(a, 255)
		a = append(a, 255)
		a = append(a, 0)
		a = append(a, 1)
	case EffectWave:
		a = append(a, 1)
	case EffectBreath:
		a = append(a, 255)
		a = append(a, 0)
		a = append(a, 0)
	case EffectBreathDual:
		a = append(a, 255)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 255)
	case EffectStarlight:
		a = append(a, 100)
		a = append(a, 255)
		a = append(a, 0)
		a = append(a, 0)
	case EffectStarlightDual:
		a = append(a, 100)
		a = append(a, 255)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 0)
		a = append(a, 255)
	case EffectStarlightRandom:
		a = append(a, 100)
	case EffectRipple:
		a = append(a, 0.0)
	case EffectRippleRandom:
		a = append(a, 0.0)
	}

	return a
}

// SetEffect activates an Effect
func (d *Device) SetEffect(effect Effect, args ...interface{}) {
	var call *dbus.Call

	if len(args) == 0 {
		args = defaultEffectArgs(effect)
	}

	switch effect {
	case EffectRipple:
		fallthrough
	case EffectRippleRandom:
		call = d.dbusObject.Call(fmt.Sprintf("razer.device.lighting.custom.set%s", effectToDBusMethod(effect)), 0, args...)
	default:
		call = d.dbusObject.Call(fmt.Sprintf("razer.device.lighting.chroma.set%s", effectToDBusMethod(effect)), 0, args...)
	}

	dbusCall(call)
}
