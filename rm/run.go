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
                //mod := *Mod
                //mod.BeforeInit = nil
                //mod.AfterInit = nil
                //for i := range mod.Commands {
                //    mod.Commands[i].Action = nil
                //}
                //for i := range mod.DataTypes {
                //    mod.DataTypes[i].RdbLoad = nil
                //    mod.DataTypes[i].RdbSave = nil
                //    mod.DataTypes[i].AofRewrite = nil
                //    mod.DataTypes[i].Digest = nil
                //    mod.DataTypes[i].Free = nil
                //}
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
    app.Run(strings.Split("redismodule info", " "))
    //app.RunAndExitOnError()
}
