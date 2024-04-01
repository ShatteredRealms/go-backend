package helpers

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func BindEnvsToStruct(obj interface{}) {
	viper.AutomaticEnv()

	val := reflect.ValueOf(obj)
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		key := field.Name
		if field.Anonymous {
			key = ""
		}
		bindRecursive(key, val.Field(i))
	}
}

func bindRecursive(key string, val reflect.Value) {
	if val.Kind() != reflect.Struct {
		env := "SRO_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		viper.MustBindEnv(key, env)
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		newKey := field.Name
		if field.Anonymous {
			newKey = ""
		} else if key != "" {
			newKey = "." + newKey
		}

		bindRecursive(key+newKey, val.Field(i))
	}
}
