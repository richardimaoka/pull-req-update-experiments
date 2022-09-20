package main

import (
	"fmt"
	"strings"
)

type CommandStringer interface {
	String() string
}

type SingleCommand struct {
	Title   string
	Comment string
	Command string
}

func (c *SingleCommand) String() string {
	if c.Comment == "" {
		return c.Command
	} else {
		return c.Comment + "\n" + c.Command
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
		return c.Comment + "\n" + strings.Join(c.Commands, "\n")
	}
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
	var commands []CommandStringer

	commands = append(commands, &MultiCommands{
		Comment: "# 準備: GitHub レポジトリの作成",
		Commands: []string{
			"mkdir pull-req-update-experiments",
			"cd pull-req-update-experiments",
			"git init",
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "# GitHub repository create",
		Command: "gh repo create pull-req-update-experiments --public --source=. --remote=origin",
	})

	mainBranch := "developer"
	filename := "pull-req-no-conflict.txt"
	commands = append(commands, &MultiCommands{
		Comment: "# 準備: GitHub レポジトリの初期構成",
		Commands: []string{
			fmt.Sprintf(`git switch %s`, mainBranch),
			genAbcFile(filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "%s"`, filename),
			fmt.Sprintf(`git push origin %s`, mainBranch),
		},
	})

	pr1branch := "pr-update-1"
	commands = append(commands, &MultiCommands{
		Commands: []string{
			fmt.Sprintf(`git switch %s`, pr1branch),
			fmt.Sprintf(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "%s"`, filename),
			fmt.Sprintf(`git push --set-upstream origin "%s"`, pr1branch),
			fmt.Sprintf(`gh pr create --title %s --body "" --base %s --head %s`, pr1branch, mainBranch, pr1branch),
		},
	})

	// bytes, err := json.Marshal(commands)

	fmt.Printf("#!/bin/sh\n\n")
	for _, cmd := range commands {
		fmt.Printf("%v\n\n", cmd.String())
	}
}
