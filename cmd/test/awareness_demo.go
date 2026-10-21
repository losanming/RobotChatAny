package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"new-robot-self/global"
	"os"
	"time"

	"new-robot-self/internal/awareness"
	"new-robot-self/pkg/location"
	amap "new-robot-self/pkg/location/provider/amap"
)

func main() {
	lat := flag.Float64("lat", 30.45, "WGS-84 latitude")
	lon := flag.Float64("lon", 120.66, "WGS-84 longitude")
	deviceID := flag.String("device", "dev-001", "device id")
	source := flag.String("source", "gps", "location source gps/wifi/cell/ip")
	needAround := flag.Bool("around", true, "enable place around search")
	keywords := flag.String("k", "", "around keywords")
	types := flag.String("types", "050000|060000", "around typecodes, '|' separated")
	radius := flag.Int("radius", 3000, "around radius meters")
	wext := flag.String("wext", "base", "weather extensions: base|all")
	to := flag.Int("timeout", 3000, "amap http timeout ms")
	flag.Parse()

	amapKey := global.AmapKey
	if amapKey == "" {
		fmt.Println("AMAP_KEY env is empty; please export your AMap Web service key")
		os.Exit(1)
	}

	cli, err := amap.NewClient(amapKey, *to)
	if err != nil {
		fmt.Println("new amap client error:", err)
		os.Exit(1)
	}

	svc, err := awareness.NewService(cli,
		awareness.WithDefaultAroundRadius(3000),
		awareness.WithWeatherTimeoutMs(2500),
	)
	if err != nil {
		fmt.Println("new awareness service error:", err)
		os.Exit(1)
	}

	in := awareness.Input{
		RawLoc: location.RawLocation{
			TS:       time.Now().Unix(),
			LatWGS84: *lat,
			LonWGS84: *lon,
			Source:   *source,
			DeviceID: *deviceID,
		},
		NeedAround:        *needAround,
		AroundKeywords:    *keywords,
		AroundTypes:       *types,
		AroundRadius:      *radius,
		WeatherExtensions: *wext,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := svc.Gather(ctx, in)
	if err != nil {
		fmt.Println("gather error:", err)
		os.Exit(1)
	}

	fmt.Println("===== PROMPT =====")
	fmt.Println(res.Prompt)

	fmt.Println("===== JSON =====")
	b, _ := json.MarshalIndent(res.Ctx, "", "  ")
	fmt.Println(string(b))
}
