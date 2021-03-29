package env

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

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

		log.Printf("\t\t<<< Type of input Value - t: [%+v]\n", t)
		log.Printf("\t\t<<< Struct Field - f: [%+v]\n", f)
		log.Printf("\t\t<<< Struct Type Field - tf: [%+v]\n", tf)

		// struct field is a non-nil pointer.
		if f.Kind() == reflect.Ptr && !f.IsNil() {
			err := Parse(f.Interface()) // Parse accepts an interface type, so we get the f value as an interface
			if err != nil {
				return err
			}
			continue
		}

		// struct field itself is a struct.
		if f.Kind() == reflect.Struct && f.CanAddr() { // Addr refers to memory address.
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
		p := strings.SplitN(e, "=", 2)
		envM[p[0]] = p[1]
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
	log.Printf("\t\tTag Info: [%+v]\n", ti)

	fieldVal := ti.defVal
	if ti.envName != "" {
		envVal, ok := envVars[ti.envName]
		if ok {
			log.Printf("\t\tEnv Info: [%s]\n", envVal)
			return envVal, nil
		}
		if ti.required {
			return "", fmt.Errorf("%s requires environment variable %s to be set", sf.Name, ti.envName)
		}
	}
	return fieldVal, nil
}

// setField sets a struct field's value. It accepts the field Value, the field Type and the value
// to set the field to.
func setValue(f reflect.Value, t reflect.StructField, val string) error {
	tt := t.Type
	ff := f
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem() // retrieving element type
		ff = f.Elem()  // returns the value the pointer points to
	}

	parseF, ok := defaultParsers[tt.Kind()]
	if !ok {
		return fmt.Errorf("no parser found for %s", t.Name)
	}

	vv, err := parseF(val)
	if err != nil {
		return fmt.Errorf("failed to parse value %s for field %s", val, t.Name)
	}

	ff.Set(reflect.ValueOf(vv).Convert(tt)) // converts reflect.ValueOf(vv) into type corresponding to ff
	log.Printf("\t\t>>> SET var: ff - %+v\n\n", ff)

	return nil
}
