package location

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreset(t *testing.T) {
	var tests = []struct {
		ID       string
		Preset   int
		Expected Bound
	}{
		{
			"Light 1",
			Top,
			Bound{ID: "Light 1", X: 0, Y: 1, Width: borderLength, Height: borderThickness},
		},
		{
			"Light 2",
			Bottom,
			Bound{ID: "Light 2", X: 0, Y: -1, Width: borderLength, Height: borderThickness},
		},
		{
			"Light 3",
			Left,
			Bound{ID: "Light 3", X: -1, Y: 0, Width: borderThickness, Height: borderLength},
		},
		{
			"Light 4",
			Right,
			Bound{ID: "Light 4", X: 1, Y: 0, Width: borderThickness, Height: borderLength},
		},
		{
			"Light 5",
			Whole,
			Bound{ID: "Light 5", X: 0, Y: 0, Width: borderLength, Height: borderLength},
		},
	}

	for _, td := range tests {
		b := Preset(td.ID, td.Preset)
		assert.Equal(t, td.Expected, b)
	}

}

func TestCenterPoint(t *testing.T) {
	var width = 1024
	var height = 768

	var tests = []struct {
		ID       string
		X        float64
		Y        float64
		Expected image.Point
	}{
		{"light 1", -1, 1, image.Point{0, 0}},
		{"light 2", 0, 0, image.Point{512, 384}},
		{"light 3", 1, -1, image.Point{1024, 768}},
		{"light 4", .5, -.5, image.Point{768, 576}},
	}

	for _, td := range tests {
		b := Bound{ID: td.ID, X: td.X, Y: td.Y}

		assert.Equal(t, td.Expected, b.CenterPoint(width, height))
	}
}

func TestRectangle(t *testing.T) {
	var width = 1024
	var height = 768

	var tests = []struct {
		Bound    Bound
		Expected image.Rectangle
	}{
		{
			// Top left corner out of bounds
			Bound{ID: "Light 1", X: -1, Y: 1, Width: 25, Height: 25},
			image.Rect(0, 0, 256, 192),
		},
		{
			// Bottom right corner out of bounds
			Bound{ID: "Light 2", X: 1, Y: -1, Width: 25, Height: 25},
			image.Rect(768, 576, 1024, 768),
		},
		{
			// Centered in bounds
			Bound{ID: "Light 3", X: 0, Y: 0, Width: 25, Height: 25},
			image.Rect(384, 288, 640, 480),
		},
	}

	for _, td := range tests {
		b := td.Bound.Rectangle(width, height)
		assert.Equal(t, td.Expected, b)
	}
}
