package config

import (
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func BindEnvsToStruct(obj interface{}) {
	viper.AutomaticEnv()

	val := reflect.ValueOf(obj)
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		bindRecursive(val.Type().Field(i).Name, val.Field(i))
	}
}

func bindRecursive(name string, val reflect.Value) {
	if val.Kind() != reflect.Struct {
		viper.BindEnv(name, strings.ReplaceAll(name, ".", "_"))
		return
	}

	for i := 0; i < val.NumField(); i++ {
		bindRecursive(name+"."+val.Type().Field(i).Name, val.Field(i))
	}
}
