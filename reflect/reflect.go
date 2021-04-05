package reflect

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type (
	tagInfo struct {
		envName    string
		required   bool
		envDefault string
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

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")
	log.Printf("\t*** t.Name(): [%+v]\n", t.Name())
	log.Printf("\t*** v: [%+v] v.Kind() [%+v]\n", v, v.Kind())

	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		tf := t.Field(i) // eg tf.Type.Name() == LatLng

		fmt.Println("--------------------------------------------")
		log.Printf("\t\t*** Field ... f.Type() == f.Type().Name(): [%+v], tf.Type == tf.Type.Name(): [%+v]\n",
			f.Type(), tf.Type)

		log.Printf("\t\t\t<<< Struct Type Field tf: [%+v] (tf.Type [%+v]) - tf.Type.Kind() [%+v]\n",
			tf, tf.Type, tf.Type.Kind())
		log.Printf("\t\t\t>>> Struct Field f: [%+v] (f.Type().Kind() [%+v]) - f.Kind() [%+v]\n",
			f, f.Type().Kind(), f.Kind())

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

		ti := processTag(tf)
		log.Printf("\t\t\t--> Tag Info: [%+v]\n", ti)
		fmt.Println("--------------------------------------------")
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")

	return nil
}

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

	d, okD := sf.Tag.Lookup("envDefault")
	if okD {
		ti.envDefault = d
	}
	return ti
}
