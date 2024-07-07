# KML

Convert KML files to deck markup

## Functions

The package has these functions:
```
BoundingBox(g Geometry, color, style string)                                        // makes a bounding box

Deckbegin(bgcolor string)                                                           // begin deck
Deckend()                                                                           // end deck

Deckshbegin(bgcolor string)                                                         // begin deck, decksh markup
Deckshend()                                                                         // end deck, decksh markup

DeckPoint(x, y []float64, color string, lw float64)                                 // make circles, deck markup
Deckpolygon(x, y []float64, color string, g Geometry)                               // make a polygon, deck markup
Deckpolyline(x, y []float64, lw float64, color string, g Geometry)                  // make a polyline, deck markup

DeckshPoint(x, y []float64, color string, lw float64)                               // make circles, decksh markup
Deckshpolygon(x, y []float64, color string, g Geometry)                             // make polygon, decksh markup
Deckshpolyline(x, y []float64, lw float64, color string, g Geometry)                // make polyline, decksh markup

Deckshape(shape, style string, x, y []float64, lw float64, color string, g Geometry // make markup 

DumpCoords(x, y []float64)                                                          // print raw coordinates
ParseCoords(s string, g Geometry) ([]float64, []float64)                            // extract and map coordinates
ParsePlainCoords(s string) ([]float64, []float64)                                   // extract coordinates
```

There are three example clients:

## geodeck -- convert lat/long pairs to deck/decksh markup

geodeck reads space separated decimal lat/long pairs from stdin or specified files, and emits deck/decksh markup representing the path to stdout.
Typically other programs will generate the input, for example the ```fitscsvcoord``` command [reads CSV files with FIT data](https://developer.garmin.com/fit/fitcsvtool/).
Note that ```fitscsvcoord``` converts from "semicircle" units to decimal latitude and longitude.

```
java -jar $JARLOC/FitCSVTool.jar -iso8601 path.fit |
grep position_lat | 
csvread  -plain=t 4 7 10  | 
awk 'length($1) == 20 {printf "%.6f %.6f\n",  $2 * (180/2^31) , $3 * (180/2^31)}'

```

```
$ fitscsvcoord path.csv > path.coord
$ geodeck [options] path.coord > path.dsh
```

The ```--info``` option reports information on the center and bounding box of the coordinates without deck generation.
The reported options may be used in subsequent calls to geodeck or used in other tools like [```create-static-map```](https://github.com/flopp/go-staticmaps/tree/master/create-static-map)

```
$ geodeck --info path.coord
--center=40.6291415,-74.4224255 -bbox="40.636468,-74.4292|40.621815,-74.415651" --longmin=-74.4292 --longmax=-74.415651 --latmin=40.621815 --latmax=40.636468
```

By default, the bounding box is determined from the input data, to override, turn off autobbox, and specify the bounding box directly.

```
$ geodeck -autobbox=f --longmin=-74.4292 --longmax=-74.415651 --latmin=40.621815 --latmax=40.636468
```

## Options

```
Usage of geodeck:
  -autobbox
      autoscale according to input values (default true)
  -bbox string
      bounding box color ("" no box)
  -bgcolor string
      background color (default "white")
  -color string
      line color (default "black")
  -fulldeck
      make a full deck
  -info
      only report center and bounding box info
  -latmax float
      latitude x maxmum (default 90)
  -latmin float
      latitude x minimum (default -90)
  -linewidth float
      line width (default 0.1)
  -longmax float
      longitude y maximum (default 180)
  -longmin float
      longitude y minimum (default -180)
  -shape string
      polygon (fill), polyline (line), circle (dot) (default "polyline")
  -style string
      deck, decksh, plain (default "decksh")
  -xmax float
      canvas x maxmum (default 95)
  -xmin float
      canvas x minimum (default 5)
  -ymax float
      canvas y maximum (default 95)
  -ymin float
      canvas y minimum (default 5)

```


## World

![kml-world-outline](worldoutline.png)

```./world world.kml | pdfdeck -stdout  -pagesize 1600x1000 - > worldoutline.pdf```

![kml-world](world.png)

```./world  -shape=fill -bgcolor=lightblue -color=brown world.kml | pdfdeck -stdout  -pagesize 1600x1000 - > world.pdf```

![kml-zoom](slave-route.png)

```./world -latmin=-20 -latmax=35 -longmin=-100 -longmax=20 -shape=fill -bgcolor=lightsteelblue -color=sienna world.kml | pdfdeck -stdout -pagesize 1600x900 - > slave-route.pdf```

### options
```
  -bbox string
      bounding box color ("" no box)
  -bgcolor string
      background color
  -color string
      fill or line color (default "black")
      (specify opacity with name:op)
  -fulldeck
      make a full deck (default true)
  -latmax float
      latitude x maxmum (default 90)
  -latmin float
      latitude x minimum (default -90)
  -linewidth float
      line width (default 0.1)
  -longmax float
      longitude y maximum (default 180)
  -longmin float
      longitude y minimum (default -180)
  -shape string
      polygon, polyline (default "polyline")
  -style string
      deck, decksh, plain (default "deck")
  -xmax float
      canvas x maxmum (default 95)
  -xmin float
      canvas x minimum (default 5)
  -ymax float
      canvas y maximum (default 95)
  -ymin float
      canvas y minimum (default 5)

```

The included KML files are from the [opendatasoft site](https://public.opendatasoft.com/explore/dataset/world-administrative-boundaries/export/)

## usmap

![kml-example](us-states.png)

```./usmap -linewidth=0.075 -bbox=blue  cb_2018_us_state_5m.kml | pdfdeck -stdout - > states.pdf```

![kml-counties](us-counties.png)

```./usmap -linewidth=0.075 -bbox=blue  cb_2018_us_county_20m.kml | pdfdeck -stdout - > counties.pdf```

![kml-filled](filled.png)

```./usmap -color "hsv(240,100,30)" -bbox blue  -shape fill   cb_2018_us_nation_20m.kml | pdfdeck -stdout - > nation.pdf```

### options
```
  -bbox string
      bounding box color ("" no box)
  -bgcolor string
      background color
  -color string
      fill or line color (default "black")
      (specify opacity with name:op)
  -fulldeck
      make a full deck (default true)
  -latmax float
      latitude x maxmum (default 50)
  -latmin float
      latitude x minimum (default 24)
  -linewidth float
      line width (default 0.1)
  -longmax float
      longitude y maximum (default -67)
  -longmin float
      longitude y minimum (default -125)
  -shape string
      polygon or polyline (default "polyline")
  -style string
      deck, decksh, or plain (default "deck")
  -xmax float
      canvas x maxmum (default 95)
  -xmin float
      canvas x minimum (default 5)
  -ymax float
      canvas y maximum (default 80)
  -ymin float
      canvas y minimum (default 10)
```

The data in the repository is from the [US Census](https://www.census.gov/geographies/mapping-files/time-series/geo/kml-cartographic-boundary-files.html)


