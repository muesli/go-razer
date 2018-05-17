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
	primary    = flag.String("color", "#ff0000", "Sets the primary keyboard color")
	secondary  = flag.String("secondary", "#00ff00", "secondary color (for 'dual' effect modes)")
	brightness = flag.Float64("brightness", -1, "brightness (between 0 and 100)")
	effect     = flag.String("effect", "", "effect mode (reactive, wave, spectrum, breath[dual,random], starlight[dual,random], ripple[random])")
	theme      = flag.String("theme", "", "theme mode (happy, rainbow)")
	top        = flag.Bool("top", false, "top mode")
)

func topMode(d razer.Device) {
	rows, cols := d.MatrixDimensions()
	k := d.Keys()

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
}

func rainbowTheme(d razer.Device) {
	rows, cols := d.MatrixDimensions()
	k := d.Keys()

	pal := colorful.FastHappyPalette(rows)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			k.Key(y, x).Color = pal[y]
			d.SetKeys(k)
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func happyTheme(d razer.Device) {
	k := d.Keys()

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
}

func main() {
	flag.Parse()

	primaryColor, err := colorful.Hex(*primary)
	if err != nil {
		panic(err)
	}
	secondaryColor, err := colorful.Hex(*secondary)
	if err != nil {
		panic(err)
	}

	devs, err := razer.Devices()
	if err != nil {
		panic(err)
	}
	if len(devs) == 0 {
		log.Fatalln("No Razer devices found.")
	}
	d := devs[0]

	fmt.Printf("Found %s\n", d.String())
	if *brightness >= 0 {
		d.SetBrightness(*brightness)
	}

	if *theme == "rainbow" {
		rainbowTheme(d)
	} else if *theme == "happy" {
		happyTheme(d)
	} else if *effect != "" {
		d.SetEffect(razer.NewEffect(razer.StringToEffectType(*effect), primaryColor, secondaryColor))
	} else if *top {
		topMode(d)
	} else {
		d.SetEffect(razer.NewEffect(razer.EffectStatic, primaryColor, secondaryColor))
	}
}
