package req

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	schemaIndent = "\t"
)

// GetSchema writes to `dst` a JSON-like schema of type `typ`. It returns nil
// iff it finished writing the schema. It returns an error instead.
func GetSchema(dst io.Writer, typ reflect.Type) error {
	if _, err := dst.Write([]byte("{\n")); err != nil {
		return err
	}
	if err := getSchemaBody(dst, typ); err != nil {
		return err
	}
	if _, err := dst.Write([]byte("}")); err != nil {
		return err
	}
	return nil
}

func getSchemaBody(dst io.Writer, typ reflect.Type) error {
	for i, fc := 0, typ.NumField(); i < fc; i++ {
		fld := typ.FieldByIndex([]int{i})
		if err := getSchemaField(dst, fld, i+1 == fc); err != nil {
			return err
		}
	}
	return nil
}

func getSchemaField(dst io.Writer, fld reflect.StructField, last bool) error {
	if err := getSchemaFieldKey(dst, fld); err != nil {
		return err
	}
	if err := getSchemaFieldValue(dst, fld, last); err != nil {
		return err
	}
	return nil
}

func getSchemaFieldKey(dst io.Writer, fld reflect.StructField) error {
	jsonName := getJsonName(fld)
	key := fmt.Sprintf(schemaIndent+`"%s": `, jsonName)
	_, err := dst.Write([]byte(key))
	return err
}

func getSchemaFieldValue(dst io.Writer, fld reflect.StructField, last bool) error {
	var decorated string
	var sep string
	if isFieldNonZero(fld) {
		decorated = "¡%s!"
	} else {
		decorated = "¿%s?"
	}
	if last {
		sep = ""
	} else {
		sep = ","
	}
	decorated = fmt.Sprintf(decorated, fld.Type.String())
	value := fmt.Sprintf("%s%s\n", decorated, sep)
	_, err := dst.Write([]byte(value))
	return err
}

func isFieldNonZero(fld reflect.StructField) bool {
	tag := fld.Tag.Get(StructTagName)
	tags := strings.Split(tag, ",")
	for _, t := range tags {
		if t == StructTagNonZero {
			return true
		}
	}
	return false
}
