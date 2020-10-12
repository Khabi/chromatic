package location

import (
	"image"
)

// Preset options for sampling
const (
	Top = iota
	Bottom
	Left
	Right
	Whole
)

// Default length and widths for presets.
const (
	borderThickness = 5
	borderLength    = 100
)

type Bounds []Bound

// Bound determine a place on the screen to sample for colors
// X/Y center point falls on a grid of -1 to 1
// Width and Height are percentages of the screen to draw a box
// around the centerpoint.
//          1
//          |
//          |
//  -1 -----+----- 1
//          |
//          |
//          -1
type Bound struct {
	ID     int     // identifier for the location
	X      float64 // x-axis centerpoint between -1 and 1
	Y      float64 // y-axis centerpoint between -1 and 1
	Width  int     // width in % of the rectangle around the centerpoint
	Height int     // height in % of the rectangle around the centerpoint
}

// Preset is a helper to create some common bounds.
func Preset(id int, preset int) Bound {
	switch preset {
	case Top:
		return Bound{ID: id, X: 0, Y: 1, Width: borderLength, Height: borderThickness}
	case Bottom:
		return Bound{ID: id, X: 0, Y: -1, Width: borderLength, Height: borderThickness}
	case Left:
		return Bound{ID: id, X: -1, Y: 0, Width: borderThickness, Height: borderLength}
	case Right:
		return Bound{ID: id, X: 1, Y: 0, Width: borderThickness, Height: borderLength}
	case Whole:
		return Bound{ID: id, X: 0, Y: 0, Width: borderLength, Height: borderLength}

	default:
		return Bound{}
	}
}

// Rectangle takes a bounding box and converts it to a rectangle.
// If the bounding box center places any part of the box outside
// the given width and height, the box is adjusted over to fix inside.
func (b Bound) Rectangle(width int, height int) image.Rectangle {
	var tl image.Point // topleft of rectangle
	var br image.Point // bottom right of rectangle

	boxWidth := (width * b.Width) / 100
	boxHeight := (height * b.Height) / 100

	center := b.CenterPoint(width, height)

	// Set the X axis
	switch {
	case center.X-(boxWidth/2) < 0:
		tl.X = 0
		br.X = boxWidth
	case center.X+int(boxWidth/2) > width:
		tl.X = width - boxWidth
		br.X = width
	default:
		tl.X = center.X - int(boxWidth/2)
		br.X = center.X + int(boxWidth/2)
	}

	// Set the Y axis
	switch {
	case center.Y-(boxHeight/2) < 0:
		tl.Y = 0
		br.Y = boxHeight
	case center.Y+(boxHeight/2) > height:
		tl.Y = height - boxHeight
		br.Y = height
	default:
		tl.Y = center.Y - (boxHeight / 2)
		br.Y = center.Y + (boxHeight / 2)
	}

	return image.Rectangle{tl, br}
}

// CenterToPoint takes the center and converts it to a point location
// in a box with the given width and height.
func (b Bound) CenterPoint(width int, height int) image.Point {
	x := int((b.X + 1) * (float64(width) / 2))
	y := int(((-1 * b.Y) + 1) * (float64(height) / 2))
	return image.Point{x, y}
}
