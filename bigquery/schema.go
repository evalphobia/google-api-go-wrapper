package bigquery

import (
	"errors"
	"reflect"
	"strings"
	"time"

	SDK "google.golang.org/api/bigquery/v2"
)

var (
	// TagName is used for struct tag which defining table column name
	TagName = "bigquery"

	errNotStructType     = errors.New("not a struct")
	errNotMapType        = errors.New("not a map")
	errInvalidStructType = errors.New("invalid struct for schema")
	errInvalidType       = errors.New("invalid type field for schema")
)

func convertToSchema(schemaStruct interface{}) (*SDK.TableSchema, error) {
	vv := reflect.ValueOf(schemaStruct)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	switch vt.Kind() {
	case reflect.Struct:
		// valid type
	case reflect.Map:
		return convertToSchemaFromMap(schemaStruct)
	default:
		return nil, errNotStructType
	}

	schema := &SDK.TableSchema{}
	for i, max := 0, vt.NumField(); i < max; i++ {
		f := vt.Field(i)
		if f.PkgPath != "" {
			continue // skip private field
		}

		tag, opts := parseTag(f, TagName)
		if tag == "-" {
			continue // skip `-` tag
		}

		v := vv.Field(i)
		if opts.has("squash") {
			convertToSchema(v.Interface())
			continue
		}

		name := getNameFromTag(f, TagName)
		fs, err := createFieldSchema(v)
		switch {
		case err != nil:
			return nil, err
		case opts.has("nullable"):
			fs.Mode = "nullable"
		}

		fs.Name = name
		schema.Fields = append(schema.Fields, fs)
	}

	return schema, nil
}

func convertToSchemaFromMap(schemaStruct interface{}) (*SDK.TableSchema, error) {
	vv := reflect.ValueOf(schemaStruct)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	if vt.Kind() != reflect.Map {
		return nil, errNotMapType
	}

	schema := &SDK.TableSchema{}
	for _, key := range vv.MapKeys() {
		v := vv.MapIndex(key)

		name, ok := key.Interface().(string)
		if !ok {
			continue
		}
		fs, err := createFieldSchema(reflect.ValueOf(v.Interface()))
		if err != nil {
			return nil, err
		}

		fs.Name = name
		schema.Fields = append(schema.Fields, fs)
	}

	return schema, nil
}

func createFieldSchema(v reflect.Value) (*SDK.TableFieldSchema, error) {
	fs := &SDK.TableFieldSchema{
		Mode: "required",
	}
	switch v.Kind() {
	case reflect.Bool:
		fs.Type = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fs.Type = "integer"
	case reflect.Float32, reflect.Float64:
		fs.Type = "float"
	case reflect.String:
		fs.Type = "string"
	case reflect.Array, reflect.Slice:
		fs.Mode = "repeated"
	case reflect.Struct:
		switch {
		case v.Type().ConvertibleTo(reflect.TypeOf(time.Time{})):
			fs.Type = "timestamp"
		default:
			return nil, errInvalidStructType
		}
	default:
		return nil, errInvalidType
	}

	return fs, nil
}

// getNameFromTag return the value in tag or field name in the struct field
func getNameFromTag(f reflect.StructField, tagName string) string {
	tag, _ := parseTag(f, tagName)
	if tag != "" {
		return tag
	}
	return f.Name
}

// getTagValues returns tag value of the struct field
func getTagValues(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}

// parseTag returns the first tag value of the struct field
func parseTag(f reflect.StructField, tag string) (string, tagOpation) {
	return splitTags(getTagValues(f, tag))
}

// splitTags returns the first tag value and rest slice
func splitTags(tags string) (string, tagOpation) {
	res := strings.Split(tags, ",")
	return res[0], res[1:]
}

// tagOpation is wrapper struct for rest tag values
type tagOpation []string

// has checks the value exists in the rest values or not
func (t tagOpation) has(tag string) bool {
	for _, opt := range t {
		if opt == tag {
			return true
		}
	}
	return false
}

func convertStructToMap(schemaStruct interface{}) (map[string]interface{}, error) {
	vv := reflect.ValueOf(schemaStruct)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	if vt.Kind() != reflect.Struct {
		return nil, errNotStructType
	}

	data := make(map[string]interface{})
	for i, max := 0, vt.NumField(); i < max; i++ {
		f := vt.Field(i)
		if f.PkgPath != "" {
			continue // skip private field
		}

		tag, opts := parseTag(f, TagName)
		if tag == "-" {
			continue // skip `-` tag
		}

		v := vv.Field(i)
		if opts.has("squash") {
			list, err := convertStructToMap(v.Interface())
			if err != nil {
				return data, err
			}

			for k, v := range list {
				data[k] = v
			}
			continue
		}

		val := v.Interface()
		if opts.has("nullable") && isZero(val) {
			continue
		}

		data[getNameFromTag(f, TagName)] = val
	}

	return data, nil
}

func isZero(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return v == false
	case string:
		return v == ""
	case zeroable:
		return v.IsZero()
	}
	return false
}

type zeroable interface {
	IsZero() bool
}

func isStruct(value interface{}) bool {
	vv := reflect.ValueOf(value)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	return vt.Kind() == reflect.Struct
}

func getSliceData(data interface{}) ([]interface{}, bool) {
	vv := reflect.ValueOf(data)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	if vt.Kind() != reflect.Slice {
		return nil, false
	}

	size := vv.Len()
	list := make([]interface{}, size)
	for i := 0; i < size; i++ {
		list[i] = vv.Index(i).Interface()
	}
	return list, true
}
