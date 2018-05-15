package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/shirou/gopsutil/cpu"

	"github.com/muesli/go-razer"
)

var (
	col        = flag.String("color", "", "Sets the entire keyboard to this color")
	top        = flag.Bool("top", false, "top mode")
	effect     = flag.String("effect", "", "effect mode (wave, spectrum, breath, starlight, ripple)")
	brightness = flag.Float64("brightness", -1, "brightness (between 0 and 100)")
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

	if *effect != "" {
		d.SetEffect(razer.StringToEffect(*effect))
	} else if *top {
		for {
			cpuUsage, err := cpu.Percent(0, false)
			if err != nil {
				panic(err)
			}
			// fmt.Println("CPU usage:", cpuUsage[0])

			for x := 0; x < cols; x++ {
				var hue = ((1.5 + (float64(x) / float64(cols-1))) * 120)
				c := colorful.Hsl(hue, 1.0, 0.5)

				if x > int(float64(cols-1)*(cpuUsage[0]/100)) {
					c = colorful.Hsl(hue, 1.0, 0.015)
				}

				for y := 0; y < rows; y++ {
					k.SetColor(y, x, c)
				}
			}
			d.SetKeys(k)

			time.Sleep(200 * time.Millisecond)
		}
	} else if *col != "" {
		c, err := colorful.Hex(*col)
		if err != nil {
			panic(err)
		}
		k.SetStaticColor(c)
		d.SetKeys(k)
	} else {
		pal := colorful.FastHappyPalette(rows)
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				k.SetColor(y, x, pal[y])
				d.SetKeys(k)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
