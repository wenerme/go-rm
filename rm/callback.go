package rm

var commands = make([]*Command, 0)
var moduleTypes = make([]*DataType, 0)

func getCommand(id int) *Command {
    return commands[id]
}
func getDataType(id int) *DataType {
    return moduleTypes[id]
}
func commandId(cmd*Command) int {
    id := len(commands)
    commands = append(commands, cmd)
    return id
}

func dataTypeId(mt*DataType) int {
    id := len(moduleTypes)
    moduleTypes = append(moduleTypes, mt)
    return id
}
