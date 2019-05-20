package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vaksi/messaging/internal/http/handler"
	"github.com/vaksi/messaging/internal/repositories"
	"github.com/vaksi/messaging/internal/services"
	"github.com/vaksi/messaging/pkg/kafka"

	"github.com/vaksi/messaging/configs"
	internalHttp "github.com/vaksi/messaging/internal/http"
	"github.com/vaksi/messaging/pkg/mysql"
)

type httpCmd struct {
	stop <-chan bool

	BaseCmd       *cobra.Command
	configuration *configs.Config
	filename      string
}

func NewHttpCmd(
	configuration *configs.Config,
) *httpCmd {
	return NewHttpCmdSignaled(configuration, nil)
}

func NewHttpCmdSignaled(
	configuration *configs.Config,
	stop <-chan bool,
) *httpCmd {
	cc := &httpCmd{stop: stop}
	cc.configuration = configuration
	cc.BaseCmd = &cobra.Command{
		Use:   "http",
		Short: "Used to run the http service",
		RunE:  cc.server,
	}
	fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
	fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
	cc.BaseCmd.Flags().AddFlagSet(fs)
	return cc
}

func (h *httpCmd) server(cmd *cobra.Command, args []string) (err error) {
	if len(h.filename) > 1 {
		h.configuration = configs.New(h.filename,
			"./configs",
			"../configs",
			"../../configs",
			"../../../configs")
	}

	// init sql
	conn := fmt.Sprintf(
		mysql.MysqlDataSourceFormat,
		h.configuration.MariaDB.User,
		h.configuration.MariaDB.Password,
		h.configuration.MariaDB.Host,
		h.configuration.MariaDB.Port,
		h.configuration.MariaDB.DbName,
		h.configuration.MariaDB.Charset,
	)

	//db.OpenConnection(conn, injector.Config)
	db := mysql.NewMySQL()
	db.OpenConnection(conn, h.configuration)
	db.SetConnMaxLifetime(h.configuration.MariaDB.MaxLifeTime)
	db.SetMaxIdleConn(h.configuration.MariaDB.MaxIdleConnection)
	db.SetMaxOpenConn(h.configuration.MariaDB.MaxOpenConnection)

	// init kafka
	// Set kafka
	kafkaAddrs := strings.Split(h.configuration.Kafka.BrokerList, ",")
	k := kafka.NewKafka(kafkaAddrs) // Set kafka

	// init service
	messageRepo := repositories.NewMessageRepository(db, *k, h.configuration)

	messageSvc := services.NewMessageService(messageRepo)

	messageHandler := handler.NewMessageHandler(messageSvc)

	// inject routes
	route := &internalHttp.Routes{
		Config:              h.configuration,
		MsgHandler: messageHandler,
	}
	router := route.NewRoutes()

	// Description µ micro service
	fmt.Println(
		fmt.Sprintf(
			WelkomText,
			h.configuration.App.Port,
			strings.Join([]string{
				h.configuration.Log.Dir,
				h.configuration.Log.Filename}, "/"),
		))

	//tableRoute(router) // Prettier Route Pattern

	return h.serve(router)
}

func (h *httpCmd) serve(router http.Handler) error {
	errCh := make(chan error, 1)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	addr := net.JoinHostPort("",
		strconv.Itoa(h.configuration.App.Port))
	s := StartWebServer(
		addr,
		h.configuration.App.ReadTimeout,
		h.configuration.App.WriteTimeout,
		router,
	)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logrus.Info(
				"Server gracefully ListenAndServe",
				logrus.Fields{})
			errCh <- err
		}
		<-h.stop
	}()

	if h.stop != nil {
		select {
		case err := <-errCh:
			logrus.Info(
				"Server gracefully h stop stopped",
				logrus.Fields{})
			return err
		case <-h.stop:
		case <-quit:
		}
	} else {
		select {
		case err := <-errCh:
			logrus.Info(
				"Server gracefully stopped",
				logrus.Fields{})
			return err
		case <-quit:
		}
	}
	return nil
}

// StartWebServer starts a web server
func StartWebServer(addr string, readTimeout, writeTimeout int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
}
