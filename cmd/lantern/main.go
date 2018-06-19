package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut/palette"
	"github.com/shirou/gopsutil/cpu"

	"github.com/muesli/go-razer"
)

var (
	brightness = flag.Float64("brightness", -1, "brightness (between 0 and 100)")
	effect     = flag.String("effect", "", "effect mode (reactive, wave, spectrum, breath[dual,random], starlight[dual,random], ripple[random])")
	primary    = flag.String("color", "#ff0000", "Sets the primary keyboard color")
	secondary  = flag.String("secondary", "#00ff00", "secondary color (for 'dual' effect modes)")
	theme      = flag.String("theme", "", "theme mode (happy, soft, warm, rainbow, random)")
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
				c = colorful.Hsl(hue, 1.0, 0.02)
			}

			for y := 0; y < rows; y++ {
				k.Key(y, x).Color = c
			}
		}
		d.SetKeys(k)

		time.Sleep(300 * time.Millisecond)
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

func happyTheme(d razer.Device, theme string) {
	k := d.Keys()
	k.SetColor(color.RGBA{255, 0, 0, 0})

	var pal []colorful.Color
	switch theme {
	case "happy":
		pal = colorful.FastHappyPalette(9)
	case "warm":
		pal = colorful.FastWarmPalette(9)
	case "soft":
		pal, _ = colorful.SoftPalette(9)
	case "random":
		rand.Seed(time.Now().UTC().UnixNano())
		pal, _ = colorful.HappyPalette(9)
	case "monokai":
		for _, c := range palette.Monokai.Colors() {
			cf, _ := colorful.MakeColor(c.Color)
			fmt.Println("Color:", c.Name, cf.Hex())
			pal = append(pal, cf)
		}
	default:
		p := palette.AllPalettes().Filter(theme)
		if len(p) == 0 {
			log.Fatalf("Could not find colors by that name")
		}
		for _, c := range p {
			cf, _ := colorful.MakeColor(c.Color)
			fmt.Println("Color:", c.Name, cf.Hex())
			pal = append(pal, cf)
		}
	}

	k.FnKeys.SetColor(pal[0%len(pal)])
	k.Numerics.SetColor(pal[1%len(pal)])
	k.Cursor.SetColor(pal[2%len(pal)])
	k.Symbols.SetColor(pal[3%len(pal)])
	k.Commandos.SetColor(pal[4%len(pal)])
	k.Actions.SetColor(pal[5%len(pal)])
	k.Letters.SetColor(pal[6%len(pal)])
	k.Arrows.SetColor(pal[7%len(pal)])
	k.Special.SetColor(pal[8%len(pal)])
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
	if d.Brightness() == 0 && *brightness == -1 {
		*brightness = 100
	}
	if *brightness >= 0 {
		d.SetBrightness(*brightness)
	}

	if *theme == "rainbow" {
		rainbowTheme(d)
	} else if *theme != "" {
		happyTheme(d, *theme)
	} else if *effect != "" {
		d.SetEffect(razer.NewEffect(razer.StringToEffectType(*effect), primaryColor, secondaryColor))
	} else if *top {
		topMode(d)
	} else {
		d.SetEffect(razer.NewEffect(razer.EffectStatic, primaryColor, secondaryColor))
	}
}
