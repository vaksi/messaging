/*  configuration.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 10, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 10/10/18 17:25
 */

package configuration

import (
    "fmt"
    "io/ioutil"
    "os"
    "testing"

    . "github.com/smartystreets/goconvey/convey"
    "github.com/spf13/viper"
)

func TestInitConfiguration(t *testing.T) {
    Convey("InitConfiguration should be success", t, func() {
        pathDir := "./configs"
        // Create Dir File
        if _, err := os.Stat(pathDir); os.IsNotExist(err) {
            os.Mkdir(pathDir, os.ModePerm)
        }

        // Create File Config
        // any approach to require this configuration into your program.
        var yamlExample = []byte(`
            Hacker: true
            name: steve
            hobbies:
                - skateboarding
                - snowboarding
                - go
            clothing:
                jacket: leather
                trousers: denim
            age: 35
            eyes : brown
            beard: true
        `)
        err := ioutil.WriteFile(pathDir + "/app.yaml", yamlExample, os.ModePerm)
        if err != nil {
            fmt.Println(err)
        }

        So(InitConfiguration("configs", "app"), ShouldBeNil)

        So(viper.GetString("Hacker"), ShouldEqual, "true")
    })

    Convey("InitConfiguration should be Error", t, func() {
        So(InitConfiguration("configs-false", "app-false"), ShouldBeError)
    })
}
