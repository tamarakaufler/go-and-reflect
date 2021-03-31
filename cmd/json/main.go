package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	// "github.com/tamarakaufler/go-and-reflect/env"
	// refl "github.com/tamarakaufler/go-and-reflect/reflect"
)

type User2 struct {
	Name    string   `json:"user_name"`
	Age     float32  `json:"user_age"`
	Note    []rune   `json:"note"`
	NI      []int32  `json:"ni"`
	Address Address2 `json:"address"`

	nationalInsurance string
}

type Address2 struct {
	Street   []rune `json:"user_address_street"`
	City     []rune `json:"user_address_city"`
	Postcode string `json:"user_address_postcode"`
	LatLng   LatLng `json:"latlng"`
}

type LatLng struct {
	Lat float64 `env:"USER_ADDRESS_LAT" envDefault:"40.0000" json:"lat"`
	Lng float64 `env:"USER_ADDRESS_LNG" envDefault:"-115.1111" json:"lng"`
}

type User3 struct {
	Name    string   `json:"user_name"`
	Age     float32  `json:"user_age"`
	Note    []rune   `json:"note"`
	NI      []int32  `json:"ni"`
	Address Address2 `json:"address"`
}

type tagInfo struct {
	name      string
	omitempty bool
	omit      bool
}

func main() {
	os.Setenv("USER_NAME", "Rebecca")
	os.Setenv("USER_ADDRESS_STREET", "16 St Mary's Close")
	os.Setenv("USER_ADDRESS_CITY", "St Albans")
	os.Setenv("USER_ADDRESS_POSTCODE", "AL3")
	os.Setenv("USER_AGE", "45")

	log.Println("######################### json ############################")

	log.Printf("User2 and User3 are identical data structures. User2 has a custom marshalling/unmarshalling, User3 uses defaults.\n\n")

	u := User2{
		Name: "Amy",
		Age:  35,
		Note: []rune("Illustrator"),
		NI:   []int32{9, 8, 7, 6, 5, 4, 3, 2, 1},
		Address: Address2{
			Street:   []rune("123 Tyttenhanger"),
			City:     []rune("St Albans"),
			Postcode: "AL4",
		},
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marshalled User2 with custom marshaller: %s\n\n", string(b))

	u3 := User3{
		Name: "Lucien",
		Age:  35,
		Note: []rune("Jazz guitarist"),
		NI:   []int32{9, 8, 7, 6, 5, 4, 3, 2, 1},
		Address: Address2{
			Street:   []rune("2A Matheson Road"),
			City:     []rune("London"),
			Postcode: "W14",
		},
	}
	b, err = json.Marshal(u3)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marshalled User3 with default marshaller: %s\n\n", string(b))
}

func (u User2) MarshalJSON() ([]byte, error) {
	dst := map[string]interface{}{}
	dst, err := processMe(u, dst)
	if err != nil {
		return nil, err
	}

	log.Printf("dst to marshall: %+v\n", dst)
	return json.Marshal(&dst)
}

func (u User2) UnmarshalJSON(b []byte) error {
	log.Printf("\tUser2 to unmarshal: %+v\n", u)

	return json.Unmarshal(b, &u)
	//return unmarshalMe(u, b)
}

func processMe(src interface{}, dst map[string]interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
	log.Printf("==> Name: %s , src to marshall: %+v\n", t.Name(), src)

	var newF interface{}

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		log.Printf("ft: %+v\n", ft)
		log.Printf("fv: %+v\n", fv)
		log.Printf("kind: %+v, canAddr: %+v \n\n", fv.Kind(), fv.CanAddr())

		// skip unexported field
		if ft.PkgPath != "" {
			log.Printf("\tis unexported: %+v\n", ft)
			continue
		}

		ti := processTag(ft)
		log.Printf("\tTAG: %+v\n", ti)
		if ti.name == "" || ti.omit || ti.omitempty {
			continue
		}

		if fv.Kind() == reflect.Struct {
			log.Printf("\tSTRUCT --> vv: %+v\n\n", fv)
			fDst := map[string]interface{}{}

			newF, err := processMe(fv.Interface(), fDst)
			if err != nil {
				return nil, err
			}
			dst[ti.name] = newF
			return dst, nil
		}
		log.Printf("\tNot a struct: %+v\n", fv)

		log.Printf("* fv: %+v\n", fv)
		switch fft := fv.Interface().(type) {
		case []rune:
			log.Printf("\t** []rune: ft.Name=%s, \n", ft.Name)
			if ft.Name == "NI" {
				newF = fft
			} else {
				newF = string(fft)
			}
			log.Printf("\t* []rune: %+v (%s)\n", newF, ft.Name)
		default:
			log.Printf("* ffv: %+v\n", fft)
			newF = fft
			log.Printf("\t* default: %+v\n", newF)
		}
		dst[ti.name] = newF
	}

	return dst, nil
}

func marshalMe(src interface{}, dst map[string]interface{}) ([]byte, error) {
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
	log.Printf("==> Name: %s , src to marshall: %+v\n", t.Name(), src)

	var newF interface{}

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		log.Printf("ft: %+v\n", ft)
		log.Printf("fv: %+v\n", fv)
		log.Printf("kind: %+v, canAddr: %+v \n\n", fv.Kind(), fv.CanAddr())

		// skip unexported field
		if ft.PkgPath != "" {
			log.Printf("\tis unexported: %+v\n", ft)
			continue
		}

		ti := processTag(ft)
		log.Printf("\tTAG: %+v\n", ti)
		if ti.name == "" || ti.omit || ti.omitempty {
			continue
		}

		if fv.Kind() == reflect.Struct {
			log.Printf("\t--> vv: %+v\n\n", fv)

			newF, err := marshalMe(fv.Interface(), dst)
			if err != nil {
				return nil, err
			}
			dst[ti.name] = newF

			continue
		}
		log.Printf("\tNot a struct: %+v\n", fv)

		log.Printf("* fv: %+v\n", fv)
		switch fft := fv.Interface().(type) {
		case []rune:
			log.Printf("\t** []rune: ft.Name=%s, \n", ft.Name)
			if ft.Name == "NI" {
				newF = fft
			} else {
				newF = string(fft)
			}
			log.Printf("\t* []rune: %+v (%s)\n", newF, ft.Name)
		default:
			log.Printf("* ffv: %+v\n", fft)
			newF = fft
			log.Printf("\t* default: %+v\n", newF)
		}

		dst[ti.name] = newF

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
