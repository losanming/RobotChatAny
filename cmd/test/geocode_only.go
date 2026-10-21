package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"new-robot-self/global"
	"os"
	"strconv"
	"strings"
	"time"

	amap "new-robot-self/pkg/location/provider/amap"
)

// Standalone geocoding tester. Run with:
//
//	go run ./cmd/test/geocode_only.go -address "北京市朝阳区阜通东大街6号" -city 110000
func main() {
	address := flag.String("address", "", "structured address, e.g. 北京市朝阳区阜通东大街6号 (required)")
	city := flag.String("city", "", "optional city limit: name/pinyin/citycode/adcode")
	timeoutMs := flag.Int("timeout_ms", 1500, "http timeout in ms")
	flag.Parse()

	if *address == "" {
		fmt.Println("-address is required")
		os.Exit(2)
	}

	key := global.AmapKey
	if key == "" {
		fmt.Println("global.AmapKey is empty; please configure your AMap key")
		os.Exit(2)
	}

	cli, err := amap.NewClient(key, *timeoutMs)
	if err != nil {
		fmt.Println("create client error:", err)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutMs+500)*time.Millisecond)
	defer cancel()
	res, err := cli.Geocode(ctx, amap.GeoParams{Address: *address, City: *city, Output: "json"})
	if err != nil {
		fmt.Println("Geocode error:", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))

	fmt.Printf("\n--- Geocode Summary ---\n")
	fmt.Printf("Provider=%s TS=%d Count=%d Address=%s City=%s\n", res.Provider, res.TS, res.Count, res.Address, res.City)
	for i, it := range res.Items {
		fmt.Printf("[%d] %s %s %s %s adcode=%s loc=%s level=%s\n", i+1, it.Country, it.Province, it.City, it.District, it.Adcode, it.Location, it.Level)
	}

	if len(res.Items) > 0 {
		loc := strings.TrimSpace(res.Items[0].Location)
		if loc != "" {
			parts := strings.Split(loc, ",")
			if len(parts) == 2 {
				lon, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				lat, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err1 == nil && err2 == nil {
					// Call reverse geocode with GCJ-02 coords (AMap returns GCJ-02)
					rctx, rcancel := context.WithTimeout(context.Background(), time.Duration(*timeoutMs+500)*time.Millisecond)
					defer rcancel()
					rres, rerr := cli.ReverseGeocode(rctx, lat, lon, amap.RegeoParams{Extensions: "base", Radius: 1000})
					if rerr != nil {
						fmt.Println("ReverseGeocode error:", rerr)
					} else {
						rb, _ := json.MarshalIndent(rres, "", "  ")
						fmt.Printf("\n--- ReverseGeocode (from first geocode location) ---\n")
						fmt.Println(string(rb))
						fmt.Printf("\nAddress=%s Province=%s City=%s District=%s Adcode=%s\n", rres.Address, rres.Province, rres.City, rres.District, rres.Adcode)
					}
				}
			}
		}
	}
}
