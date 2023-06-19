package configuration

import (
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"os"
	"reflect"
)


type Configuration struct {
	Prefix string `env:"PREFIX" default:"/"`
	Server string `env:"SERVER" default:""`
	Port   string `env:"PORT" default:"80"`
}

func (c *Configuration) Host() string {
	log.Dev("Host: %s:%s\n", c.Server, c.Port)
	return c.Server + ":" + c.Port
}

var DefaultConfiguration = &defaultConfiguration

var defaultConfiguration = Configuration{}

func (c *Configuration) ApplyDefaultsOnEmpty() *Configuration {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	c.applyDefault(v, t, "Prefix")
	c.applyDefault(v, t, "Server")
	c.applyDefault(v, t, "Port")

	return c
}

func (c *Configuration) applyDefault(v reflect.Value, t reflect.Type, f string) {
	var newVal string

	value := v.FieldByName(f)
	if !value.IsZero() || value.Kind() != reflect.String {
		return
	}

	field, ok := t.FieldByName(f)
	if !ok {
		return
	}

	newVal = os.Getenv(field.Tag.Get("env"))
	if newVal == "" {
		newVal = field.Tag.Get("default")
	}
	value.SetString(newVal)
}
