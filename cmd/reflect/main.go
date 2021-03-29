package main

import (
	"log"
	"os"

	"github.com/tamarakaufler/go-and-reflect/env"
)

type Configuration struct {
	Name    string  `env:"USER_NAME" defVal:"Lucien"`
	Age     float32 `env:"USER_AGE" defVal:"23.5"`
	Address Address

	nationalInsurance string
}

type Address struct {
	Street   string `env:"USER_ADDRESS_STREET,required"`
	City     string `env:"USER_ADDRESS_CITY,required"`
	Postcode string `env:"USER_ADDRESS_POSTCODE,required"`
	LatLng   LatLng
}

type LatLng struct {
	Lat float64 `env:"USER_ADDRESS_LAT" defVal:"40.0000"`
	Lng float64 `env:"USER_ADDRESS_LNG" defVal:"-115.1111"`
}

func main() {
	os.Setenv("USER_NAME", "Rebecca")
	os.Setenv("USER_ADDRESS_STREET", "16 St Mary's Close")
	os.Setenv("USER_ADDRESS_CITY", "St Albans")
	os.Setenv("USER_ADDRESS_POSTCODE", "AL3")
	os.Setenv("USER_AGE", "45")

	cfg := &Configuration{}

	// err := refl.Parse(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("=====================================================")
	// cfg = &Configuration{
	// 	Name:    "Marianne",
	// 	Age:     33,
	// 	Address: Address{},
	// }
	// err = refl.Parse(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=====================================================")

}
