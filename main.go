package bytag

import (
	"reflect"
	"strings"
)

// Bind binds the content of data into the struct s
func Bind(tagName string, s interface{}, data map[string]interface{}) {
	if s == nil {
		return
	}
	
	t := reflect.TypeOf(s)
	
	if t.Kind() != reflect.Ptr {
		return
	}
	
	t = t.Elem()
	
	if t.Kind() != reflect.Struct {
		return
	}
	
	v := reflect.ValueOf(s).Elem()
	
	for i := 0; i < t.NumField(); i ++ {
		f := t.Field(i)
		fi := ParseField(tagName, f)
		
		fv := v.FieldByName(fi.Name)
		ft := fv.Type()
		
		if v, ok := data[fi.Alias]; ok {
			vt := reflect.TypeOf(v)
			
			if vt.AssignableTo(ft) {
				fv.Set(reflect.ValueOf(v))
			} else {
			
			}
		}
	}
}

type FieldInfo struct {
	Alias string
	
	Name string
}

func ParseField(tagName string, f reflect.StructField) *FieldInfo {
	var parts []string
	alias := f.Name
	
	tag, tagOk := f.Tag.Lookup(tagName)
	if tagOk {
		partsTemp := strings.Split(tag, ",")
		parts = make([]string, 0, len(partsTemp))
		for i := 0; i < len(partsTemp); i++ {
			part := strings.TrimSpace(partsTemp[i])
			if len(part) != 0 {
				parts = append(parts, part)
			}
		}
	}
	
	if len(parts) != 0 {
		alias = parts[0]
		// TODO parse other tags
	}
	
	return &FieldInfo{
		Alias:alias,
		Name: f.Name,
	}
}