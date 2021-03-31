package main

import (
	"log"
	"os"

	"github.com/tamarakaufler/go-and-reflect/env"
)

type User struct {
	Name    string  `env:"USER_NAME" envDefault:"Lucien"`
	Age     float32 `env:"USER_AGE" envDefault:"23.5"`
	Address Address

	nationalInsurance string
}

type Address struct {
	Street   string `env:"USER_ADDRESS_STREET,required"`
	City     string `env:"USER_ADDRESS_CITY,required"`
	Postcode string `env:"USER_ADDRESS_POSTCODE,required"`
	LatLng   LatLng
}

type User2 struct {
	Name    string  `json:"user_name"`
	Age     float32 `json:"user_age"`
	Note    []rune  `json:"note"`
	Address Address2

	nationalInsurance string
}

type Address2 struct {
	Street   []rune `json:"user_address_street"`
	City     []rune `json:"user_address_city"`
	Postcode string `json:"user_address_postcode"`
	LatLng   LatLng
}

type LatLng struct {
	Lat float64 `env:"USER_ADDRESS_LAT" envDefault:"40.0000"`
	Lng float64 `env:"USER_ADDRESS_LNG" envDefault:"-115.1111"`
}

type User3 struct {
	Name    string  `json:"user_name"`
	Age     float32 `json:"user_age"`
	Note    []rune  `json:"note"`
	Address Address2

	nationalInsurance string
}

// type tagInfo struct {
// 	name      string
// 	omitempty bool
// 	omit      bool
// }

func main() {
	os.Setenv("USER_NAME", "Rebecca")
	os.Setenv("USER_ADDRESS_STREET", "16 St Mary's Close")
	os.Setenv("USER_ADDRESS_CITY", "St Albans")
	os.Setenv("USER_ADDRESS_POSTCODE", "AL3")
	os.Setenv("USER_AGE", "45")

	log.Println("######################### env ############################")

	cfg := &User{}

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(" After parsing: cfg ... %+v\n", cfg)
}
