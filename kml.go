// kmldeck reads KML files and produces deck/decksh markup
package kml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	linefmt = "<line xp1=\"%.3f\" yp1=\"%.3f\" xp2=\"%.3f\" yp2=\"%.3f\" sp=\"%.3f\" color=\"%s\"/>\n"
	textfmt = "<text align=\"c\" sp=\"1.2\" xp=\"%.3f\" yp=\"%.3f\">(%.2f, %.2f)</text>\n"
	rectfmt = "<rect xp=\"%.3f\" yp=\"%.3f\" wp=\"%.3f\" hp=\"%.3f\" color=\"%s\" opacity=\"10\"/>\n"
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
	for i, c := range f {
		coords := strings.Split(c, ",")
		x[i], _ = strconv.ParseFloat(coords[0], 64)
		y[i], _ = strconv.ParseFloat(coords[1], 64)
		x[i] = vmap(x[i], g.Longmin, g.Longmax, g.Xmin, g.Xmax)
		y[i] = vmap(y[i], g.Latmin, g.Latmax, g.Ymin, g.Ymax)
	}
	return x, y
}

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// Poly makes decksh markup for a polygon or polyline given x, y slices
func Poly(name string, x, y []float64, lw float64, color string) {
	style := fmt.Sprintf("%.2f %s", lw, color)
	fmt.Printf("%s \"%.3f", name, x[0])
	for i := 1; i < len(x); i++ {
		fmt.Printf(" %.3f", x[i])
	}
	fmt.Printf(" %.3f\" ", x[len(x)-1])

	fmt.Printf(" \"%.3f", y[0])
	for i := 1; i < len(y); i++ {
		fmt.Printf(" %.3f", y[i])
	}
	fmt.Printf(" %.3f\" %s\n", y[len(y)-1], style)
}

// Deckpolyline makes deck markup for a ployline given x, y coordinate slices
func Deckpolyline(x, y []float64, lw float64, color string) {
	lx := len(x)
	if lx < 2 {
		return
	}
	for i := 0; i < lx-1; i++ {
		//fmt.Printf(linefmt, x[i], y[i], x[i+1], y[i+1], lw, color)
		line(x[i], y[i], x[i+1], y[i+1], lw, color)
	}
	//fmt.Printf(linefmt, x[0], y[0], x[lx-1], y[lx-1], lw, color)
	line(x[0], y[0], x[lx-1], y[lx-1], lw, color)
}

func line(x1, y1, x2, y2, lw float64, color string) {
	if x1 < 0.0 || x2 < 0.0 || x1 > 100 || x2 > 100 {
		return
	}
	if y1 < 0.0 || y2 < 0.0 || y1 > 100 || y2 > 100 {
		return
	}
	fmt.Printf(linefmt, x1, y1, x2, y2, lw, color)
}

// Deckpolygon makes deck markup for a polygon given x, y coordinates slices
func Deckpolygon(x, y []float64, lw float64, color string) {
	nc := len(x)
	if nc < 3 || nc != len(y) {
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

// Deckshape makes either a set of polylines or polygons given a slice of coordinates
func Deckshape(name string, x, y []float64, lw float64, color string) {
	switch name {
	case "polyline", "line":
		Deckpolyline(x, y, lw, color)
	case "polygon", "fill":
		Deckpolygon(x, y, lw, color)
	default:
		Deckpolyline(x, y, lw, color)
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
