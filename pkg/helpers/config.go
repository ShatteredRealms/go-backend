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
		env := field.Name
		if field.Anonymous {
			env = ""
		}
		bindRecursive(key, env, val.Field(i))
	}
}

func bindRecursive(key, env string, val reflect.Value) {
	if val.Kind() != reflect.Struct {
		env = "SRO_" + strings.ToUpper(env)
		_ = viper.BindEnv(key, env)
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		newKey := field.Name
		newEnv := field.Name
		if key != "" {
			newKey = "." + newKey
		}
		if field.Anonymous {
			newEnv = ""
		} else if env != "" {
			newEnv = "_" + newEnv
		}

		bindRecursive(key+newKey, env+newEnv, val.Field(i))
	}
}
