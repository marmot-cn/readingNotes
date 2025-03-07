package main 

import ( 
	"image" 
	"image/color" 
	"image/gif" 
	"io"
	"math" 
	"math/rand" 
	"time"
	"log" 
	"net/http" 
	"strconv"
)

var palette = []color.Color{color.White, color.Black} 

const ( 
	whiteIndex = 0 // first color in palette 
	blackIndex = 1 // next color in palette 
)

func main() { 
	// The sequence of images is deterministic unless we seed 
	// the pseudo-random number generator using the current time. 
	// Thanks to Randall McPherson for pointing out the omission. 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil { 
			log.Print(err) 
		}
		cycles := r.Form["cycles"][0]

		c, _:= strconv.ParseFloat(cycles, 64) 
		lissajous(w, c)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil)) 
	rand.Seed(time.Now().UTC().UnixNano()) 
}

func lissajous(out io.Writer, c float64) { 

	const ( 
		res = 0.001 // angular resolution 
		size = 100 // image canvas covers [-size..+size] 
		nframes = 64 // number of animation frames 
		delay = 8 // delay between frames in 10ms units 
	)

	cycles := c
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator 
	anim := gif.GIF{LoopCount: nframes} 
	phase := 0.0 // phase difference 
	for i := 0; i < nframes; i++ { 
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) 
		img := image.NewPaletted(rect, palette) 
		for t := 0.0; t < cycles*2*math.Pi; t += res { 
			x := math.Sin(t) 
			y := math.Sin(t*freq + phase) 
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex) 
		}

		phase += 0.1 
		anim.Delay = append(anim.Delay, delay) 
		anim.Image = append(anim.Image, img) 
	}

	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors 
}

//执行 
//root@daa41eb01821:/go/01# go build lissajous.go
//root@daa41eb01821:/go/01# ./lissajous >out.gif