package initialize

import (
	"log"
	"new-robot-self/pkg/location/adcode"
	public_city "new-robot-self/public/cityCode"
)

func InitAdcode() error {
	st, err := adcode.Load("asset/cities_adcode.json")
	if err != nil {
		log.Fatalf("load adcode failed: %v", err)
		return err
	}
	public_city.CityStore = st
	return nil
}
