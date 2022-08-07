package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/ajstarks/kml"
)

// KML Structure
type Kml struct {
	XMLName  xml.Name `xml:"kml"`
	Text     string   `xml:",chardata"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document struct {
		Text   string `xml:",chardata"`
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
		Placemark []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name"`
			Polygon struct {
				Text            string `xml:",chardata"`
				OuterBoundaryIs struct {
					Text       string `xml:",chardata"`
					LinearRing struct {
						Text        string `xml:",chardata"`
						Coordinates string `xml:"coordinates"`
					} `xml:"LinearRing"`
				} `xml:"outerBoundaryIs"`
				InnerBoundaryIs struct {
					Text       string `xml:",chardata"`
					LinearRing struct {
						Text        string `xml:",chardata"`
						Coordinates string `xml:"coordinates"`
					} `xml:"LinearRing"`
				} `xml:"innerBoundaryIs"`
			} `xml:"Polygon"`
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
					OuterBoundaryIs struct {
						Text       string `xml:",chardata"`
						LinearRing struct {
							Text        string `xml:",chardata"`
							Coordinates string `xml:"coordinates"`
						} `xml:"LinearRing"`
					} `xml:"outerBoundaryIs"`
					InnerBoundaryIs []struct {
						Text       string `xml:",chardata"`
						LinearRing struct {
							Text        string `xml:",chardata"`
							Coordinates string `xml:"coordinates"`
						} `xml:"LinearRing"`
					} `xml:"innerBoundaryIs"`
				} `xml:"Polygon"`
			} `xml:"MultiGeometry"`
		} `xml:"Placemark"`
	} `xml:"Document"`
}

// readData loads the KML structure from an io.Reader
func readData(filename string) (Kml, error) {
	var data Kml
	r, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	err = xml.NewDecoder(r).Decode(&data)
	r.Close()
	return data, err
}

// kmldeck makes deck or decksh markup from coordinates
func kmldeck(data Kml, m kml.Geometry, linewidth float64, color, shape, style string) {
	// for every placemark, get the coordinates of the polygons
	for _, pms := range data.Document.Placemark {
		px, py := kml.ParseCoords(pms.Polygon.OuterBoundaryIs.LinearRing.Coordinates, m) // single polygons
		kml.Deckshape(shape, style, px, py, linewidth, color, m)
		mpolys := pms.MultiGeometry.Polygon // multiple polygons
		for _, p := range mpolys {
			mx, my := kml.ParseCoords(p.OuterBoundaryIs.LinearRing.Coordinates, m)
			kml.Deckshape(shape, style, mx, my, linewidth, color, m)
		}
	}
}

// kmldump prints coordinates contained in a KML document
func kmldump(data Kml) {
	// for every placemark, get the coordinates of the polygons
	for _, pms := range data.Document.Placemark {
		px, py := kml.ParsePlainCoords(pms.Polygon.OuterBoundaryIs.LinearRing.Coordinates) // single polygons
		kml.DumpCoords(px, py)
		mpolys := pms.MultiGeometry.Polygon // multiple polygons
		for _, p := range mpolys {
			mx, my := kml.ParsePlainCoords(p.OuterBoundaryIs.LinearRing.Coordinates)
			kml.DumpCoords(mx, my)
		}
	}
}

// begin begins a deck or decksh document
func begin(style, color string) {
	switch style {
	case "deck":
		kml.Deckbegin(color)
	case "decksh":
		kml.Deckshbegin(color)
	}
}

// end ends a deck or decksh document
func end(style string) {
	switch style {
	case "deck":
		kml.Deckend()
	case "decksh":
		kml.Deckshend()
	}
}

func main() {

	var mapgeo kml.Geometry
	var fulldeck bool
	var linewidth float64
	var color, bbox, shape, bgcolor, style string

	// options
	flag.Float64Var(&mapgeo.Xmin, "xmin", 5, "canvas x minimum")
	flag.Float64Var(&mapgeo.Xmax, "xmax", 95, "canvas x maxmum")
	flag.Float64Var(&mapgeo.Ymin, "ymin", 5, "canvas y minimum")
	flag.Float64Var(&mapgeo.Ymax, "ymax", 95, "canvas y maximum")
	flag.Float64Var(&mapgeo.Latmin, "latmin", -90, "latitude x minimum")
	flag.Float64Var(&mapgeo.Latmax, "latmax", 90, "latitude x maxmum")
	flag.Float64Var(&mapgeo.Longmin, "longmin", -180, "longitude y minimum")
	flag.Float64Var(&mapgeo.Longmax, "longmax", 180, "longitude y maximum")
	flag.Float64Var(&linewidth, "linewidth", 0.1, "line width")
	flag.StringVar(&color, "color", "black", "line color")
	flag.StringVar(&bbox, "bbox", "", "bounding box color (\"\" no box)")
	flag.StringVar(&shape, "shape", "polyline", "polygon, polyline")
	flag.StringVar(&style, "style", "deck", "deck, decksh, plain")
	flag.StringVar(&bgcolor, "bgcolor", "", "background color")
	flag.BoolVar(&fulldeck, "fulldeck", true, "make a full deck")
	flag.Parse()

	// add deck/slide markup, if specified
	if fulldeck {
		begin(style, bgcolor)
	}
	for _, filename := range flag.Args() {
		// read data
		data, err := readData(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		// make a bounding box, if specified
		if len(bbox) > 0 {
			kml.BoundingBox(mapgeo, bbox, style)
		}
		switch style {
		case "deck", "decksh":
			kmldeck(data, mapgeo, linewidth, color, shape, style)
		case "plain", "dump":
			kmldump(data)
		}
	}
	// end the deck, if specified
	if fulldeck {
		end(style)
	}
}
