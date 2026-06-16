//+build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}	

func init() {
    build.SetBeforeBuildCallback(func(c build.Config) (build.Config, error) {
        fmt.Printf("Running before build callback with config: %+v\n", c)
        c.EnableCGo = true
        return c, nil
    })
}

// Default configures the default target.
var Default = build.BuildAll
