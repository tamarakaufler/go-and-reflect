package env

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//var envVars map[string]string

// func init() {
// 	//envVars = GetEnvVars()
// }

type (
	parseFunc func(string) (interface{}, error)
	tagInfo   struct {
		envName  string
		required bool
		defVal   string
	}
)

var (
	envVars        = GetEnvVars()
	defaultParsers = map[reflect.Kind]parseFunc{
		reflect.String: func(s string) (interface{}, error) {
			return s, nil
		},
		reflect.Bool: func(s string) (interface{}, error) {
			p, err := strconv.ParseBool(s)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Float32: func(s string) (interface{}, error) {
			p, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Float64: func(s string) (interface{}, error) {
			p, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Int: func(s string) (interface{}, error) {
			p, err := strconv.ParseInt(s, 0, 0)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Int8: func(s string) (interface{}, error) {
			p, err := strconv.ParseInt(s, 0, 8)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Int32: func(s string) (interface{}, error) {
			p, err := strconv.ParseInt(s, 0, 32)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
		reflect.Int64: func(s string) (interface{}, error) {
			p, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return nil, err
			}
			return p, nil
		},
	}
)

// Parse expects the provided data structure and reports on its content. The input must be
// a pointer to a struct.
func Parse(c interface{}) error {
	log.Println("-----------------------------------------------------")
	defer func() {
		log.Println("-----------------------------------------------------")
	}()

	// creates a new initialised concrete type stored in the provided interface c.
	v := reflect.ValueOf(c)

	// the provided concrete type must be a pointer.
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("input %+v must be a pointer", c)
	}

	// now we need the value that the interface v contains.
	e := v.Elem()
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("the dynamic type of the input %+v must be a struct", e)
	}

	envVars = GetEnvVars()

	err := parse(e)
	if err != nil {
		return err
	}

	return nil
}

// parse accepts a struct value.
func parse(v reflect.Value) error {
	log.Printf("PARSE INPUT: v [%+v]\n", v.Type().Name())
	t := v.Type() // type of the struct, eg Address (=> t.Name() = Address)

	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)

		if !f.CanSet() {
			continue
		}

		tf := t.Field(i) // eg tf.Type.Name() == LatLng
		log.Printf("\t*** Field ... f.Type().Name(): [%+v],\nt: [%+v], tf: [%+v]\n\n",
			f.Type().Name(), t, tf)

		log.Printf("\t\t>>> Struct Field f: [%+v]\n", f)
		log.Printf("\t\t<<< Struct Type Field tf: [%+v]\n", tf)

		// struct field is a non-nil pointer.
		if f.Kind() == reflect.Ptr && !f.IsNil() {
			log.Println("\t\t== pointer and not nil ...")

			log.Printf("\t\t== f type name: [%+v]\n", f.Type().Name())
			err := Parse(f.Interface()) // Parse accepts an interface type, so we get the f value as an interface
			if err != nil {
				return err
			}
			continue
		}

		// struct field itself is a struct.
		// Addr refers to memory address
		//if f.Kind() == reflect.Struct && f.CanAddr() && f.Type().Name() == "" {
		if f.Kind() == reflect.Struct && f.CanAddr() {
			log.Println("\t\t-- struct and addressable ...")

			log.Printf("\t\t-- f.Addr(): %+v ... f.Type().Name(): [%+v] ... t.Name(): [%+v]\n",
				f.Addr(), f.Type().Name(), t.Name())
			// f.Addr() returns a pointer to the struct f and Interface() returns the value of f as an interface.
			err := Parse(f.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}
		fieldV, err := getValue(tf, envVars)
		if err != nil {
			return err
		}
		log.Printf("\t\t@@@ obtained field value: [%s]\n", fieldV)

		err = setValue(f, tf, fieldV)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetEnvVars() map[string]string {
	envs := os.Environ()

	envM := make(map[string]string, len(envs))
	for _, e := range envs {
		//log.Printf("\t\t+++ 1 ENV Info: [%+v]\n", e)

		p := strings.SplitN(e, "=", 2)
		envM[p[0]] = p[1]
		//log.Printf("\t\t+++ 2 ENV Info: %s = %s\n", p[0], p[1])
	}
	return envM
}

// processTag accepts v.Type().Field(i), where v is reflect.Value, value that an interface contains.
// v.Type().Field(i) contains also struct field tags metadata.
func processTag(sf reflect.StructField) tagInfo {
	ti := tagInfo{}

	t, okE := sf.Tag.Lookup("env")
	if okE {
		p := strings.Split(t, ",")
		ti.envName = p[0]
		if len(p) == 2 && p[1] == "required" {
			ti.required = true
			return ti
		}
	}

	d, okD := sf.Tag.Lookup("defVal")
	if okD {
		ti.defVal = d
	}
	return ti
}

func getValue(sf reflect.StructField, envVars map[string]string) (string, error) {
	ti := processTag(sf)
	log.Printf("\t\t### Tag Info: [%+v]\n", ti)

	fieldVal := ti.defVal
	if ti.envName != "" {
		envVal, ok := envVars[ti.envName]
		//log.Printf("\t\t### Env var: ti.envName=%s ... envVal=%s\n", ti.envName, envVal)

		if ok {
			return envVal, nil
		}
		if ti.required {
			return "", fmt.Errorf("%s requires environment variable %s to be set", sf.Name, ti.envName)
		}
	}
	return fieldVal, nil
}

func setValue(f reflect.Value, sf reflect.StructField, val string) error {

	return nil
}
