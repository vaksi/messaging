package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/vaksi/messaging/configs"
	"math/rand"
	"time"
)

// MySqlFactory is an abstract for sql database
type MySqlFactory interface {
	OpenConnection(connString string, config *configs.Config)
	Close() error
	GetDB() (*DB, error)
	QueryRow(query string, args ...interface{}) (*sql.Row, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	SetConnMaxLifetime(int)
	SetMaxIdleConn(int)
	SetMaxOpenConn(int)
}

const MYSQL = "mysql"

type fallbackFunc func(error) error

type DB struct {
	*sql.DB
	config       *configs.Config
	commandName  *string
	retryCount   int
	fallbackFunc func(error) error
}

func NewMySQL() MySqlFactory {
	return &DB{}
}

var MysqlDataSourceFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local"

// GetDB gets database connection
func (r *DB) GetDB() (*DB, error) {
	err := r.Ping()
	if err != nil {
		logrus.Error("Get DB", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	return r, nil
}

// OpenConnection gets a handle for a database
func (r *DB) OpenConnection(connString string, config *configs.Config) {
	db, err := open(MYSQL, connString, config)
	if err != nil {
		panic(err)
	}
	r.DB = db.DB
	_, err = r.GetDB()
	if err != nil {
		panic(err.Error())
	}
}

func open(driverName, connString string, config *configs.Config) (*DB, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		panic(err.Error())
	}
	return &DB{
		DB:     db,
		config: config,
	}, nil
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (r *DB) SetConnMaxLifetime(connMaxLifetime int) {
	r.DB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func (r *DB) SetMaxIdleConn(maxIdleConn int) {
	r.DB.SetMaxIdleConns(maxIdleConn)
}

// SetMaxOpenConns sets the maximum amount of time a connection may be reused.
func (r *DB) SetMaxOpenConn(maxOpenConn int) {
	r.DB.SetMaxOpenConns(maxOpenConn)
}

func (r *DB) Close() error {
	return r.DB.Close()
}

func (r *DB) Begin() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return r.DB.Exec(query, args...)
}

// FetchRows the fetch data rows
func (r *DB) Query(query string, args ...interface{}) (rs *sql.Rows, err error) {
	if err = r.callBreaker(func() error {
		if r.DB == nil {
			err = errors.New("the database connection is nil")
			logrus.Error(err, logrus.Fields{
				"query": query,
				"args":  args,
			})
			return err
		}
		if rs, err = r.DB.Query(query, args...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logrus.Error(err.Error(), logrus.Fields{
			"query": query,
			"args":  args,
		})
	}
	return rs, err
}

// FetchRow the fetch data row
func (r *DB) QueryRow(query string, args ...interface{}) (rs *sql.Row, err error) {
	if err = r.callBreaker(func() (err error) {
		if r.DB == nil {
			err = errors.New("the database connection is nil")
			logrus.Error(err, logrus.Fields{
				"query": query,
				"args":  args,
			})
			return err
		}
		rs = r.DB.QueryRow(query, args...)
		return nil
	}); err != nil {
		logrus.Error(err, logrus.Fields{
			"query": query,
			"args":  args,
		})
	}
	return rs, err
}

// SetCommandBreaker the circuit breaker
func (r *DB) SetCommandBreaker(commandName string, timeout, maxConcurrent int, args ...interface{}) *DB {

	r.commandName = &commandName
	r.retryCount = r.config.CB.Retry
	if len(args) == 1 {
		switch args[0].(type) {
		case fallbackFunc:
			r.fallbackFunc = args[0].(fallbackFunc)
		}
	}

	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		MaxConcurrentRequests: maxConcurrent,
		Timeout:               timeout,
		ErrorPercentThreshold: 25,
	})

	return r
}

// callBreaker command circuit breaker
func (r *DB) callBreaker(fn func() error) error {
	var err error
	if r.commandName == nil {
		return fn()
	}
	cn := *r.commandName
	for i := 0; i <= r.retryCount; i++ {
		err = hystrix.Do(cn, func() error {
			return fn()
		}, r.fallbackFunc)
		if err != nil {
			var backOffTime time.Duration
			if i <= 0 {
				backOffTime = 0 * time.Millisecond
			} else {
				// rand.Int63n(nc.interval*1000)
				backOffTime = (time.Duration(int64(2/time.Millisecond)) * time.Millisecond) + (time.Duration(rand.Int63n(5*1000)) * time.Millisecond)
			}
			time.Sleep(backOffTime)
			continue
		}
		break
	}
	return err
}

// GetQueryTimeout for circuit breaker
func (r *DB) GetQueryTimeout() int {
	if timeout := r.config.CB.Timeout; timeout > 1 {
		return timeout
	}
	return 1000
}

// GetDefaultMaxConcurrent circuit breaker
func (r *DB) GetDefaultMaxConcurrent() int {
	if concurrent := r.config.CB.Concurrent; concurrent > 1 {
		return concurrent
	}
	return 100
}
