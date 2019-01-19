package rm

import (
	"bytes"
	"strings"
)

// Flags for create command
type CmdFlag string

const (
	// The command may modify the data set (it may also read from it).
	CF_WRITE CmdFlag = "write"

	// The command returns data from keys but never writes.
	CF_READONLY = "readonly"

	// The command is an administrative command (may change replication or perform similar tasks).
	CF_ADMIN = "admin"

	// The command may use additional memory and should be denied during out of memory conditions.
	CF_DENY_OOM = "deny-oom"

	// Don't allow this command in Lua scripts.
	CF_DENY_SCRIPT = "deny-script"

	// Allow this command while the server is loading data.
	// Only commands not interacting with the data set
	// should be allowed to run in this mode. If not sure
	// don't use this flag.
	CF_ALLOW_LOADING = "allow-loading"

	// The command publishes things on Pub/Sub channels.
	CF_PUBSUB = "pubsub"

	// The command may have different outputs even starting from the same input arguments and key values.
	CF_RANDOM = "random"

	// The command is allowed to run on slaves that don't
	// serve stale data. Don't use if you don't know what
	// this means.
	CF_ALLOW_STALE = "allow-stale"

	// Don't propoagate the command on monitor. Use this if the command has sensible data among the arguments.
	// The command time complexity is not greater
	CF_NO_MONITOR = "no-monitor"

	//The command time complexity is not greater
	//than O(log(N)) where N is the size of the collection or
	//anything else representing the normal scalability
	//issue with the command.
	CF_FAST = "fast"

	//The command implements the interface to return
	//the arguments that are keys. Used when start/stop/step
	//is not enough because of the command syntax.
	CF_GETKEYS_API = "getkeys-api"

	//The command should not register in Redis Cluster
	//since is not designed to work with it because, for
	//example, is unable to report the position of the
	//keys, programmatically creates key names, or any
	//other reason.
	CF_NO_CLUSTER = "no-cluster"
)

func BuildCommandFlag(f ...CmdFlag) string {
	buf := bytes.NewBufferString("")
	for _, v := range f {
		buf.WriteString(string(v))
		buf.WriteRune(' ')
	}
	flags := strings.TrimRight(buf.String(), " ")
	return flags
}
