package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ajstarks/kml"
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

// readData loads the KML structure from an io.Reader
func readData(r io.Reader) (Kml, error) {
	var data Kml
	err := xml.NewDecoder(r).Decode(&data)
	return data, err
}

func main() {

	var mapgeo kml.Geometry
	var fulldeck bool
	var linewidth float64
	var color, bbox, shape, bgcolor string

	// options
	flag.Float64Var(&mapgeo.Xmin, "xmin", 5, "canvas x minimum")
	flag.Float64Var(&mapgeo.Xmax, "xmax", 95, "canvas x maxmum")
	flag.Float64Var(&mapgeo.Ymin, "ymin", 10, "canvas y minimum")
	flag.Float64Var(&mapgeo.Ymax, "ymax", 80, "canvas y maximum")
	flag.Float64Var(&mapgeo.Latmin, "latmin", 24, "latitude x minimum")
	flag.Float64Var(&mapgeo.Latmax, "latmax", 50, "latitude x maxmum")
	flag.Float64Var(&mapgeo.Longmin, "longmin", -125, "longitude y minimum")
	flag.Float64Var(&mapgeo.Longmax, "longmax", -67, "longitude y maximum")
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
		kml.Deckbegin(bgcolor)
	}
	// make a bounding box, if specified
	if len(bbox) > 0 {
		kml.BoundingBox(mapgeo, bbox)
	}
	// for every placemark, get the coordinates of the polygons
	for _, pms := range data.Document.Folder.Placemark {
		px, py := kml.ParseCoords(pms.Polygon.OuterBoundaryIs.LinearRing.Coordinates, mapgeo) // single polygons
		kml.Deckshape(shape, px, py, linewidth, color)
		mpolys := pms.MultiGeometry.Polygon // multiple polygons
		for _, p := range mpolys {
			mx, my := kml.ParseCoords(p.OuterBoundaryIs.LinearRing.Coordinates, mapgeo)
			kml.Deckshape(shape, mx, my, linewidth, color)
		}
	}
	// end the deck, if specified
	if fulldeck {
		kml.Deckend()
	}
}