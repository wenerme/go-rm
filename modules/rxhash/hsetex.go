package rxhash

import "github.com/wenerme/go-rm/rm"

func init() {
	commands = append(commands, CreateCommand_HSETEX())
}
func CreateCommand_HSETEX() rm.Command {
	return rm.Command{
		Desc: `HSETEX key field value
Set field to value ony if field is already exists`,
		Name:   "hsetex",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(args) != 4 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()

			key := ctx.OpenKey(args[1], rm.WRITE)
			if key.KeyType() != rm.KEYTYPE_EMPTY && key.KeyType() != rm.KEYTYPE_HASH {
				ctx.ReplyWithError(rm.ERRORMSG_WRONGTYPE)
				return rm.ERR
			}
			ctx.ReplyWithLongLong(int64(key.HashSet(rm.HASH_XX, args[2], args[3])))
			return rm.OK
		},
	}
}
