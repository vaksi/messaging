package configs

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var config = []byte(`
app:
  name: massaging
  port: 8081
  env: dev #dev|stg|prod
  timezone: Asia/Jakarta
  debug: true  #in production set false
  read_timeout: 10 # in second
  write_timeout: 10 # in second
  callback_timeout: 10 # in second

api:
  prefix: "/massaging/v1"

MariaDB:
  db_name: massaging
  host: localhost
  port: 3306
  user: root
  charset: utf8

logrus:
  dir: "logs"
  filename: "massaging.logrus"
`)

func initConfig(t *testing.T, c []byte) {
	viper.Reset()
	viper.New()
	viper.SetConfigType("yaml")
	r := bytes.NewReader(c)

	var (
		err error
		n   int64
	)
	buf := new(bytes.Buffer)
	n, err = buf.ReadFrom(r)

	assert.IsType(t, int64(345), n)
	assert.Nil(t, err)
	err = viper.ReadConfig(buf)
	assert.Nil(t, err)
}

func TestNewSuccess(t *testing.T) {
	initConfig(t, config)
	var constants Constants
	err := viper.Unmarshal(&constants)
	assert.Nil(t, err)

	c := Config{}
	c.Constants = constants
	configuration := New("app.test",
		"./configs", "../configs", "../../configs")
	assert.Equal(t, &c, configuration)

}

func TestNewFailYAML(t *testing.T) {
	f := func() {
		New("app.forfailtest",
			"./configs", "../configs", "../../configs")
	}
	assert.Panics(t, f)
	assert.PanicsWithValue(t, "1 error(s) decoding:\n\n* '' has invalid keys: appfalse", f)
}

func TestNewConfigNotFound(t *testing.T) {
	f := func() {
		New("app.test", "")
	}
	assert.Panics(t, f)
	assert.PanicsWithValue(t, "Config File \"app.test\" Not Found in \"[]\"", f)
}

func TestInitViperConfigNotFound(t *testing.T) {
	f := func() {
		initViper("app.test")
	}
	assert.Panics(t, f)
	assert.PanicsWithValue(t, "Config File \"app.test\" Not Found in \"[]\"", f)
}
