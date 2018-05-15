package razer

import (
	"image/color"
)

// Keys represents all keys (divided in rows) on the keyboard
type Keys struct {
	rows []Row
}

// A Row is a row of keys
type Row []color.Color

// SetColor sets the color of a single key
func (k *Keys) SetColor(row, col int, color color.Color) {
	k.rows[row][col] = color
}

// SetStaticColor sets a single color on all keys
func (k *Keys) SetStaticColor(color color.Color) {
	for y, r := range k.rows {
		for x := range r {
			k.SetColor(y, x, color)
		}
	}
}

// Keys returns the key-matrix for this device
func (d *Device) Keys() Keys {
	k := Keys{}

	rows, cols := d.MatrixDimensions()
	for y := 0; y < rows; y++ {
		r := Row{}
		for x := 0; x < cols; x++ {
			r = append(r, color.RGBA{})
		}
		k.rows = append(k.rows, r)
	}

	return k
}

// SetKeys sets the colors for the entire keyboard
func (d *Device) SetKeys(k Keys) {
	for i, r := range k.rows {
		var m []byte
		m = append(m, byte(i))
		m = append(m, 0)
		m = append(m, byte(len(r)-1))

		m = append(m, r.message()...)
		d.setKeyRow(m)
	}

	d.activateCustomSettings()
}

// setKeyRow sets the colors for an entire row of keys
func (d *Device) setKeyRow(row []byte) {
	dbusCall(d.dbusObject.Call("razer.device.lighting.chroma.setKeyRow", 0, row))
}

func (r Row) message() []byte {
	var m []byte

	for _, v := range r {
		var r, g, b uint32
		if v != nil {
			r, g, b, _ = v.RGBA()
		}
		m = append(m, uint8(r>>8))
		m = append(m, uint8(g>>8))
		m = append(m, uint8(b>>8))
	}

	return m
}
