package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type SingleCommand struct {
	Title   string
	Comment string
	Command string
}

func (c *SingleCommand) String() string {
	if c.Comment == "" {
		return c.Command
	} else {
		return "# " + c.Comment + "\n" + c.Command
	}
}

type MultiCommands struct {
	Title    string
	Comment  string
	Commands []string
}

func (c *MultiCommands) String() string {
	if c.Comment == "" {
		return strings.Join(c.Commands, "\n")
	} else {
		return "# " + c.Comment + "\n" + strings.Join(c.Commands, "\n")
	}
}

func WriteMarkdown(w io.Writer, cmdBlocks []fmt.Stringer) {
	for _, cmdBlk := range cmdBlocks {
		var commands []string
		switch v := cmdBlk.(type) {
		case *SingleCommand:
			commands = append(commands, v.Command)
		case *MultiCommands:
			commands = append(commands, v.Commands...)
		}

		fmt.Fprintln(w, "```sh:コピペして実行")
		for _, cmdString := range commands {
			fmt.Fprintln(w, cmdString)
		}
		fmt.Fprint(w, "```\n\n")
	}
}

func RunCommands(cmdBlocks []fmt.Stringer) {
	input := bufio.NewScanner(os.Stdin)

	for _, cmdBlk := range cmdBlocks {
		var commands []string
		switch v := cmdBlk.(type) {
		case *SingleCommand:
			commands = append(commands, v.Command)
		case *MultiCommands:
			commands = append(commands, v.Commands...)
		}

		for _, cmdString := range commands {
			fmt.Println("### Executing the following command ###")
			fmt.Println(cmdString)
			fmt.Print("[y/n] ")

			input.Scan()
			text := input.Text()
			switch text {
			case "y":
				fmt.Println("executing")
				execCmd := exec.Command("sh", "-c", cmdString)
				output, _ := execCmd.CombinedOutput()
				fmt.Println(string(output))
			case "n":
				fmt.Println("skipping")
			default:
				fmt.Print("[y/n] ")
			}
		}
	}

}
