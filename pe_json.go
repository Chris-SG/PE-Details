package pe_details

import (
	"fmt"
	"reflect"
)

func (o PEOffset) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, fmt.Sprintf("\"%#x\"", o)...)
	return b, nil
}

func (ch COFFHeader) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(ch)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (sf OHStandardFields) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(sf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (wf OHWindowsFields32) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (wf OHWindowsFields64) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (idd ImageDataDirectory) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(idd)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (st SectionTable) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(st)
	b = append(b, fmt.Sprintf("\"%s\":\"%+#v\",", v.Type().Field(0).Name, v.Field(0).Interface())...)
	for i := 1; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (cr CoffRelocation) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(cr)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (cs CoffSymbol) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(cs)
	for i := 0; i < v.NumField()-1; i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (as AuxSymbol) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(as)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}