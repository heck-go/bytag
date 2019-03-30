package bytag

import (
	"reflect"
	"strings"
)

// Bind binds the content of data into the struct s
func Bind(tagName string, s interface{}, data map[string]interface{}) interface{} {
	if s == nil {
		return nil
	}
	
	t := reflect.TypeOf(s)
	tk := t.Kind()
	
	if tk != reflect.Ptr {
		return nil
	}
	
	t = t.Elem()
	tk = t.Kind()
	
	if tk != reflect.Struct {
		return nil
	}
	
	v := reflect.ValueOf(s).Elem()
	
	for i := 0; i < t.NumField(); i ++ {
		f := t.Field(i)
		fi := ParseField(tagName, f)
		
		fv := v.FieldByName(fi.Name)
		fk := fv.Kind()
		ft := fv.Type()
		
		if v, ok := data[fi.Alias]; ok {
			vt := reflect.TypeOf(v)
			vk := vt.Kind()
			
			if vt.AssignableTo(ft) {
				fv.Set(reflect.ValueOf(v))
				continue
			}
			
			if fk == reflect.Struct && vk == reflect.Map {
				Bind(tagName, fv.Addr().Interface(), v.(map[string]interface{}))
				continue
			}
		}
	}
	
	return s
}

func BindSlice(tagName string, s reflect.Value, data []interface{}) {
	// sk := s.Kind()
	et := s.Type().Elem()
	ek := et.Kind()
	
	ret := reflect.MakeSlice(et, s.Len(), s.Cap())
	vet := reflect.TypeOf(data).Elem()
	if vet.AssignableTo(et) {
		for i := 0; i < s.Len(); i++ {
			ret.Index(i).Set(reflect.ValueOf(data[i]))
		}
	} else if ek == reflect.Struct {
		for i := 0; i < s.Len(); i++ {
			// v := Bind(tagName, ret.Index(i).Addr().Interface(), data[i].(map[string]interface{}))
			v := Bind(tagName, ret.Index(i).Addr().Interface(), data[i].(map[string]interface{}))
			ret.Index(i).Set(reflect.ValueOf(v))
		}
	}
}

type FieldInfo struct {
	Alias string
	
	Name string
}

// ParseField parses [FieldInfo] for the given struct field [f] from struct tag with name [tagName]
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