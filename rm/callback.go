package rm

var commands = make([]*Command, 0)
var moduleTypes = make([]*ModuleType, 0)

func getCommand(id int) *Command {
    return commands[id]
}
func getModuleType(id int) *ModuleType {
    return moduleTypes[id]
}
func commandId(cmd*Command) int {
    id := len(commands)
    commands = append(commands, cmd)
    return id
}

func moduleTypeId(mt*ModuleType) int {
    id := len(moduleTypes)
    moduleTypes = append(moduleTypes, mt)
    return id
}
