package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/shirou/gopsutil/cpu"

	"github.com/muesli/go-razer"
)

var (
	col        = flag.String("color", "", "Sets the entire keyboard to this color")
	brightness = flag.Float64("brightness", -1, "brightness (between 0 and 100)")
	top        = flag.Bool("top", false, "top mode")
	theme      = flag.String("theme", "", "theme mode")
	effect     = flag.String("effect", "", "effect mode (wave, spectrum, breath, starlight, ripple)")
)

func main() {
	flag.Parse()

	devs, err := razer.Devices()
	if err != nil {
		panic(err)
	}
	if len(devs) == 0 {
		log.Fatalln("No Razer devices found.")
	}
	d := devs[0]

	fmt.Printf("Found %s\n", d.String())
	rows, cols := d.MatrixDimensions()
	k := d.Keys()

	if *brightness >= 0 {
		d.SetBrightness(*brightness)
	}

	if *theme == "happy" {
		k.SetColor(color.RGBA{255, 0, 0, 0})
		pal := colorful.FastHappyPalette(9)
		for _, key := range k.FnKeys {
			key.Color = pal[0]
		}
		for _, key := range k.Numerics {
			key.Color = pal[1]
		}
		for _, key := range k.Letters {
			key.Color = pal[6]
		}
		for _, key := range k.Symbols {
			key.Color = pal[3]
		}
		for _, key := range k.Commandos {
			key.Color = pal[4]
		}
		for _, key := range k.Actions {
			key.Color = pal[5]
		}
		for _, key := range k.Cursor {
			key.Color = pal[2]
		}
		for _, key := range k.Arrows {
			key.Color = pal[7]
		}
		for _, key := range k.Special {
			key.Color = pal[8]
		}
		d.SetKeys(k)
	} else if *effect != "" {
		d.SetEffect(razer.StringToEffect(*effect))
	} else if *top {
		base := 2.5
		for {
			cpuUsage, err := cpu.Percent(0, false)
			if err != nil {
				panic(err)
			}
			// fmt.Println("CPU usage:", cpuUsage[0])

			for x := 0; x < cols; x++ {
				var hue = ((base - (float64(x) / float64(cols-1))) * 120)
				c := colorful.Hsl(hue, 1.0, 0.5)

				if x > int(float64(cols-1)*(cpuUsage[0]/100)) {
					c = colorful.Hsl(hue, 1.0, 0.025)
				}

				for y := 0; y < rows; y++ {
					k.Key(y, x).Color = c
				}
			}
			d.SetKeys(k)

			time.Sleep(200 * time.Millisecond)
			base -= 0.015
			if base < 0 {
				base = 3
			}
		}
	} else if *col != "" {
		c, err := colorful.Hex(*col)
		if err != nil {
			panic(err)
		}

		var a []interface{}
		r, g, b, _ := c.RGBA()
		a = append(a, uint8(r>>8))
		a = append(a, uint8(g>>8))
		a = append(a, uint8(b>>8))

		d.SetEffect(razer.EffectStatic, a...)
	} else {
		pal := colorful.FastHappyPalette(rows)
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				k.Key(y, x).Color = pal[y]
				d.SetKeys(k)
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
}
