package razer

import (
	"fmt"
	"image/color"
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

type Effect struct {
	Name string
	Type EffectType

	primary   color.Color
	secondary color.Color
}

// Effect type safety
type EffectType int

func NewEffect(e EffectType, primary, secondary color.Color) Effect {
	return Effect{
		Name:      effectToDBusMethod(e),
		Type:      e,
		primary:   primary,
		secondary: secondary,
	}
}

// StringToEffectType converts a string to an Effect
func StringToEffectType(s string) EffectType {
	switch strings.ToLower(s) {
	case "static":
		return EffectStatic
	case "reactive":
		return EffectReactive
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
	case "wave":
		return EffectWave
	}

	return EffectNone
}

// effectToDBusMethod returns the name of the DBus method matching the Effect
func effectToDBusMethod(e EffectType) string {
	switch e {
	case EffectStatic:
		return "Static"
	case EffectReactive:
		return "Reactive"
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
	case EffectWave:
		return "Wave"
	}

	return "None"
}

func (e Effect) arguments() []interface{} {
	var a []interface{}

	switch e.Type {
	case EffectStatic:
		fallthrough
	case EffectBreath:
		a = append(a, colorToEffectArg(e.primary)...)
	case EffectReactive:
		a = append(a, colorToEffectArg(e.primary)...)
		a = append(a, 1)
	case EffectBreathDual:
		a = append(a, colorToEffectArg(e.primary)...)
		a = append(a, colorToEffectArg(e.secondary)...)
	case EffectStarlight:
		a = append(a, 100)
		a = append(a, colorToEffectArg(e.primary)...)
	case EffectStarlightDual:
		a = append(a, 100)
		a = append(a, colorToEffectArg(e.primary)...)
		a = append(a, colorToEffectArg(e.secondary)...)
	case EffectStarlightRandom:
		a = append(a, 100)
	case EffectRipple:
		a = append(a, colorToEffectArg(e.primary)...)
		a = append(a, 0.0)
	case EffectRippleRandom:
		a = append(a, 0.0)
	case EffectWave:
		a = append(a, 1)
	}

	return a
}

// SetEffect activates an Effect
func (d *Device) SetEffect(effect Effect) {
	var call *dbus.Call

	switch effect.Type {
	case EffectRipple:
		fallthrough
	case EffectRippleRandom:
		call = d.dbusObject.Call(fmt.Sprintf("razer.device.lighting.custom.set%s", effect.Name), 0, effect.arguments()...)
	default:
		call = d.dbusObject.Call(fmt.Sprintf("razer.device.lighting.chroma.set%s", effect.Name), 0, effect.arguments()...)
	}

	dbusCall(call)
}

// colorToEffectArg converts a color to a valid effect argument
func colorToEffectArg(c color.Color) []interface{} {
	var a []interface{}
	r, g, b, _ := c.RGBA()
	a = append(a, uint8(r>>8))
	a = append(a, uint8(g>>8))
	a = append(a, uint8(b>>8))

	return a
}
