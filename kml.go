// kmldeck reads KML files and produces deck/decksh markup
package kml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	linefmt       = "<line xp1=\"%.3f\" yp1=\"%.3f\" xp2=\"%.3f\" yp2=\"%.3f\" sp=\"%.3f\" color=\"%s\"/>\n"
	textfmt       = "<text align=\"c\" sp=\"1.2\" xp=\"%.3f\" yp=\"%.3f\">(%.2f, %.2f)</text>\n"
	rectfmt       = "<rect xp=\"%.3f\" yp=\"%.3f\" wp=\"%.3f\" hp=\"%.3f\" color=\"%s\" opacity=\"10\"/>\n"
	deckshlinefmt = "line %.3f %.3f %.3f %.3f %.2f \"%s\"\n"
)

// geometry defines the canvas and map boundaries
type Geometry struct {
	Xmin, Xmax       float64
	Ymin, Ymax       float64
	Latmin, Latmax   float64
	Longmin, Longmax float64
}

// parseCoords makes x, y slices from the string data contained in the kml coordinate element
// (lat,long,elevation separated by commas, each coordinate separated by spaces)
func ParseCoords(s string, g Geometry) ([]float64, []float64) {
	f := strings.Fields(s)
	n := len(f)
	x := make([]float64, n)
	y := make([]float64, n)
	var xp, yp float64
	for i, c := range f {
		coords := strings.Split(c, ",")
		xp, _ = strconv.ParseFloat(coords[0], 64)
		yp, _ = strconv.ParseFloat(coords[1], 64)
		x[i] = vmap(xp, g.Longmin, g.Longmax, g.Xmin, g.Xmax)
		y[i] = vmap(yp, g.Latmin, g.Latmax, g.Ymin, g.Ymax)
	}
	return x, y
}

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func filter(x, y []float64, g Geometry) ([]float64, []float64) {
	nc := len(x)
	if nc != len(y) {
		return x, y
	}
	xp := []float64{}
	yp := []float64{}
	for i := 0; i < nc; i++ {
		if x[i] >= g.Xmin && x[i] <= g.Xmax && y[i] >= g.Ymin && y[i] <= g.Ymax {
			xp = append(xp, x[i])
			yp = append(yp, y[i])
		}
	}
	return xp, yp
}

// Deckpolygon makes deck markup for a polygon given x, y coordinates slices
func Deckpolygon(x, y []float64, color string, g Geometry) {
	nc := len(x)
	if nc < 3 || nc != len(x) {
		return
	}
	end := nc - 1
	fmt.Printf("<polygon color=\"%s\" xc=\"%.3f", color, x[0])
	for i := 1; i < nc; i++ {
		fmt.Printf(" %.3f", x[i])
	}
	fmt.Printf(" %.3f\" ", x[end])
	fmt.Printf("yc=\"%.3f", y[0])
	for i := 1; i < nc; i++ {
		fmt.Printf(" %.3f", y[i])
	}
	fmt.Printf(" %.3f\"/>\n", y[end])
}

// Deckshpoly makes decksh markup for a polygon or polyline given x, y slices
func Deckshpolygon(x, y []float64, color string, g Geometry) {
	nc := len(x)
	if nc < 3 || nc != len(y) {
		return
	}
	end := nc - 1
	fmt.Printf("polygon \"%.3f", x[0])
	for i := 1; i < len(x); i++ {
		fmt.Printf(" %.3f", x[i])
	}
	fmt.Printf(" %.3f\" ", x[end])

	fmt.Printf(" \"%.3f", y[0])
	for i := 1; i < len(y); i++ {
		fmt.Printf(" %.3f", y[i])
	}
	fmt.Printf(" %.3f\" %s\n", y[end], color)
}

// Deckpolyline makes deck markup for a ployline given x, y coordinate slices
func Deckpolyline(x, y []float64, lw float64, color string, g Geometry) {
	lx := len(x)
	if lx < 2 {
		return
	}
	for i := 0; i < lx-1; i++ {
		deckline(x[i], y[i], x[i+1], y[i+1], lw, color, g)
	}
	deckline(x[0], y[0], x[lx-1], y[lx-1], lw, color, g)
}

// Deckshpolyline makes decksh markup for a polyline given x, y coordinate slices
func Deckshpolyline(x, y []float64, lw float64, color string, g Geometry) {
	lx := len(x)
	if lx < 2 {
		return
	}
	for i := 0; i < lx-1; i++ {
		deckshline(x[i], y[i], x[i+1], y[i+1], lw, color, g)
	}
	deckshline(x[0], y[0], x[lx-1], y[lx-1], lw, color, g)
}

// deckline makes a line in deck markup
func deckline(x1, y1, x2, y2, lw float64, color string, g Geometry) {
	if x1 >= g.Xmin && x2 <= g.Xmax && y1 >= g.Ymin && y2 <= g.Ymax {
		fmt.Printf(linefmt, x1, y1, x2, y2, lw, color)
	}
}

// deckshline makes a line in decksh markup
func deckshline(x1, y1, x2, y2, lw float64, color string, g Geometry) {
	if x1 >= g.Xmin && x2 <= g.Xmax && y1 >= g.Ymin && y2 <= g.Ymax {
		fmt.Printf(deckshlinefmt, x1, y1, x2, y2, lw, color)
	}
}

// Deckshape makes either a set of polylines or polygons given a slice of coordinates
func Deckshape(name string, x, y []float64, lw float64, color string, g Geometry) {
	switch name {
	case "polyline", "line":
		Deckpolyline(x, y, lw, color, g)
	case "dline":
		Deckshpolyline(x, y, lw, color, g)
	case "polygon", "fill":
		Deckpolygon(x, y, color, g)
	case "dpoly":
		Deckshpolygon(x, y, color, g)
	}
}

// Deckbegin begins a deck
func Deckbegin(bgcolor string) {
	if bgcolor == "" {
		fmt.Printf("<deck><slide>")
	} else {
		fmt.Printf("<deck><slide bg=\"%s\">", bgcolor)
	}
}

// Deckend ends a deck
func Deckend() {
	fmt.Printf("</slide></deck>")
}

// BoundingBox makes a lat/long bounding box, labeled at the corners
func BoundingBox(g Geometry, color string) {
	w := g.Xmax - g.Xmin
	h := g.Ymax - g.Ymin
	x := g.Xmin + (w / 2)
	y := g.Ymin + (h / 2)
	fmt.Printf(textfmt, g.Xmin, g.Ymin, g.Longmin, g.Latmin) // lower left
	fmt.Printf(textfmt, g.Xmax, g.Ymin, g.Longmax, g.Latmin) // lower right
	fmt.Printf(textfmt, g.Xmax, g.Ymax, g.Longmax, g.Latmax) // upper right
	fmt.Printf(textfmt, g.Xmin, g.Ymax, g.Longmin, g.Latmax) // upper right
	fmt.Printf(rectfmt, x, y, w, h, color)
}
