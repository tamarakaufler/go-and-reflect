package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type User2 struct {
	Name    string   `json:"user_name"`
	Age     float32  `json:"user_age"`
	Note    []rune   `json:"note"`
	NI      []int32  `json:"ni"`
	Address Address2 `json:"address"`

	nationalInsurance string //nolint:structcheck,unused
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

var customFieldsM = map[string]struct{}{
	"note":                {},
	"user_address_street": {},
	"user_address_city":   {},
}

func main() {
	log.Println("######################### json ############################")
	log.Println(`User2 and User3 are identical data structures. User2 has a custom marshalling/unmarshalling,
	User3 uses defaults.`)

	log.Println("========================= custom marshal ============================")
	u2 := User2{
		Name: "Amy",
		Age:  35,
		Note: []rune("Illustrator"),
		NI:   []int32{9, 8, 7, 6, 5, 4, 3, 2, 1},
		Address: Address2{
			Street:   []rune("123 Tyttenhanger"),
			City:     []rune("St Albans"),
			Postcode: "AL4",
			LatLng: LatLng{
				Lat: 40.0000,
				Lng: -115.0000,
			},
		},
	}
	b2, err := json.Marshal(u2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marshalling User2 (custom marshalling): %+v\n\n", u2)
	log.Printf("Marshalled User2 with custom marshaller: %s\n\n", string(b2))

	log.Println("========================= default marshal ============================")
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
	b3, err := json.Marshal(u3)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marshalled User3 with default marshaller: %s\n\n", string(b3))

	log.Println("========================= custom unmarshal ============================")
	u2 = User2{}
	errUnm := json.Unmarshal(b2, &u2)
	if errUnm != nil {
		log.Fatal(errUnm)
	}
	log.Printf("Unmarshalled User2 with custom marshaller: %+v\n\n", u2)
}

func (u User2) MarshalJSON() ([]byte, error) {
	dst := map[string]interface{}{}
	dst, err := marshalMe(u, dst)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&dst)
}

func (u *User2) UnmarshalJSON(b []byte) error {
	log.Printf("\tbytes to unmarshal User2: %s\n", string(b))

	return unmarshalMe(b, u)
}

func marshalMe(src interface{}, dst map[string]interface{}) (map[string]interface{}, error) {
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

			newF, err := marshalMe(fv.Interface(), fDst)
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

func unmarshalMe(b []byte, dst interface{}) error {
	src := map[string]interface{}{}

	err := json.Unmarshal(b, &src)
	if err != nil {
		return err
	}

	log.Printf("unmarshaled bytes into array: %+v\n\n", src)

	src = processMap(src)

	msCfg := &mapstructure.DecoderConfig{
		// DecodeHook doed not work recursively for maps of maps.
		// DecodeHook: mapstructure.ComposeDecodeHookFunc(
		// ),
		WeaklyTypedInput: false,
		TagName:          "json",
		Result:           dst,
	}

	dec, err := mapstructure.NewDecoder(msCfg)
	if err != nil {
		return err
	}

	err = dec.Decode(&src)
	if err != nil {
		return err
	}
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

func processMap(src map[string]interface{}) map[string]interface{} {

	for k, v := range src {
		log.Printf("\t\tkey: %s, value: %v, (type): %T\n", k, v, v)

		//switch vv := v.(type) {
		switch vv := v.(type) {
		case map[string]interface{}:
			s := processMap(vv)
			src[k] = s
		default:
			ms := stringToRuneSlice(k, v)
			src[k] = ms
		}
	}
	return src
}

// stringToRuneSlice takes value and converts it accordingly for correct unmarshalling.
// Input is field name (map key) and the corresponding value.
func stringToRuneSlice(f string, i interface{}) interface{} {
	log.Printf(">>> field = %s, value = %+v, (type) = %T\n", f, i, i)

	_, ok := customFieldsM[f]

	switch v := i.(type) {
	case string:
		if !ok {
			return i
		}
		return []rune(v)
	default:
		return i
	}
}
