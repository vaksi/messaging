package cmd

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/vaksi/messaging/configs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var (
	wg sync.WaitGroup
)

func TestNewHttp(t *testing.T) {
	assert := require.New(t)
	stop := make(chan bool)
	wg.Add(1)
	configuration := configs.New("app.test",
		"./configs", "../configs", "../../../configs")

	http := NewHttpCmd(configuration)
	http.stop = stop

	go func() {
		defer wg.Done()
		err := http.BaseCmd.Execute()
		assert.NoError(err)
	}()

	stop <- true
	wg.Wait()
}

func TestHttp(t *testing.T) {
	assert := require.New(t)
	stop := make(chan bool)
	wg.Add(1)
	configuration := configs.New("app.test",
		"./configs", "../configs", "../../../configs")

	cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd
	go func() {
		defer wg.Done()
		_, err := cmd.ExecuteC()
		assert.NoError(err)
	}()

	stop <- true
	wg.Wait()
}

func TestHttpFail(t *testing.T) {
	var (
		err           error
		stop          = make(chan bool)
		configuration *configs.Config
	)

	wg.Add(1)
	f := func() {
		configuration = configs.New("appss",
			"./configs", "../configs", "../../configs")
	}
	assert.Panics(t, f)
	cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd

	go func() {
		defer wg.Done()
		fn := func() {
			_, err = cmd.ExecuteC()
		}
		assert.Panics(t, fn)
		<-stop
	}()

	assert.NoError(t, err)
	stop <- true
	wg.Wait()
}

func TestNewHttpCmdWithFilename(t *testing.T) {
	stop := make(chan bool)
	configuration := configs.New("app.test", ConfigPath...)

	wg.Add(1)
	os.Args = []string{"main", "http", "-f", "app"}

	cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd
	go func() {
		defer wg.Done()

		_, err := cmd.ExecuteC()
		assert.NoError(t, err)
	}()

	stop <- true
	wg.Wait()
	os.Args = []string{""}
}

func TestListenAndServe(t *testing.T) {
	var err error
	stop := make(chan bool)

	configuration := configs.New("app.test", ConfigPath...)

	cc := &httpCmd{stop: stop}
	cc.configuration = configuration
	cc.BaseCmd = &cobra.Command{
		Use:   "http",
		Short: "Used to run the http service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			mux := new(mux.Router)
			return cc.serve(mux)
		},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = cc.BaseCmd.Execute()
	}()
	assert.NoError(t, err)
	stop <- true
	wg.Wait()
}
