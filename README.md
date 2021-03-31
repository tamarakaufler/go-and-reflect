# go-and-reflect (on data structures and life)

Learning more about Go reflection.

# Go reflection

## Concept

Go reflection/introspection during runtime of an application revolves around three concepts:

- reflect.Type ........... provides information about the data structure's name (User, Address etc), its field names
                           (Name, Age), field tags metadata.
- reflect.Kind ........... tells if a data structure is a struct, pointer, map, string etc
- reflect.Value .......... holds information about the values inspected data structure field

**reflect.Type** provides information about data structure and its fields for identification and metadata
                 (User, Address, Age, field tags etc).

**reflect.Kind** identifies whether the data structure is a struct, pointer, string, float32 etc.

**reflect.Value** informs about the values fields hold and, if allowed, provides a setter.

```
   reflect.Type
        |
        | reflect.Kind
        |    |
type User struct {
    Name string ....... field name Name                            examples of reflection
    Age float32 ....... field name Age ............ reflect.Type:  {Name:Age PkgPath: Type:float32
                                                                   Tag:env:"USER_AGE" envDefault:"23.5"
                                                                   Offset:16 Index:[1] Anonymous:false}
                                                    reflect.Value: 33
    Address ........... an embedded field, has no field name
}

type Address struct { ............................. reflect.Type:  {Name:Address PkgPath: Type:main.Address Tag:
                                                                   Offset:24 Index:[2] Anonymous:false}
                                                    reflect.Value: {Street: City: Postcode: LatLng:{Lat:0 Lng:0}}
    Street string
    City string
    Postcode string
    LatLng
}

type LatLng struct {
    Latitude, Longitude float64
}
```

## Details

### reflect package - investigating data structures

First step towards the env package.

### env package - poor man's carloos0/env

### cmd/json/main.go - custom marshalling/unmarshalling.

The concrete types User2 and User3 have identical fields:

```
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
```

User2 has a custom marshaller and unmarshaller. Note, Street and City fields are rune slices, however they are required to be marsaled into strings and unmarshaled from strings to rune slices. Other fields, that are either rune or int32 slices are required to be marshalled and unmashalled as such.

Go treats []rune and []int32 the same because rune is an alias for int32. The custom marshaller marshals rune slices into string while leaving the NI field as is, ie it is marshaled as []int32.

#### marshalling

Custom marshaling for User2 implements (User2).MarshalJSON. The implementation uses
pure reflection.

#### unmarshalling

Custom unmarshaling for User2 implements (*User2).UnmarshalJSON. The implementation
uses type assertion and mapstructure library.
