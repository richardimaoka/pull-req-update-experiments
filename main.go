package main

import "fmt"

type Command struct {
	command string
}

func genAbcFile(filename string) string {
	return fmt.Sprintf(
		`cat << EOF > %s
a

b

c
EOF`, filename)
}

func main() {
	var commands []Command
	filename := "pull-req-update-commit.txt"
	commands = append(commands, Command{command: genAbcFile(filename)})
	commands = append(commands, Command{command: "git add --all"})
	commands = append(commands, Command{command: fmt.Sprintf(`git commit -m "%s"`, filename)})
	fmt.Println(commands)
}
