package chromatic

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
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
					stream.Set(l)
				}
				fps.Incr(1)
			}
		}
	}
}

func Get(frame image.Image, bounds location.Bounds) map[int]colorful.Color {
	res := map[int]colorful.Color{}

	fb := frame.Bounds()
	width := fb.Max.X
	height := fb.Max.Y

	for _, b := range bounds {
		rect := b.Rectangle(width, height)
		section := image.NewRGBA(rect)
		draw.Draw(section, rect, frame, rect.Min, draw.Src)

		clr := extract.Average(section)
		fmt.Println(clr)
		res[b.ID] = clr
	}

	return res
}
