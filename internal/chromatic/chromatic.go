package chromatic

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
	"sync"
	"time"

	"github.com/GetVivid/huego"
	"github.com/Khabi/chromatic/internal/extract"
	"github.com/Khabi/chromatic/internal/location"
	"github.com/korandiz/v4l"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/paulbellamy/ratecounter"
	"github.com/sirupsen/logrus"
)

type State int

var fps *ratecounter.RateCounter

const (
	Running State = iota
	Paused
	Stop
	Status
)

func (s State) String() string {
	return [...]string{"running", "paused", "stoping"}[s]
}

type ServerStatus struct {
	State string
	FPS   int64
}

func Run(command <-chan State, status chan ServerStatus, video *v4l.Device, hue *huego.EntertainmentGroup, bounds location.Bounds) {
	fps = ratecounter.NewRateCounter(1 * time.Second)

	var stream *huego.EntertainmentStream
	defer stream.StopStream()
	var state = Paused
	var err error
	for {
		select {
		case cmd := <-command:
			switch cmd {
			case Running:
				state = Running
				stream, err = hue.StartStream()
				if err != nil {
					logrus.WithError(err).Error("unable to capture")
					os.Exit(1)
				}
				err = video.TurnOn()
				if err != nil {
					logrus.WithError(err).Error("unable to capture")
					os.Exit(1)
				}
				logrus.Info("starting capture")

			case Paused:
				state = Paused
				video.TurnOff()
				stream.StopStream()
				logrus.Info("pausing capture")
			case Stop:
				logrus.Info("stopping")
				return
			case Status:
				logrus.Info("fetching status")
				status <- ServerStatus{
					state.String(),
					fps.Rate(),
				}
			}

		default:
			if state == Running {
				buf, _ := video.Capture()
				b := make([]byte, buf.Size())
				buf.Read(b)
				img, _, err := image.Decode(bytes.NewReader(b))
				if err != nil {
					logrus.WithError(err).Error("unable to decode frame")
				}
				results := Get(img, bounds)
				l := make(map[int][]float32)
				for id, clr := range results {
					c1, c2, c3 := clr.Xyy()
					l[id] = []float32{float32(c1), float32(c2), float32(c3)}
				}
				fmt.Println(l)
				stream.Set(l)

				fps.Incr(1)
			}
		}
	}
}

type Processor struct {
	ID    int
	Color colorful.Color
}

var wg sync.WaitGroup

func Get(frame image.Image, bounds location.Bounds) map[int]colorful.Color {
	res := map[int]colorful.Color{}

	fb := frame.Bounds()
	width := fb.Max.X
	height := fb.Max.Y

	retChan := make(chan Processor, len(bounds))

	for _, b := range bounds {
		wg.Add(1)
		go func(bound location.Bound, frame image.Image, results chan Processor) {
			defer wg.Done()
			rect := bound.Rectangle(width, height)
			section := image.NewRGBA(rect)
			draw.Draw(section, rect, frame, rect.Min, draw.Src)

			//m := resize.Resize(50, 0, section, resize.Lanczos3)

			clr := extract.Average(section)
			//r, g, z, a := clr.RGBA()
			//fmt.Println(b.ID, r>>8, g>>8, z>>8, a>>8)
			res := Processor{
				ID:    bound.ID,
				Color: clr,
			}
			results <- res
		}(b, frame, retChan)
	}

	wg.Wait()
	close(retChan)

	for m := range retChan {
		res[m.ID] = m.Color
	}

	return res
}
