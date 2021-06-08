package cmd

import (
	"encoding/json"
	"github.com/itxaka/luet-mtree/pkg/action"
	"os"
)

// Basic stub to call the action in the package, does nothing really
func newEventCmd(args []string) error {
	event := args[0]
	payload := args[1]

	eventDispatcher := action.NewEventDispatcherAction(event, payload)
	out, err := eventDispatcher.Run()

	// As this is part of being a luet plugin we need to output to console ONLY the results in json formatting so luet
	// can parse it.
	// Thankfully our eventDispatcher returns a nice map that can be dumped to json format easily :D
	outJson, _ := json.Marshal(out)
	os.Stdout.Write(outJson)

	// Let the root cmd be the one to set the exit status as success/failure
	return err
}