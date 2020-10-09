package extract

import (
	"image"
	"image/color"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// Average returns the average color of an image.
// This is most effective in small images.  The bigger the
// image, the more likely its going to move towards black
// or brown.
func Average(i image.Image) colorful.Color {
	var r, g, b, a uint32     // RGBA values
	var pr, pg, pb, pa uint32 // RGBA values at a specific point

	bounds := i.Bounds()
	pixels := uint32(bounds.Dy() * bounds.Dx())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pr, pg, pb, pa = i.At(x, y).RGBA()
			r += pr
			g += pg
			b += pb
			a += pa
		}
	}

	r /= pixels
	g /= pixels
	b /= pixels
	a /= pixels

	avgColor, _ := colorful.MakeColor(
		color.NRGBA{
			uint8(r / 0x101),
			uint8(g / 0x101),
			uint8(b / 0x101),
			uint8(a / 0x101),
		},
	)
	return avgColor
}

// ByCount orders the cluster from most observations to last
// cluster[0] will always be the most prominent.
type ByCount []clusters.Cluster

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return len(a[i].Observations) > len(a[j].Observations) }

// Prominent tries to pull the most prominent color out of an
// image by using kmeans.  The bigger the image given here, the longer
// it will take to process.  Images should be downscaled to something
// like 50x50 to get a speedy response
func Prominent(i image.Image) colorful.Color {
	var o clusters.Observations
	var r, g, b, a uint32     // RGBA values
	var pr, pg, pb, pa uint32 // RGBA values at a specific point

	bounds := i.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pr, pg, pb, pa = i.At(x, y).RGBA()
			o = append(o, clusters.Coordinates{
				float64(pr),
				float64(pg),
				float64(pb),
				float64(pa),
			})
		}
	}

	km, _ := kmeans.NewWithOptions(.001, nil)
	clusters, _ := km.Partition(o, 2)

	sort.Sort(ByCount(clusters))
	r = uint32(clusters[0].Center[0])
	g = uint32(clusters[0].Center[1])
	b = uint32(clusters[0].Center[2])
	a = uint32(clusters[0].Center[3])

	prominentColor, _ := colorful.MakeColor(
		color.NRGBA{
			uint8(r / 0x101),
			uint8(g / 0x101),
			uint8(b / 0x101),
			uint8(a / 0x101),
		},
	)
	return prominentColor
}
