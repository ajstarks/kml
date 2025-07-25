package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/kml"
)

// config: a bag of configuration options
type config struct {
	fulldeck, info, autobbox                                      bool
	shapesize, textsize                                           float64
	textcolor, color, bbox, shape, bgcolor, style, text, fieldsep string
}

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// mapData maps raw lat/long coordinates to canvas coordinates
func mapData(x, y []float64, g kml.Geometry) ([]float64, []float64) {
	for i := 0; i < len(x); i++ {
		x[i] = vmap(x[i], g.Longmin, g.Longmax, g.Xmin, g.Xmax)
		y[i] = vmap(y[i], g.Latmin, g.Latmax, g.Ymin, g.Ymax)
	}
	return x, y
}

// readData reads lat/long pairs (separated by white space) from a file
func readData(r io.Reader) ([]float64, []float64, error) {
	x := []float64{}
	y := []float64{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		f := strings.Fields(t)
		if len(f) != 2 {
			continue
		}
		xp, err := strconv.ParseFloat(f[1], 64) // long
		if err != nil {
			continue
		}
		yp, err := strconv.ParseFloat(f[0], 64) // lat
		if err != nil {
			continue
		}
		x = append(x, xp)
		y = append(y, yp)
	}
	return x, y, s.Err()
}

// readLoc reads lat/long pairs and an optional name (separated by white space) from a file
func readLoc(r io.Reader, sep byte) (kml.Locdata, error) {
	var data kml.Locdata
	s := bufio.NewScanner(r)
	ff := func(c rune) bool { return c == rune(sep) }
	for s.Scan() {
		t := s.Text()
		f := strings.FieldsFunc(t, ff)
		if len(f) < 2 { // not enough fields
			continue
		}
		xp, err := strconv.ParseFloat(f[1], 64) // long
		if err != nil {
			continue
		}
		yp, err := strconv.ParseFloat(f[0], 64) // lat
		if err != nil {
			continue
		}
		data.X = append(data.X, xp)
		data.Y = append(data.Y, yp)
		if len(f) > 2 { // if name is present
			data.Name = append(data.Name, f[2])
		}
	}
	return data, s.Err()
}

// bboxData returns minima and maxima from data
func bboxData(x, y []float64) (float64, float64, float64, float64) {
	maxx := -180.0
	minx := 180.0
	maxy := -90.0
	miny := 90.0

	for i := 0; i < len(x); i++ {
		xp, yp := x[i], y[i]
		if xp > maxx {
			maxx = xp
		}
		if xp < minx {
			minx = xp
		}
		if yp > maxy {
			maxy = yp
		}
		if yp < miny {
			miny = yp
		}
	}
	return minx, maxx, miny, maxy
}

// process input and options, making markup
func process(filename string, dest io.Writer, c config, mapgeo kml.Geometry) {

	// read from stdin by default, if specified, open a file
	r := os.Stdin
	if len(filename) > 0 {
		var rerr error
		r, rerr = os.Open(filename)
		if rerr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", rerr)
			return
		}
	}
	// read coordinates
	//x, y, err := readData(r)
	loc, err := readLoc(r, c.fieldsep[0])
	x := loc.X
	y := loc.Y
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer r.Close()

	// if specified, only show bbox info
	if c.info {
		minx, maxx, miny, maxy := bboxData(x, y)
		centerLon := minx + (maxx-minx)/2
		centerLat := miny + (maxy-miny)/2
		fmt.Fprintf(os.Stdout, "--center=%v,%v --bbox=\"%v,%v|%v,%v\" --longmin=%v --longmax=%v --latmin=%v --latmax=%v\n",
			centerLat, centerLon, maxy, minx, miny, maxx, minx, maxx, miny, maxy)
		return
	}
	// if specified adjust mapping to source data bounding box
	if c.autobbox {
		mapgeo.Longmin, mapgeo.Longmax, mapgeo.Latmin, mapgeo.Latmax = bboxData(x, y)
	}
	// add slide markup, if specified
	if c.fulldeck {
		if len(filename) > 0 {
			fmt.Fprintln(dest, "// "+filename)
		}
		beginslide(dest, c.bgcolor, c.style)
	}
	// draw a bounding box, if specified
	if len(c.bbox) > 0 {
		kml.BoundingBox(mapgeo, c.bbox, c.style)
	}
	// map to deck canvas, make the drawing
	x, y = mapData(x, y, mapgeo)
	kml.Deckshape(c.shape, c.style, x, y, c.shapesize, c.color, mapgeo)

	if len(c.text) > 0 {
		kml.DeckText(c.text, c.style, x, y, loc.Name, c.textsize, c.textcolor)
	}
	// end the slide, if specified
	if c.fulldeck {
		endslide(dest, c.style)
	}
}

// begindeck makes the beginning markup
func begindeck(dest io.Writer, style string) {
	if style == "deck" {
		fmt.Fprintln(dest, "<deck>")
	} else {
		fmt.Fprintln(dest, "deck")
	}
}

// enddeck makes the ending markup
func enddeck(dest io.Writer, style string) {
	if style == "deck" {
		fmt.Fprintln(dest, "</deck>")
	} else {
		fmt.Fprintln(dest, "edeck")
	}
}

// beginslide makes the markup for the beginning of a slide
func beginslide(dest io.Writer, bgcolor, style string) {
	if style == "deck" {
		fmt.Fprintf(dest, "<slide bg=%q>\n", bgcolor)
	} else {
		fmt.Fprintf(dest, "slide %q\n", bgcolor)
	}
}

// endslides makes the markup for the ending of a slide
func endslide(dest io.Writer, style string) {
	if style == "deck" {
		fmt.Fprintln(dest, "</slide>")
	} else {
		fmt.Fprintln(dest, "eslide")
	}
}
func main() {
	var mapgeo kml.Geometry
	var cfg config

	// coordinate mapping options
	flag.Float64Var(&mapgeo.Xmin, "xmin", 5, "canvas x minimum")
	flag.Float64Var(&mapgeo.Xmax, "xmax", 95, "canvas x maxmum")
	flag.Float64Var(&mapgeo.Ymin, "ymin", 5, "canvas y minimum")
	flag.Float64Var(&mapgeo.Ymax, "ymax", 95, "canvas y maximum")
	flag.Float64Var(&mapgeo.Latmin, "latmin", -90, "latitude x minimum")
	flag.Float64Var(&mapgeo.Latmax, "latmax", 90, "latitude x maxmum")
	flag.Float64Var(&mapgeo.Longmin, "longmin", -180, "longitude y minimum")
	flag.Float64Var(&mapgeo.Longmax, "longmax", 180, "longitude y maximum")
	// config options
	flag.BoolVar(&cfg.info, "info", false, "only report center and bounding box info")
	flag.BoolVar(&cfg.autobbox, "autobbox", true, "autoscale according to input values")
	flag.Float64Var(&cfg.shapesize, "shapesize", 0.25, "shapesize (linewidth or dotsize)")
	flag.StringVar(&cfg.color, "color", "black", "line color")
	flag.StringVar(&cfg.bbox, "bbox", "", "bounding box color (\"\" no box)")
	flag.StringVar(&cfg.shape, "shape", "polyline", "polygon (fill), polyline (line), circle (dot)")
	flag.StringVar(&cfg.style, "style", "decksh", "deck, decksh, plain")
	flag.StringVar(&cfg.text, "text", "", "use text labels")
	flag.Float64Var(&cfg.textsize, "textsize", 0.5, "textsize")
	flag.StringVar(&cfg.textcolor, "textcolor", "black", "textcolor")
	flag.StringVar(&cfg.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&cfg.fieldsep, "fs", " ", "data field separator")
	flag.BoolVar(&cfg.fulldeck, "fulldeck", false, "make a full deck")

	flag.Parse()

	dest := os.Stdout
	// don't do any generation if info only
	if cfg.info {
		cfg.fulldeck = false
	}
	// add deck markup, if specified
	if cfg.fulldeck {
		begindeck(dest, cfg.style)
	}
	// for every file (or stdin if no files are specified), make markup
	if len(flag.Args()) == 0 {
		process("", dest, cfg, mapgeo)
	} else {
		for _, filename := range flag.Args() {
			process(filename, dest, cfg, mapgeo)
		}
	}
	if cfg.fulldeck {
		enddeck(dest, cfg.style)
	}
}
