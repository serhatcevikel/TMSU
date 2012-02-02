/*
Copyright 2011 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

type HelpCommand struct{}

func (HelpCommand) Name() string {
	return "help"
}

func (HelpCommand) Synopsis() string {
	return "List commands or show help for a particular command"
}

func (HelpCommand) Description() string {
	return `tmsu help [OPTION] [COMMAND]

Shows help summary or, where COMMAND is specified, help for COMMAND.

    --list    list commands`
}

func (command HelpCommand) Exec(args []string) error {
	switch len(args) {
	case 0:
		command.summary()
	default:
	    if (args[0] == "--list") {
	        command.listCommands()
        } else {
		    command.describeCommand(args[0])
        }
	}

	return nil
}

func (HelpCommand) summary() {
	fmt.Println("TMSU")
	fmt.Println()

	var maxWidth int = 0
	commandNames := make([]string, 0, len(commands))
	for _, command  := range commands {
        commandName := command.Name()
		maxWidth = int(math.Max(float64(maxWidth), float64(len(commandName))))
		commandNames = append(commandNames, commandName)
	}

	sort.Strings(commandNames)

	for _, commandName := range commandNames {
		command, _ := commands[commandName]

		commandSummary := command.Synopsis()
		if commandSummary == "" {
			continue
		}

		fmt.Printf("  %-"+strconv.Itoa(maxWidth)+"v  %v\n", command.Name(), commandSummary)
	}
}

func (HelpCommand) listCommands() {
	for _, command  := range commands {
	    if command.Synopsis() == "" {
	        continue
        }

        fmt.Println(command.Name())
	}
}

func (HelpCommand) describeCommand(commandName string) {
	command := commands[commandName]
	if command == nil {
		fmt.Printf("No such command '%v'.\n", commandName)
		return
	}

	fmt.Println(command.Description())
}
