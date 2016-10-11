package rm

import (
	"github.com/urfave/cli"
	"text/template"
	"os"
	"encoding/json"
	"strings"
)

// Handler main invoke
func Run() {
	app := cli.NewApp()
	app.Name = Mod.Name
	app.Email = Mod.Email
	app.Author = Mod.Author
	app.Copyright = Mod.Copyright
	app.Version = Mod.SemVer
	app.Action = func(ctx *cli.Context) error {
		tpl, err := template.New("info").Parse(`{{"" -}}
Redis module {{.Name}} version {{.Version}} semver {{.SemVer}}
{{- if .Author}} created by {{.Author}}
    {{- if .Email}} <{{.Email}}>{{end -}}
{{- end}}
{{- if .Website}}
Know more from {{.Website}}{{end}}

This should be build by

    go build -v -buildmode=c-shared

Then you can load this redismodule by

    redis-server --loadmodule {{.Name}} --loglevel debug

`)
		if err != nil {
			return err
		}
		err = tpl.Execute(os.Stdout, Mod)
		if err != nil {
			return err
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:"info",
			Action:func(ctx *cli.Context) error {
				b, err := json.MarshalIndent(Mod, "  ", "  ")
				if err != nil {
					return err
				}
				os.Stdout.Write(b)
				os.Stdout.Sync()
				return nil
			},
		},
	}
	if (false) {
		// Test info output
		app.Run(strings.Split("redismodule info", " "))
	}
	app.RunAndExitOnError()
}
