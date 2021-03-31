package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
)

type User struct {
	Name    string  `env:"USER_NAME" envDefault:"Lucien" json:"user_name"`
	Age     float32 `env:"USER_AGE" envDefault:"23.5" json:"user_age"`
	Address Address

	nationalInsurance string
}

type Address struct {
	Street   string `env:"USER_ADDRESS_STREET,required" json:"user_address_street"`
	City     string `env:"USER_ADDRESS_CITY,required" json:"user_address_city"`
	Postcode string `env:"USER_ADDRESS_POSTCODE,required" json:"user_address_postcode"`
	LatLng   LatLng
}

type LatLng struct {
	Lat float64 `env:"USER_ADDRESS_LAT" envDefault:"40.0000" json:"user_address_lat"`
	Lng float64 `env:"USER_ADDRESS_LNG" envDefault:"-115.1111" json:"user_address_lng"`
}

type User2 struct {
	Name    string  `json:"user_name"`
	Age     float32 `json:"user_age"`
	Address Address2

	nationalInsurance string
}

type Address2 struct {
	Street   []rune `json:"user_address_street"`
	City     []rune `json:"user_address_city"`
	Postcode string `json:"user_address_postcode"`
	LatLng   LatLng
}

type tagInfo struct {
	name      string
	omitempty bool
	omit      bool
}

func (u User2) MarshalJSON() ([]byte, error) {
	log.Printf("\tUser2 to marshal: %+v\n", u)

	return marshalMe(u)
}

func (u User2) UnmarshalJSON(b []byte) error {
	log.Printf("\tUser2 to unmarshal: %+v\n", u)

	return unmarshalMe(u, b)
}

func marshalMe(src interface{}) ([]byte, error) {
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)

	var dst interface{}  // will be the new adjusted Address
	var newF interface{} // the adjusted field of the adjusted Address

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		// skip unexported field
		if !fv.CanSet() {
			continue
		}

		ti := processTag(ft)
		if ti.omit || ti.omitempty {
			continue
		}

		switch fft := fv.Interface().(type) {
		case []rune:
			newF = string(fft)
		default:
			newF = fft
		}

		fv.Set(reflect.ValueOf(newF))
	}

	log.Printf("dst to marshall: %+v\n", dst)
	return json.Marshal(&dst)
}

func unmarshalMe(i interface{}, b []byte) error {

	return nil
}

// processTag processes struct field tag.
func processTag(sf reflect.StructField) tagInfo {
	ti := tagInfo{}

	// json tag does not exist.
	tag, ok := sf.Tag.Lookup("json")
	if !ok {
		ti.name = sf.Name
		return ti
	}

	// json tag - indicates omitting the field.
	if tag == "-" {
		ti.omit = true
		return ti
	}

	// json tag omitempty indicates omitting the field if the field is empty.
	parts := strings.SplitN(tag, ",", 2)
	if len(parts) > 1 {
		if parts[1] == "omitempty" {
			ti.omitempty = true
			return ti
		}
		tag = parts[0]
	}

	ti.name = tag
	return ti
}
