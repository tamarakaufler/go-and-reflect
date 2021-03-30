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