package main

import (
	"context"
	"fmt"
	"new-robot-self/global"
	"new-robot-self/internal/awareness"
	loc "new-robot-self/pkg/location"
	"new-robot-self/pkg/location/provider/amap"
	"os"
	"time"

	"new-robot-self/internal/brain"
)

func main() {
	ctx := context.Background()
	newBrain := brain.NewBrain()

	amapKey := global.AmapKey
	if amapKey == "" {
		fmt.Println("AMAP_KEY env is empty; please export your AMap Web service key")
		os.Exit(1)
	}

	cli, err := amap.NewClient(amapKey, 5000)
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

	gps_service, err := brain.NewGPSService(ctx, newBrain, svc)
	if err != nil {
		fmt.Println("new gps service error:", err)
		return
	}
	gps_service.Start()
	defer gps_service.Stop()
	gps_service.ReportRaw("erkqll13", loc.RawLocation{
		TS:       time.Now().Unix(),
		LatWGS84: 34.19422779316223,
		LonWGS84: 117.22516690979421,
		DeviceID: "asldallr",
	})
	time.Sleep(time.Second * 5)
	fmt.Println("after brain:", newBrain)
}
