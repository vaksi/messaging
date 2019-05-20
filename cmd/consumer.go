package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vaksi/messaging/configs"
	"github.com/vaksi/messaging/internal/repositories"
	"github.com/vaksi/messaging/internal/services"
	"github.com/vaksi/messaging/pkg/kafka"
	"github.com/vaksi/messaging/pkg/logrus"
	"github.com/vaksi/messaging/pkg/mysql"
)

var (
	messagingConsumer     = "messaging-consumer"
	messagingConsumerInfo = `
Name     : %s
Topic    : %s
Group ID : %s
------------------------------------------------------------------------------
`
)

// MessagingConsumer consumer for bulk upload
type msgConsumer struct {
	stop <-chan bool

	BaseCmd       *cobra.Command
	filename string
	config *configs.Config
	MessagingService services.IMessageService
}

func NewConsumerCmd(
	configuration *configs.Config,
) *msgConsumer {
	return NewConsumerCmdSignaled(configuration, nil)
}

func NewConsumerCmdSignaled(
	configuration *configs.Config,
	stop <-chan bool,
) *msgConsumer {
	cc := &msgConsumer{stop: stop}
	cc.config = configuration
	cc.BaseCmd = &cobra.Command{
		Use:   "consumer",
		Short: "Used to run the http service",
		Run:  cc.RunConsumerMessaging,
	}
	fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
	fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
	cc.BaseCmd.Flags().AddFlagSet(fs)
	return cc
}


// NewMessagingConsumer creates new MessagingConsumer
func NewMessagingConsumer() *msgConsumer {
	return &msgConsumer{}
}

func (mc *msgConsumer) initService() {
	if len(mc.filename) > 1 {
		mc.config = configs.New(mc.filename,
			"./configs",
			"../configs",
			"../../configs",
			"../../../configs")
	}

	// init sql
	conn := fmt.Sprintf(
		mysql.MysqlDataSourceFormat,
		mc.config.MariaDB.User,
		mc.config.MariaDB.Password,
		mc.config.MariaDB.Host,
		mc.config.MariaDB.Port,
		mc.config.MariaDB.DbName,
		mc.config.MariaDB.Charset,
	)

	//db.OpenConnection(conn, injector.Config)
	db := mysql.NewMySQL()
	db.OpenConnection(conn, mc.config)
	db.SetConnMaxLifetime(mc.config.MariaDB.MaxLifeTime)
	db.SetMaxIdleConn(mc.config.MariaDB.MaxIdleConnection)
	db.SetMaxOpenConn(mc.config.MariaDB.MaxOpenConnection)


	// init kafka
	// Set kafka
	kafkaAddrs := strings.Split(mc.config.Kafka.BrokerList, ",")
	k := kafka.NewKafka(kafkaAddrs) // Set kafka

	// init service
	messageRepo := repositories.NewMessageRepository(db, *k, mc.config)

	messageSvc := services.NewMessageService(messageRepo)

	mc.MessagingService = messageSvc
}

func (mc *msgConsumer) kafkaMessageHandler(key []byte, msg []byte) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Error(err)
		}
	}()

	err := mc.MessagingService.ReceiveMessage(key, msg)
	if err != nil {
		logrus.Error(err)
	}
}

func (buc *msgConsumer) kafkaErrorHandler(err error) {
	logrus.Error(err)
}

func (buc *msgConsumer) kafkaNotificationHandler(notification interface{}) {
	logrus.Debug(notification)
}

func (buc *msgConsumer) printInfo(topics []string, groupID string) {
	fmt.Printf(messagingConsumerInfo, messagingConsumer, topics, groupID)
}

// Run Consumer Message
func (mc *msgConsumer) RunConsumerMessaging(cmd *cobra.Command, args []string)  {

	mc.initService()

	var (
		topics  = []string{mc.config.Kafka.MessagingConsumer.Topic}
		groupID = mc.config.Kafka.MessagingConsumer.Group
		broker  = strings.Split(mc.config.Kafka.BrokerList, ",")
	)

	mc.printInfo(topics, groupID)
	done := make(chan bool)
	k := kafka.NewKafka(broker)

	err := k.Consume(groupID, topics, mc.kafkaMessageHandler, mc.kafkaErrorHandler, mc.kafkaNotificationHandler, done)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("shutting down")
	os.Exit(0)
}
