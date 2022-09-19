package main

import (
	"fmt"
)

type Command struct {
	Command string
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
	commands = append(commands, Command{Command: "mkdir pull-req-update-experiments"})
	commands = append(commands, Command{Command: "cd pull-req-update-experiments"})
	commands = append(commands, Command{Command: "git init"})

	commands = append(commands, Command{Command: "gh repo create pull-req-update-experiments --public --source=. --remote=origin"})

	mainBranch := "developer"
	filename := "pull-req-no-conflict.txt"
	commands = append(commands, Command{Command: fmt.Sprintf(`git switch %s`, mainBranch)})
	commands = append(commands, Command{Command: genAbcFile(filename)})
	commands = append(commands, Command{Command: "git add --all"})
	commands = append(commands, Command{Command: fmt.Sprintf(`git commit -m "%s"`, filename)})
	commands = append(commands, Command{Command: fmt.Sprintf(`git push origin %s"`, mainBranch)})

	pr1branch := "pr-update-1"
	commands = append(commands, Command{Command: fmt.Sprintf(`git switch %s`, pr1branch)})
	commands = append(commands, Command{Command: fmt.Sprintf(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename)})
	commands = append(commands, Command{Command: "git add --all"})
	commands = append(commands, Command{Command: fmt.Sprintf(`git commit -m "%s"`, filename)})

	// bytes, err := json.Marshal(commands)

	for _, cmd := range commands {
		fmt.Printf("%v\n", cmd.Command)
	}
}
