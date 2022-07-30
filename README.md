# KML

Convert KML files to deck markup

![kml-example](us-states.png)

```./kml -linewidth=0.075 -bbox=blue < cb_2018_us_states_5m.kml | pdfdeck -stdout - > states.pdf```

![kml-counties](us-counties.png)

```./kml -linewidth=0.075 -bbox=blue < cb_2018_us_county_20m.kml | pdfdeck -stdout - > counties.pdf```

![kml-filled](filled.png)

```./kml -color "hsv(240,100,30)" -bbox blue  -shape fill  < cb_2018_us_nation_20m.kml | pdfdeck - > nation.pdf```

## options
```
  -bbox string
      bounding box color ("" no box)
  -color string
      line color (default "black")
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
