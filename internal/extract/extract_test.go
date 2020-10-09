package extract

import (
	"image"
	_ "image/jpeg"
	"os"
	"path"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/nfnt/resize"
	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {
	var tests = []struct {
		name     string
		image    string
		expected colorful.Color
	}{
		{
			"average red",
			"red.jpg",
			colorful.Color{
				R: 0.9882352941176471,
				G: 0.00784313725490196,
				B: 0,
			},
		},
		{
			"average blue",
			"blue.jpg",
			colorful.Color{
				R: 0.054901960784313725,
				G: 0.3411764705882353,
				B: 0.5882352941176471,
			},
		},
		{
			"average purple",
			"purple.jpg",
			colorful.Color{
				R: 0.6313725490196078,
				G: 0.27450980392156865,
				B: 0.596078431372549,
			},
		},
		{
			"average avatar",
			"avatar.jpg",
			colorful.Color{
				R: 0.16861219195849547,
				G: 0.023346303501945526,
				B: 0.011673151750972763,
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			fh, err := os.Open(path.Join("../../testdata", td.image))
			assert.NoError(t, err)
			defer fh.Close()

			i, _, err := image.Decode(fh)
			assert.NoError(t, err)

			c := Average(i)
			t.Log("Hex Code: ", c.Hex())
			assert.Equal(t, td.expected, c)

		})

	}
}

// TestProminent tests the promienent color analysis
// because it uses kmeans, the values may flux a bit
// these tests may actually fail.
func TestProminent(t *testing.T) {
	var tests = []struct {
		name     string
		image    string
		expected colorful.Color
	}{
		{
			"average red",
			"red.jpg",
			colorful.Color{
				R: 0.9882352941176471,
				G: 0.00784313725490196,
				B: 0,
			},
		},
		{
			"average blue",
			"blue.jpg",
			colorful.Color{
				R: 0.054901960784313725,
				G: 0.3411764705882353,
				B: 0.5882352941176471,
			},
		},
		{
			"average purple",
			"purple.jpg",
			colorful.Color{
				R: 0.6313725490196078,
				G: 0.27450980392156865,
				B: 0.596078431372549,
			},
		},
		{
			"average avatar",
			"avatar.jpg",
			colorful.Color{
				R: 0.35294117647058826,
				G: 0.14901960784313725,
				B: 0.13333333333333333,
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			fh, err := os.Open(path.Join("../../testdata", td.image))
			assert.NoError(t, err)
			defer fh.Close()

			i, _, err := image.Decode(fh)
			assert.NoError(t, err)

			newImage := resize.Resize(50, 0, i, resize.Lanczos3)

			c := Prominent(newImage)
			t.Log("Hex Code: ", c.Hex())
			assert.Equal(t, td.expected, c)

		})

	}
}
