package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
)

const (
	AppName    = "envtpl"
)

// goreleaser convention, https://goreleaser.com/customization/builds/go/
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// options
var (
	missingKey  = "default"
	output      string
	showVersion bool
)

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var RootCmd = &cobra.Command{
	Use:   "envtpl",
	Short: "Render go templates from environment variables",
	Long:  `Render go templates from environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("%s %s (commit: %s, built at: %s, Go runtime: %s)\n",
				AppName, version, commit, date, runtime.Version())
			os.Exit(0)
		}

		var t *template.Template
		var err error

		// load template; if an argument is not specified, default to stdin
		if len(args) > 0 {
			t, err = parseFiles(args...)
			die(err)
		} else {
			bytes, err := io.ReadAll(os.Stdin)
			die(err)
			t, err = parse(string(bytes))
			die(err)
		}

		// get environment variables to supply to the template
		env := readEnv()

		// get writer for rendered output; if an output file is not
		// specified, default to stdout
		var w io.Writer
		if output == "" {
			w = os.Stdout
		} else {
			f, err := os.Create(output)
			die(err)
			defer f.Close()
			w = io.Writer(f)
		}

		// set error handling strategy for missing keys
		if missingKey != "default" {
			t = t.Option("missingkey=" + missingKey)
		}

		// render the template
		err = t.Execute(w, env)
		die(err)
	},
}

// returns a new template with custom function maps
// (sprig and environment key prefix matcher) applied
func parse(s string) (*template.Template, error) {
	return template.New("").Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).Parse(s)
}

// returns a new template with custom function maps
// (sprig and environment key prefix matcher) applied
func parseFiles(files ...string) (*template.Template, error) {
	return template.New(filepath.Base(files[0])).Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).ParseFiles(files...)
}

// returns map of environment variables
func readEnv() (env map[string]string) {
	env = make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		env[pair[0]] = pair[1]
	}
	return
}

// custom function that returns key, value for all environment variable keys matching prefix
// (see original envtpl: https://pypi.org/project/envtpl/)
func environment(prefix string) map[string]string {
	env := make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

// returns custom template functions map
func customFuncMap() template.FuncMap {
	var functionMap = map[string]interface{}{
		"environment": environment,
	}
	return template.FuncMap(functionMap)
}

func main() {
	RootCmd.Flags().StringVarP(&missingKey, "missingkey", "m", missingKey, "Strategy for dealing with missing keys: [default|zero|error]")
	RootCmd.Flags().StringVarP(&output, "output", "o", output, "The rendered output file")
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version information")

	err := RootCmd.Execute()
	die(err)
}
