# Módulo Go Redis
go-rm te permitirá escribir módulos Redis, utilizando Golang.

Leer en | [中文](./README-zh_CN.md) | [Inglés](./README.md)
----|----

## Demostración

```bash
# Asegurate de tener la última versión de Redis
# Por ejemplo, con brew puedes descargarla utilizando:
# brew reinstall redis --HEAD

# Para compilar el módulo
go get -u -v -buildmode=c-shared github.com/wenerme/go-rm/modules/hashex

# Para iniciar redis-server y cargar nuestro módulo, con logging en modo de depuración:
redis-server --loadmodule ~/go/pkg/*/github.com/wenerme/go-rm/modules/hashex* --loglevel debug
```

__Conectarse a redis-server__

```
$ redis-cli hsetget a a 5
(nil)
$ redis-cli hsetget a a 4
"5"
$ redis-cli hsetget a a 3
"3"
```

Vaya, funciona, ahora puedes distribuir este módulo Redis a tus amigos. :P

## ¿Cómo escribir un módulo Redis?

Implementar un módulo Redis es tan fácil como crear una aplicación de consola en Go, esto es todo lo que necesitas para implementar el comando de arriba, el código fuente está [aquí](https://github.com/wenerme/go-rm/blob/master/modules/hashex/hashex.go).

```go
package main

import "github.com/wenerme/go-rm/rm"

func main() {
    // En caso de que alguien intente llamar esto.
    rm.Run()
}

func init() {
    rm.Mod = CreateMyMod()
}
func CreateMyMod() *rm.Module {
    mod := rm.NewMod()
    mod.Name = "hashex"
    mod.Version = 1
    mod.Commands = []rm.Command{
        {
            Name:   "hsetget",
            Flags:  "write fast deny-oom",
            FirstKey:1, LastKey:1, KeyStep:1,
            Action: func(cmd rm.CmdContext) int {
                ctx := cmd.Ctx
                if len(cmd.Args) != 4 {
                    return ctx.WrongArity()
                }
                ctx.AutoMemory()
                // abrir la clave y asegurarse de que se trata de un hash y no está vacía
                key := ctx.OpenKey(cmd.Args[1], rm.READ | rm.WRITE)
                if key.KeyType() != rm.KEYTYPE_EMPTY && key.KeyType() != rm.KEYTYPE_HASH {
                    ctx.ReplyWithError(rm.ERRORMSG_WRONGTYPE)
                    return rm.ERR
                }
                // obtener el valor actual del elemento hash
                var val rm.String;
                key.HashGet(rm.HASH_NONE, cmd.Args[2], (*uintptr)(&val), rm.NullString())
                // definir el nuevo valor
                key.HashSet(rm.HASH_NONE, cmd.Args[2], cmd.Args[3], rm.NullString())
                if val.IsNull() {
                    ctx.ReplyWithNull()
                } else {
                    ctx.ReplyWithString(val)
                }
                return rm.OK
            },
        },
    }
    return mod
}
```

## Fantasía

* Un módulo de gestión de módulos, provee
    * mod.search
        * Búsqueda de módulos en repositorios (¿GitHub?)
        * La estructura del repositorio sería así:
        ```
        /namespace
            /module-name
                /bin
                    /darwin_amd64
                        module-name.so
                        module-name.sha
                    /linux_amd64
                module-name.go     
        ```
    * mod.get
        * Descargar el módulo a ~/.redismodule
        * Dado que el módulo está escrito en Go, podemos compilarlo para casi cualquier plataforma
        * Podemos utilizar el tag/commit para versionar el binario, entonces también sería posible descargar versiones anteriores
    * mod.install
        * Instalar el módulo descargado
    * ...
* Un módulo para gestión de cluster
    * Facilitaría la creación/gestión/monitoreo de un cluster basado en Redis 3
* Un tipo de dato JSON, para demostrar la forma de agregar un nuevo tipo de datos a Redis.
    * json.fmt key template
    * json.path key path \[pretty]
    * json.get key \[pretty]
    * json.set key value
        * Esto validaría el JSON

## Dificultades
* El código C no puede llamar a funciones Go, entonces cada callback es pregenerado
    * 200 comandos como máximo
    * 5 tipos de datos como máximo
    * Los límites son fáciles de cambiar, sólo necesitan un valor máximo apropiado
* Go no puede llamar a var_args, la llamada a esta función también es pregenerada
    * HashSet/HashGet acepta hasta 20 argumentos
    * Los límites son fáciles de cambiar, sólo necesitan un valor máximo apropiado

## TODO

* Encontrar límites apropiados para tipos de datos y var_args
