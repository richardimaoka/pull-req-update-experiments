package main

import "fmt"

type Command struct {
	command string
}

func genAbcFile(filename string) Command {
	commandString := fmt.Sprintf(
`cat << EOF > %s.txt
a

b

c
EOF`, filename)

	return Command{
		command: commandString
	}
}

func main() {
	var commands []Command
	commands = append(commands, genAbcFile())
	fmt.Println(abcText("pull-req-update-commit"))
}
