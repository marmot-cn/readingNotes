package main 

import ( 
	"image" 
	"image/color" 
	"image/gif" 
	"io"
	"math" 
	"math/rand" 
	"os" 
	"time"
)


var green = color.RGBA{0x00,0xff,0x00,0xff}
var red = color.RGBA{0xff,0x00,0x00,0xff}
var blue = color.RGBA{0x00,0x00,0xff,0xff}

var palette = []color.Color{color.White,  color.Black, green, red, blue} 

const ( 
	whiteIndex = 0 // first color in palette 
	blackIndex = 1 // next color in palette 
	greenIndex = 2
	redIndex = 3
	cyanIndex = 4
)

const totalColor = 5

func main() { 
	// The sequence of images is deterministic unless we seed 
	// the pseudo-random number generator using the current time. 
	// Thanks to Randall McPherson for pointing out the omission. 
	rand.Seed(time.Now().UTC().UnixNano()) 
	lissajous(os.Stdout) 
}

func lissajous(out io.Writer) { 
	const ( 
		cycles = 5 // number of complete x oscillator revolutions 
		res = 0.001 // angular resolution 
		size = 100 // image canvas covers [-size..+size] 
		nframes = 64 // number of animation frames 
		delay = 8 // delay between frames in 10ms units 
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator 
	anim := gif.GIF{LoopCount: nframes} 
	phase := 0.0 // phase difference 
	for i := 0; i < nframes; i++ { 
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) 
		img := image.NewPaletted(rect, palette) 
		for t := 0.0; t < cycles*2*math.Pi; t += res { 
			x := math.Sin(t) 
			y := math.Sin(t*freq + phase) 
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(rand.Intn(5))) 
		}

		phase += 0.1 
		anim.Delay = append(anim.Delay, delay) 
		anim.Image = append(anim.Image, img) 
	}

	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors 
}

//执行
//root@daa41eb01821:/go/01/homework# go build 6.go
//root@daa41eb01821:/go/01/homework# ./6 > 6.gif