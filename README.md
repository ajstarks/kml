# KML

Convert KML files to deck markup

![kml-example](us-states.png)

```./kml -linewidth=0.075 -bbox=blue < cb_2018_us_states_5m.kml | pdfdeck -stdout - > states.pdf```

![kml-counties](us-counties.png)

```/kml -linewidth=0.075 -bbox=blue < cb_2018_us_county_20m.kml | pdfdeck -stdout - > counties.pdf```
## options
```
  -bbox string
    	bounding box
  -color string
    	line color (default "black")
  -fulldeck
    	make a full deck (default true)
  -latmax float
    	latitude x maxmum (default 50)
  -latmin float
    	latitude x minimum (default 24)
  -linewidth float
    	linewidth (default 0.1)
  -longmax float
    	longitude y maximum (default -67)
  -longmin float
    	longitude y minimum (default -125)
  -xmax float
    	canvas x maxmum (default 95)
  -xmin float
    	canvas x minimum (default 5)
  -ymax float
    	canvas y maximum (default 80)
  -ymin float
    	canvas y minimum (default 10)

```
