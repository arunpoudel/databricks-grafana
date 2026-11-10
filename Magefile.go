//go:build mage
// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/grafana/grafana-plugin-sdk-go/build/buildinfo"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

// Default builds the standard plugin binaries and replaces the Linux/amd64
// binary with one that includes Databricks' Statement Execution API backend.
func Default() error {
	build.BuildAll()

	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil
	}

	pluginJSON, err := os.ReadFile("src/plugin.json")
	if err != nil {
		return err
	}
	var pluginMetadata struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(pluginJSON, &pluginMetadata); err != nil {
		return err
	}

	packageJSON, err := os.ReadFile("package.json")
	if err != nil {
		return err
	}
	var packageMetadata struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(packageJSON, &packageMetadata); err != nil {
		return err
	}

	info := buildinfo.Info{
		Time:     time.Now().UnixNano() / int64(time.Millisecond),
		PluginID: pluginMetadata.ID,
		Version:  packageMetadata.Version,
	}
	flags := make(map[string]string, 3)
	info.AppendFlags(flags)
	keys := make([]string, 0, len(flags))
	for key := range flags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ldflags := "-w -s -extldflags \"-static\""
	for _, key := range keys {
		ldflags += fmt.Sprintf(" -X '%s=%s'", key, flags[key])
	}

	cmd := exec.Command(
		"go",
		"build",
		"-tags=arrow_json_stdlib,databricks_kernel",
		"-ldflags", ldflags,
		"-o", "dist/gpx_databricks_linux_amd64",
		"./pkg",
	)
	cmd.Env = withEnvironmentVariable(os.Environ(), "CGO_ENABLED", "1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func withEnvironmentVariable(environment []string, name, value string) []string {
	prefix := name + "="
	for i, entry := range environment {
		if strings.HasPrefix(entry, prefix) {
			environment[i] = prefix + value
			return environment
		}
	}
	return append(environment, prefix+value)
}
