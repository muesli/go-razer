package razer

import (
	"image/color"
)

// Key represents a single key
type Key struct {
	Row   int
	Col   int
	Color color.Color
}

// A KeySet is a set of keys
type KeySet []*Key

// Keys represents all keys (divided in rows) on the keyboard
type Keys struct {
	Letters   KeySet
	FnKeys    KeySet
	Numerics  KeySet
	Symbols   KeySet
	Commandos KeySet
	Actions   KeySet
	Cursor    KeySet
	Arrows    KeySet
	Special   KeySet

	rows []KeySet
}

// SetColor sets a single color on all keys in the KeySet
func (k KeySet) SetColor(c color.Color) {
	for _, key := range k {
		key.Color = c
	}
}

// SetColor sets a single color on all keys
func (k *Keys) SetColor(c color.Color) {
	for _, r := range k.rows {
		r.SetColor(c)
	}
}

// Key returns a single key
func (k *Keys) Key(row, col int) *Key {
	return k.rows[row][col]
}

// KeySpan returns a span of keys
func (k *Keys) KeySpan(row, startCol, endCol int) KeySet {
	var kk KeySet

	for x := startCol; x <= endCol; x++ {
		kk = append(kk, k.rows[row][x])
	}

	return kk
}

// Keys returns the key-matrix for this device
func (d *Device) Keys() Keys {
	if len(d.keys.rows) > 0 {
		return d.keys
	}

	k := Keys{}
	rows, cols := d.MatrixDimensions()
	for y := 0; y < rows; y++ {
		r := KeySet{}
		for x := 0; x < cols; x++ {
			r = append(r, &Key{
				Row: y,
				Col: x,
			})
		}
		k.rows = append(k.rows, r)
	}

	if rows == 6 && cols == 22 {
		// Ornata
		k.FnKeys = k.KeySpan(0, 2, 14)
		k.Letters = append(k.Letters, k.KeySpan(2, 2, 11)...)
		k.Letters = append(k.Letters, k.KeySpan(3, 2, 10)...)
		k.Letters = append(k.Letters, k.KeySpan(4, 3, 9)...)
		k.Numerics = append(k.Numerics, k.KeySpan(1, 2, 11)...)
		k.Numerics = append(k.Numerics, k.KeySpan(2, 18, 20)...)
		k.Numerics = append(k.Numerics, k.KeySpan(3, 18, 20)...)
		k.Numerics = append(k.Numerics, k.KeySpan(4, 18, 20)...)
		k.Numerics = append(k.Numerics, k.KeySpan(5, 18, 19)...)
		k.Symbols = append(k.Symbols, k.KeySpan(1, 1, 1)...)
		k.Symbols = append(k.Symbols, k.KeySpan(1, 12, 13)...)
		k.Symbols = append(k.Symbols, k.KeySpan(2, 12, 13)...)
		k.Symbols = append(k.Symbols, k.KeySpan(3, 11, 13)...)
		k.Symbols = append(k.Symbols, k.KeySpan(4, 2, 2)...)
		k.Symbols = append(k.Symbols, k.KeySpan(4, 10, 12)...)
		k.Symbols = append(k.Symbols, k.KeySpan(1, 19, 21)...)
		k.Symbols = append(k.Symbols, k.KeySpan(2, 21, 21)...)
		k.Symbols = append(k.Symbols, k.KeySpan(5, 20, 20)...)
		k.Commandos = append(k.Commandos, k.KeySpan(2, 1, 1)...)
		k.Commandos = append(k.Commandos, k.KeySpan(3, 1, 1)...)
		k.Commandos = append(k.Commandos, k.KeySpan(4, 1, 1)...)
		k.Commandos = append(k.Commandos, k.KeySpan(4, 14, 14)...)
		k.Commandos = append(k.Commandos, k.KeySpan(5, 1, 3)...)
		k.Commandos = append(k.Commandos, k.KeySpan(5, 11, 14)...)
		k.Actions = append(k.Actions, k.KeySpan(0, 1, 1)...)
		k.Actions = append(k.Actions, k.KeySpan(1, 14, 14)...)
		k.Actions = append(k.Actions, k.KeySpan(3, 14, 14)...)
		k.Actions = append(k.Actions, k.KeySpan(5, 4, 10)...)
		k.Actions = append(k.Actions, k.KeySpan(4, 21, 21)...)
		k.Cursor = append(k.Cursor, k.KeySpan(1, 15, 17)...)
		k.Cursor = append(k.Cursor, k.KeySpan(2, 15, 17)...)
		k.Arrows = append(k.Arrows, k.KeySpan(4, 16, 16)...)
		k.Arrows = append(k.Arrows, k.KeySpan(5, 15, 17)...)
		k.Special = append(k.Special, k.KeySpan(0, 15, 17)...)
		k.Special = append(k.Special, k.KeySpan(1, 18, 18)...)
	}

	d.keys = k
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

	d.keys = k
	d.activateCustomSettings()
}

// setKeyRow sets the colors for an entire row of keys
func (d *Device) setKeyRow(row []byte) {
	dbusCall(d.dbusObject.Call("razer.device.lighting.chroma.setKeyRow", 0, row))
}

func (k KeySet) message() []byte {
	var m []byte

	for _, v := range k {
		var r, g, b uint32
		if v.Color != nil {
			r, g, b, _ = v.Color.RGBA()
		}
		m = append(m, uint8(r>>8))
		m = append(m, uint8(g>>8))
		m = append(m, uint8(b>>8))
	}

	return m
}
