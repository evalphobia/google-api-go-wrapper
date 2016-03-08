package bigquery

import (
	"errors"
	"reflect"
	"strings"
	"time"

	bigquery "google.golang.org/api/bigquery/v2"
)

var (
	// TagName is used for struct tag which defining table column name
	TagName = "bigquery"

	errNotStructType     = errors.New("not a struct")
	errInvalidStructType = errors.New("invalid struct for schema")
	errInvalidType       = errors.New("invalid type field for schema")
)

func convertToSchema(schemaStruct interface{}) (*bigquery.TableSchema, error) {
	vv := reflect.ValueOf(schemaStruct)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	vt := vv.Type()
	if vt.Kind() != reflect.Struct {
		return nil, errNotStructType
	}

	schema := &bigquery.TableSchema{}
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

func createFieldSchema(v reflect.Value) (*bigquery.TableFieldSchema, error) {
	fs := &bigquery.TableFieldSchema{
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
