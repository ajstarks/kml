// kml reads KML files and produces deck/decksh markup
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// KML Structure
type Kml struct {
	XMLName  xml.Name `xml:"kml"`
	Text     string   `xml:",chardata"`
	Xsd      string   `xml:"xsd,attr"`
	Gx       string   `xml:"gx,attr"`
	Atom     string   `xml:"atom,attr"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document struct {
		Text   string `xml:",chardata"`
		LookAt struct {
			Text      string `xml:",chardata"`
			Longitude string `xml:"longitude"`
			Latitude  string `xml:"latitude"`
			Range     string `xml:"range"`
			Tilt      string `xml:"tilt"`
			Heading   string `xml:"heading"`
		} `xml:"LookAt"`
		Name       string `xml:"name"`
		Visibility string `xml:"visibility"`
		Style      struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			IconStyle struct {
				Text  string `xml:",chardata"`
				Scale string `xml:"scale"`
			} `xml:"IconStyle"`
			LabelStyle struct {
				Text  string `xml:",chardata"`
				Scale string `xml:"scale"`
			} `xml:"LabelStyle"`
			LineStyle struct {
				Text            string `xml:",chardata"`
				Color           string `xml:"color"`
				Width           string `xml:"width"`
				LabelVisibility string `xml:"labelVisibility"`
			} `xml:"LineStyle"`
			PolyStyle struct {
				Text  string `xml:",chardata"`
				Color string `xml:"color"`
			} `xml:"PolyStyle"`
		} `xml:"Style"`
		Schema struct {
			Text        string `xml:",chardata"`
			Name        string `xml:"name,attr"`
			ID          string `xml:"id,attr"`
			SimpleField []struct {
				Text        string `xml:",chardata"`
				Type        string `xml:"type,attr"`
				Name        string `xml:"name,attr"`
				DisplayName string `xml:"displayName"`
			} `xml:"SimpleField"`
		} `xml:"Schema"`
		Folder struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			Name      string `xml:"name"`
			Placemark []struct {
				Text         string `xml:",chardata"`
				ID           string `xml:"id,attr"`
				Name         string `xml:"name"`
				Visibility   string `xml:"visibility"`
				Description  string `xml:"description"`
				StyleUrl     string `xml:"styleUrl"`
				ExtendedData struct {
					Text       string `xml:",chardata"`
					SchemaData struct {
						Text       string `xml:",chardata"`
						SchemaUrl  string `xml:"schemaUrl,attr"`
						SimpleData []struct {
							Text string `xml:",chardata"`
							Name string `xml:"name,attr"`
						} `xml:"SimpleData"`
					} `xml:"SchemaData"`
				} `xml:"ExtendedData"`
				MultiGeometry struct {
					Text    string `xml:",chardata"`
					Polygon []struct {
						Text            string `xml:",chardata"`
						Extrude         string `xml:"extrude"`
						Tessellate      string `xml:"tessellate"`
						AltitudeMode    string `xml:"altitudeMode"`
						OuterBoundaryIs struct {
							Text       string `xml:",chardata"`
							LinearRing struct {
								Text        string `xml:",chardata"`
								Coordinates string `xml:"coordinates"`
							} `xml:"LinearRing"`
						} `xml:"outerBoundaryIs"`
					} `xml:"Polygon"`
				} `xml:"MultiGeometry"`
				Polygon struct {
					Text            string `xml:",chardata"`
					Extrude         string `xml:"extrude"`
					Tessellate      string `xml:"tessellate"`
					AltitudeMode    string `xml:"altitudeMode"`
					OuterBoundaryIs struct {
						Text       string `xml:",chardata"`
						LinearRing struct {
							Text        string `xml:",chardata"`
							Coordinates string `xml:"coordinates"`
						} `xml:"LinearRing"`
					} `xml:"outerBoundaryIs"`
				} `xml:"Polygon"`
			} `xml:"Placemark"`
		} `xml:"Folder"`
	} `xml:"Document"`
}

const (
	linefmt = "<line xp1=\"%.3f\" yp1=\"%.3f\" xp2=\"%.3f\" yp2=\"%.3f\" sp=\"%.3f\" color=\"%s\"/>\n"
	textfmt = "<text align=\"c\" sp=\"1.2\" xp=\"%.3f\" yp=\"%.3f\">(%.2f, %.2f)</text>\n"
	rectfmt = "<rect xp=\"%.3f\" yp=\"%.3f\" wp=\"%.3f\" hp=\"%.3f\" color=\"%s\" opacity=\"10\"/>\n"
)

// geometry defines the canvas and map boundaries
type geometry struct {
	xmin, xmax       float64
	ymin, ymax       float64
	latmin, latmax   float64
	longmin, longmax float64
}

// readData loads the KML structure from an io.Reader
func readData(r io.Reader) (Kml, error) {
	var data Kml
	err := xml.NewDecoder(r).Decode(&data)
	return data, err
}

// parseCoords makes x, y slices from the string data contained in the kml coordinate element
// (lat,long,elevation separated by commas, each coordinate separated by spaces)
func parseCoords(s string, g geometry) ([]float64, []float64) {
	f := strings.Fields(s)
	n := len(f)
	x := make([]float64, n)
	y := make([]float64, n)
	for i, c := range f {
		coords := strings.Split(c, ",")
		x[i], _ = strconv.ParseFloat(coords[0], 64)
		y[i], _ = strconv.ParseFloat(coords[1], 64)
		x[i] = vmap(x[i], g.longmin, g.longmax, g.xmin, g.xmax)
		y[i] = vmap(y[i], g.latmin, g.latmax, g.ymin, g.ymax)
	}
	return x, y
}

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// poly makes decksh markup for a polygon or polyline given x, y slices
func poly(name string, x, y []float64, lw float64, color string) {
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

// deckpolyline makes deck markup for a ployline given x, y coordinate slices
func deckpolyline(x, y []float64, lw float64, color string) {
	lx := len(x)
	if lx < 2 {
		return
	}
	for i := 0; i < lx-1; i++ {
		fmt.Printf(linefmt, x[i], y[i], x[i+1], y[i+1], lw, color)
	}
	fmt.Printf(linefmt, x[0], y[0], x[lx-1], y[lx-1], lw, color)
}

// deckpolygon makes deck markup for a polygon given x, y coordinates slices
func deckpolygon(x, y []float64, lw float64, color string) {
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

// deckshape makes either a set of polylines or polygons given a slice of coordinates
func deckshape(name string, x, y []float64, lw float64, color string) {
	switch name {
	case "polyline", "line":
		deckpolyline(x, y, lw, color)
	case "polygon", "fill":
		deckpolygon(x, y, lw, color)
	default:
		deckpolyline(x, y, lw, color)
	}
}

// deckbegin begins a deck
func deckbegin(bgcolor string) {
	if bgcolor == "" {
		fmt.Printf("<deck><slide>")
	} else {
		fmt.Printf("<deck><slide bg=\"%s\">", bgcolor)
	}
}

// deckend ends a deck
func deckend() {
	fmt.Printf("</slide></deck>")
}

// boundingBox makes a lat/long bounding box, labeled at the corners
func boundingBox(g geometry, color string) {
	w := g.xmax - g.xmin
	h := g.ymax - g.ymin
	x := g.xmin + (w / 2)
	y := g.ymin + (h / 2)
	fmt.Printf(textfmt, g.xmin, g.ymin, g.longmin, g.latmin) // lower left
	fmt.Printf(textfmt, g.xmax, g.ymin, g.longmax, g.latmin) // lower right
	fmt.Printf(textfmt, g.xmax, g.ymax, g.longmax, g.latmax) // upper right
	fmt.Printf(textfmt, g.xmin, g.ymax, g.longmin, g.latmax) // upper right
	fmt.Printf(rectfmt, x, y, w, h, color)
}

func main() {

	var mapgeo geometry
	var fulldeck bool
	var linewidth float64
	var color, bbox, shape, bgcolor string

	// options
	flag.Float64Var(&mapgeo.xmin, "xmin", 5, "canvas x minimum")
	flag.Float64Var(&mapgeo.xmax, "xmax", 95, "canvas x maxmum")
	flag.Float64Var(&mapgeo.ymin, "ymin", 10, "canvas y minimum")
	flag.Float64Var(&mapgeo.ymax, "ymax", 80, "canvas y maximum")
	flag.Float64Var(&mapgeo.latmin, "latmin", 24, "latitude x minimum")
	flag.Float64Var(&mapgeo.latmax, "latmax", 50, "latitude x maxmum")
	flag.Float64Var(&mapgeo.longmin, "longmin", -125, "longitude y minimum")
	flag.Float64Var(&mapgeo.longmax, "longmax", -67, "longitude y maximum")
	flag.Float64Var(&linewidth, "linewidth", 0.1, "line width")
	flag.StringVar(&color, "color", "black", "line color")
	flag.StringVar(&bbox, "bbox", "", "bounding box color (\"\" no box)")
	flag.StringVar(&shape, "shape", "polyline", "polygon or polyline")
	flag.StringVar(&bgcolor, "bgcolor", "", "background color")
	flag.BoolVar(&fulldeck, "fulldeck", true, "make a full deck")
	flag.Parse()

	// read data
	data, err := readData(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// add deck/slide markup, if specified
	if fulldeck {
		deckbegin(bgcolor)
	}
	// make a bounding box, if specified
	if len(bbox) > 0 {
		boundingBox(mapgeo, bbox)
	}
	// for every placemark, get the coordinates of the polygons
	for _, pms := range data.Document.Folder.Placemark {
		px, py := parseCoords(pms.Polygon.OuterBoundaryIs.LinearRing.Coordinates, mapgeo) // single polygons
		deckshape(shape, px, py, linewidth, color)
		mpolys := pms.MultiGeometry.Polygon // multiple polygons
		for _, p := range mpolys {
			mx, my := parseCoords(p.OuterBoundaryIs.LinearRing.Coordinates, mapgeo)
			deckshape(shape, mx, my, linewidth, color)
		}
	}
	// end the deck, if specified
	if fulldeck {
		deckend()
	}
}
