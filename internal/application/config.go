package application

import (
	"os"
	"reflect"
)

type Config struct {
	POSTGRES_DSN string

	// SMTP
	SMTP_USERNAME        string
	SMTP_PASSWORD        string
	SMTP_HOST            string
	SMTP_SENDER_USERNAME string
}

func NewConfig() (c Config) {
	valueOf := reflect.Indirect(reflect.ValueOf(&c))
	ofType := valueOf.Type()
	for i := 0; i < valueOf.NumField(); i++ {
		typeField := ofType.Field(i)
		valueField := valueOf.Field(i)
		value, ok := os.LookupEnv(typeField.Name)
		if !ok {
			panic("Cant't find env variable \"" + typeField.Name + "\"")
		}
		valueField.SetString(value)
	}
	return
}
