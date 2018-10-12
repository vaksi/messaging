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

    "github.com/spf13/viper"
)

func InitConfiguration(configPath, configName string) (err error) {
    viper.SetConfigName(configName) // name of config file (without extension)
    viper.AddConfigPath(configPath)   // path to look for the config file in
    // viper.AddConfigPath("$HOME/appname")  // call multiple times to add many search paths
    // viper.AddConfigPath(".")               // optionally look for config in the working directory
    err = viper.ReadInConfig() // Find and read the config file
    if err != nil { // Handle errors reading the config file
        return fmt.Errorf("Fatal error config file: %s \n", err)
    }

    return
}

// // initConfig initializes the configuration
// func initConfiguration() error {
//     viper.SetConfigName("App")
//     viper.AddConfigPath("configurations")
//     if err := viper.ReadInConfig(); err != nil {
//         return err
//     }
//     return nil
// }

// // SetConfiguration sets the configuration
// func SetConfiguration(param string) error {
//     if len(param) == 0 {
//         // Get default configuration file
//         if err := initConfiguration(); err != nil {
//             return fmt.Errorf("%v", err)
//         }
//         return nil
//     }
//
//     // Get file extension
//     ext := filepath.Ext(param)
//     ext = strings.TrimPrefix(ext, ".")
//     viper.SetConfigType(ext)
//
//     // Open configuration file
//     file, err := os.Open(AbsolutePath(param))
//     if err != nil {
//         return err
//     }
//     defer file.Close()
//     if err := viper.ReadConfig(file); err != nil {
//         return fmt.Errorf("%v", err)
//     }
//
//     return nil
// }
//
// func userHomeDir() string {
//     if runtime.GOOS == "windows" {
//         home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
//         if home == "" {
//             home = os.Getenv("USERPROFILE")
//         }
//         return home
//     }
//     return os.Getenv("HOME")
// }
//
// // AbsolutePath get absolute path
// func AbsolutePath(inPath string) string {
//     if strings.HasPrefix(inPath, "$HOME") {
//         inPath = userHomeDir() + inPath[5:]
//     }
//
//     if strings.HasPrefix(inPath, "$") {
//         end := strings.Index(inPath, string(os.PathSeparator))
//         inPath = os.Getenv(inPath[1:end]) + inPath[end:]
//     }
//
//     if filepath.IsAbs(inPath) {
//         return filepath.Clean(inPath)
//     }
//
//     p, err := filepath.Abs(inPath)
//     if err == nil {
//         return filepath.Clean(p)
//     }
//
//     return ""
// }
