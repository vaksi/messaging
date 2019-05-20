package configs

import (
	"github.com/spf13/viper"
	"github.com/vaksi/messaging/pkg/logrus"
)

type Constants struct {
	// App as application config
	App struct {
		Name            string `mapstructure:"name"`
		Port            int    `mapstructure:"port"`
		ReadTimeout     int    `mapstructure:"read_timeout"`
		WriteTimeout    int    `mapstructure:"write_timeout"`
		Timezone        string `mapstructure:"timezone"`
		Debug           bool   `mapstructure:"debug"`
		Env             string `mapstructure:"env"`
		SecretKey       string `mapstructure:"secret_key"`
		CallbackTimeout int    `mapstructure:"callback_timeout"`
	}
	// Api as Endpoint Config
	Api struct {
		Prefix string `mapstructure:"prefix"`
	}
	// Log as logging config
	Log struct {
		Dir      string `mapstructure:"dir"`
		Filename string `mapstructure:"filename"`
	}
	// CB as Circuit Breakers Config
	CB struct {
		Retry      int `mapstructure:"retry_count"`
		Timeout    int `mapstructure:"db_timeout"`
		Concurrent int `mapstructure:"max_concurrent"`
	}
	// MariaDB as Sql MariaDB config
	MariaDB struct {
		DbName            string `mapstructure:"db_name"`
		Host              string `mapstructure:"host"`
		Port              int    `mapstructure:"port"`
		User              string `mapstructure:"user"`
		Password          string `mapstructure:"password"`
		MaxLifeTime       int    `mapstructure:"max_life_time"`
		MaxIdleConnection int    `mapstructure:"max_idle_connection"`
		MaxOpenConnection int    `mapstructure:"max_open_connection"`
		Charset           string `mapstructure:"charset"`
	}

	// Kafka config
	Kafka struct{
		BrokerList string `mapstructure:"brokerList"`
		GroupID string `mapstructure:"groupId"`
		MessageCountStart string `mapstructure:"messageCountStart"`
		OffsetType string `mapstructure:"offsetType"`
		Partition  string `mapstructure:"partition"`
		MessagingConsumer struct {
			Topic string `mapstructure:"topic"`
			Group string `mapstructure:"group"`
			TimeSleep int `mapstructure:"time_sleep"`
			MessagePoll int `mapstructure:"message_poll"`
			ConsumerType string `mapstructure:"consumer_type"`
		} `mapstructure:"messaging_consumer"`
	}
}

type Config struct {
	Constants
}

// New is used to generate a configuration instance which will be passed around the codebase
func New(filename string, paths ...string) *Config {
	if filename == "" {
		filename = "app"
	}
	constants := initViper(filename, paths...)
	// init logger
	logrus.Init(constants.Log.Filename,
		constants.App.Debug,
		constants.Log.Dir,
		"./"+constants.Log.Dir,
		"../"+constants.Log.Dir,
		"../../"+constants.Log.Dir)
	// init Infrastructure DB
	return &Config{
		constants,
	}
}

func initViper(filename string, paths ...string) (constants Constants) {
	vip := viper.New()
	// Search the root directory for the configuration file
	for _, path := range paths {
		vip.AddConfigPath(path)
	}
	// Configuration fileName without the .TOML or .YAML extension
	vip.SetConfigName(filename)
	if err := vip.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	vip.WatchConfig() // Watch for changes to the configuration file and recompile
	if err := vip.UnmarshalExact(&constants); err != nil {
		panic(err.Error())
	}
	return constants
}