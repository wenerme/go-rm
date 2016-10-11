package main

import (
    "github.com/urfave/cli"
    "io/ioutil"
    "text/template"
    "bytes"
    "os"
    "github.com/bradfitz/iter"
    "regexp"
    "strings"
    "reflect"
    "fmt"
    "github.com/ngaut/log"
)

func main() {
    app := cli.NewApp()
    app.Name = "gen"
    app.Commands = []cli.Command{
        {
            Name: "callback",
            Usage: "Generate Callbacks",
            Action: GenerateCallback,
        },
        {
            Name: "wrapper",
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:"t",
                    Usage:`Generate wrapper c or go`,
                    Value:"c",
                },
                cli.StringFlag{
                    Name:"f",
                    Usage:`Generated filename`,
                },
            },
            Action: GenerateWrapper,
        },
    }
    //app.Run(strings.Split("gen wrapper -t go -f ./rm/wrapper.go", " "))
    app.Run(strings.Split("gen callback", " "))
    //app.RunAndExitOnError()
    //b, err := ioutil.ReadFile("./API.md")
    //if err != nil {
    //    panic(err)
    //}
    //
    //b, err = json.Marshal(ParseApi(string(b)))
    //if err != nil {
    //    panic(err)
    //}
    //ioutil.WriteFile("api.json", b, os.ModePerm)
}

var callbackTemplate = `
#ifndef GO_RM_CALLBACK_H
#define GO_RM_CALLBACK_H

#include "./decl.h"

{{range $k, $v := .CallTypes}}

// {{.Name}} - {{.Type}}
inline int {{.Name}}_size(){return {{.Number}};}
{{/*定义回调函数*/ -}}

{{range $i, $_ := N .Number}}
{{- with $v -}}

    {{.Ret}} cb_{{.Name}}_{{$i}}({{.ParamString}}){
        {{- if ne .Ret "void"}}return {{end}}
        {{- /*调用实际方法*/ -}}
        {{.Name}}_call({{$i}},{{.ParamName}});{{- "" -}}
    }
{{end}}
{{- end}}


{{- /*定义变量*/ -}}

const {{.Type}} cb_{{.Name}}[] = {
{{range $i, $_ := N .Number}}
{{- with $v -}}
    cb_{{.Name}}_{{$i}},{{if neednl $i }}{{"\n"}}{{end}}
{{- end}}
{{- end}}
};

{{- end}}
#endif
`

type CallbackCtx struct {
    CmdFuncNumber int
    ModTypeNumber int
    CallTypes     []*CallType
}

type CallType struct {
    Name        string
    Type        string
    Ret         string
    ParamString string
    ParamName   string
    Number      int
}

func GenerateCallback(ctx *cli.Context) error {
    tpl := MustTemplate(callbackTemplate, template.FuncMap{
        "neednl": func(i int) bool {
            return i % 15 == 0 && i != 0;
        }})

    buf := bytes.NewBufferString("")
    callbackCtx := CallbackCtx{
        CmdFuncNumber: 60,
        ModTypeNumber: 6,
        // mt_rdb_load RedisModuleTypeLoadFunc
        // mt_rdb_save RedisModuleTypeSaveFunc
        // mt_aof_rewrite RedisModuleTypeRewriteFunc
        // mt_digest RedisModuleTypeDigestFunc
        // mt_free RedisModuleTypeFreeFunc

        // typedef int (*RedisModuleCmdFunc) (RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
        // typedef void *(*RedisModuleTypeLoadFunc)(RedisModuleIO *rdb, int encver);
        // typedef void (*RedisModuleTypeSaveFunc)(RedisModuleIO *rdb, void *value);
        // typedef void (*RedisModuleTypeRewriteFunc)(RedisModuleIO *aof, RedisModuleString *key, void *value);
        // typedef void (*RedisModuleTypeDigestFunc)(RedisModuleDigest *digest, void *value);
        // typedef void (*RedisModuleTypeFreeFunc)(void *value);
        CallTypes:[]*CallType{
            {
                Name: "cmd_func",
                Type: "RedisModuleCmdFunc",
                ParamString:"RedisModuleCtx *ctx, RedisModuleString **argv, int argc",
                ParamName: "ctx, argv, argc",
                Ret: "int",
                Number: 200,
            },
            {
                Name: "mt_rdb_load",
                Type: "RedisModuleTypeLoadFunc",
                ParamString: "RedisModuleIO *rdb, int encver",
                ParamName: "rdb, encver",
                Ret : "void*",
            },
            {
                Name: "mt_rdb_save",
                Type: "RedisModuleTypeSaveFunc",
                ParamString: "RedisModuleIO *rdb, void *value",
                ParamName: "rdb, value",
                Ret : "void",
            },
            {
                Name: "mt_aof_rewrite",
                Type: "RedisModuleTypeRewriteFunc",
                ParamString: "RedisModuleIO *aof, RedisModuleString *key, void *value",
                ParamName: "aof, key, value",
                Ret : "void",
            },
            {
                Name: "mt_digest",
                Type: "RedisModuleTypeDigestFunc",
                ParamString: "RedisModuleDigest *digest, void *value",
                ParamName: "digest, value",
                Ret : "void",
            },
            {
                Name: "mt_free",
                Type: "RedisModuleTypeFreeFunc",
                ParamString: "void *value",
                ParamName: "value",
                Ret : "void",
            },
        },
    }
    for _, c := range callbackCtx.CallTypes {
        if c.Number != 0 {
            continue
        }
        if c.Name[:2] == "mt" {
            c.Number = 5
        }
    }
    err := tpl.Execute(buf, callbackCtx)
    if err != nil {
        log.Fatalf("Execute template failed: %v", err)
        return err
    }
    return ioutil.WriteFile("./rm/callbacks.h", buf.Bytes(), os.ModePerm)
}

var wrapperHeaderTemplate = `
#ifndef GO_RM_WRAPPER_H
#define GO_RM_WRAPPER_H

#include "./redismodule.h"

{{range $i, $v := .Apis}}
{{if show_api $v -}}
{{commented .Desc}}
{{.Ret}} {{.Name}}(
    {{- range $j, $arg := .ArgInfos -}}
        {{.Type}} {{.Name}}{{if not (last $j $v.ArgInfos)}},{{end}}
    {{- end -}}
){
    {{- if need_ret $v}}return {{end -}}
    RedisModule_{{.Name}}(
        {{- range $j, $arg := .ArgInfos -}}
            {{$arg.Name}}{{if not (last $j $v.ArgInfos)}},{{end}}
        {{- end -}}
        );}
{{- end}}
{{end}}

#endif
`

var wrapperGoTemplate = `
package rm

//#include "./rm.h"
import "C"
import (
    "unsafe"
)

{{range $i, $v := .Apis}}
{{if show_api $v -}}
{{commented .Desc}}
// {{.Sig}}
func {{.Name}}(
    {{- range $j, $arg := .ArgInfos -}}
       {{.Name}} {{gotype .Type}}{{if not (last $j $v.ArgInfos)}},{{end}}
    {{- end -}}
)(
    {{- if need_ret $v}}{{gotype .Ret}}{{end -}}
){
    {{- if need_ret $v}}return {{end -}}
    {{- /* 类型转换 */ -}}
    {{- if need_ret $v}}{{gotype .Ret}}({{end -}}
    {{- /* 方法调用 */ -}}
    C.{{.Name}}(
        {{- range $j, $arg := .ArgInfos -}}
            {{.Name}}{{if not (last $j $v.ArgInfos)}},{{end}}
        {{- end -}}

    )
    {{- if need_ret $v}}){{end -}}
}
{{- end}}
{{end}}
`

type WrapperCtx struct {
    Apis []ApiInfo
}

func GenerateWrapper(ctx *cli.Context) error {
    tplString := wrapperHeaderTemplate
    fn := "./rm/wrapper.h"
    switch ctx.String("t") {
    case "go":
        fn = "./rm/wrapper.go"
        tplString = wrapperGoTemplate
    }
    if ctx.String("f") != "" {
        fn = ctx.String("f")
    }

    tpl := MustTemplate(tplString, template.FuncMap{
        "need_ret": func(i ApiInfo) bool {
            return i.Ret != "void"
        },
        "show_api": func(i ApiInfo) bool {
            switch i.Name {
            case "ZsetAddFlagsToCoreFlags":
            case "ZsetAddFlagsFromCoreFlags":
            case "FreeCallReply_Rec":
            default:
                return true
            }
            return false
        },
    })
    b, err := ioutil.ReadFile("./API.md")
    if err != nil {
        return err
    }
    apis := ParseApi(string(b))
    context := WrapperCtx{
        Apis: apis,
    }
    str := MustExecute(tpl, context)
    return ioutil.WriteFile(fn, []byte(str), os.ModePerm)
}

type ApiInfo struct {
    Sig      string
    Ret      string
    Name     string
    Args     string
    Desc     string
    ArgInfos []ArgInfo
}
type ArgInfo struct {
    Type string
    Name string
}

func ParseApi(md string) []ApiInfo {
    apis := make([]ApiInfo, 0)
    parts := regexp.MustCompile(`(?m)^##.*$`).Split(md, -1)
    sig := regexp.MustCompile(`(?m)^\s*$\s+([^;]+);\s*$^\s*$`)
    // void *RM_Alloc(size_t bytes);
    sigPattern := regexp.MustCompile(`(.*?)RM_([^(]+)\((.*?)\);`)
    for i, v := range parts {
        if i == 0 {
            continue
        }
        f := strings.TrimSpace(sig.FindString(v))
        match := sigPattern.FindStringSubmatch(f)
        apis = append(apis, ApiInfo{
            Sig: match[0],
            Ret: TypeMap(match[1]),
            Name:match[2],
            Args:match[3],
            Desc:strings.TrimSpace(sig.ReplaceAllString(v, "")),
            ArgInfos: ParseArgs(match[3]),
        })
    }
    return apis
}

func ParseArgs(args string) []ArgInfo {
    infos := make([]ArgInfo, 0)
    parts := regexp.MustCompile(`\s*,\s*`).Split(args, -1)
    namePattern := regexp.MustCompile(`(\w+)$`)

    for _, v := range parts {
        if v == "..." {
            continue
        }
        name := namePattern.FindStringSubmatch(v)[0]
        t := strings.TrimSpace(namePattern.ReplaceAllLiteralString(v, ""))
        infos = append(infos, ArgInfo{
            Name: name,
            Type: TypeMap(t),
        })
    }
    return infos
}

func TypeMap(t string) string {
    a := strings.TrimSpace(t)
    a = strings.Replace(a, " *", "*", -1)
    switch a {
    case "robj*":
        a = "RedisModuleString*"
    default:
        // moduleType
        if strings.HasPrefix(a, "moduleType") {
            // RedisModuleTypeSaveFunc
            a = "RedisModule" + a[len("module"):]
        }
    }
    return a
}
func GoTypeMap(f string) string {
    a := f
    a = strings.TrimSpace(strings.TrimPrefix(a, "const"))
    t := a;
    switch a {
    case "void*":
        t = "unsafe.Pointer"
    case "RedisModuleString*":
        t = "String"
    case "size_t":
        t = "int"
    case "size_t*":
        t = "*int"
    case "int*":
        t = "*int"
    case "RedisModuleCtx*":
        t = "Ctx"
    case "char*":
        t = "string"
    case "void**":
        t = "/* TODO void** */unsafe.Pointer"
    case "RedisModuleCmdFunc":
        t = "CmdFunc"
    case "RedisModuleCallReply*":
        t = "CallReply"
    case "RedisModuleKey*":
        t = "Key"
    case "RedisModuleIO*":
        t = "IO"
    case "RedisModuleType*":
        t = "/* TODO RedisModuleType* */unsafe.Pointer"
    case "mstime_t":
        t = "uint64"
    case "long long*":
        t = "*int64"
    case "double":
        t = "float64"
    case "double*":
        t = "*float64"
    case "unsigned long long":
        fallthrough
    case "uint64_t":
        t = "uint64"
    case "long long":
        fallthrough
    case "long":
        fallthrough
    case "int64_t":
        t = "int64"
    case "RedisModuleTypeLoadFunc":
    case "RedisModuleTypeSaveFunc":
    case "RedisModuleTypeRewriteFunc":
    case "RedisModuleTypeDigestFunc":
    case "RedisModuleTypeFreeFunc":
    case "int":
    default:
        fmt.Fprintf(os.Stderr, "No go type map found for %s(%v)\n", t, []byte(t))
    }

    return t
}

func MustTemplate(content string, funcs template.FuncMap) *template.Template {
    tpl := template.New("tpl")
    tpl.Funcs(template.FuncMap{
        "N": iter.N,
        "gotype": GoTypeMap,
        "last": func(x int, a interface{}) bool {
            return x == reflect.ValueOf(a).Len() - 1
        },
        "commented": func(s string) string {
            return regexp.MustCompile(`(?m)^`).ReplaceAllString(s, "// ")
        },
    })
    tpl.Funcs(funcs)

    template.Must(tpl.Parse(content))
    return tpl
}
func MustExecute(tpl *template.Template, data interface{}) string {
    buf := bytes.NewBufferString("")
    err := tpl.Execute(buf, data)
    if err != nil {
        panic(err)
    }
    return buf.String()
}
