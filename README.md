# go-and-reflect (on data structure and life)

A little exercise I undertook to get a better understanding of Go reflection.

# Go reflection

## Concept

Go reflection, introspection during runtime of an applicationrevolves around three concepts:

- reflect.Type ........... provides information about the data structure's name (User, Address etc), its field names (Name, Age), field tags metadata
- reflect.Kind ........... tells if the the data structure is a struct, pointer, map, string etc
- reflect.Value .......... hold information about the values related to the inspected data structure

   reflect.Type
        |
        |  reflect.Kind
        |      |
type Person struct {
    Name string ....... field name Name                                     examples of reflection
    Age float32 ....... field name Age ........................ reflect.Type:  {Name:Age PkgPath: Type:float32
                                                                                Tag:env:"USER_AGE" defVal:"23.5"
                                                                                Offset:16 Index:[1] Anonymous:false}
                                                                reflect.Value:  33
    Address ........... an embedded field, has no field name
}

type Address struct { ......................................... reflect.Type:  {Name:Address PkgPath: Type:main.Address Tag:
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

## Details

### reflect package - investigating data structures

First step towards the env package.

### env package - poor man's carloos0/env