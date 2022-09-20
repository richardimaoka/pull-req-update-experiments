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

func genAbcFile(filename string) string {
	return f(
		`cat << EOF > %s
a

b

c
EOF`, filename)
}

func WriteShellScript(w io.Writer, commands []fmt.Stringer) {
	fmt.Fprint(w, "#!/bin/sh\n\n")
	for _, cmd := range commands {
		fmt.Fprintf(w, "%v\n\n", cmd.String())
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

func main() {
	var commands []fmt.Stringer
	f := fmt.Sprintf

	// commands = append(commands, &MultiCommands{
	// 	Comment: "準備: GitHub レポジトリの作成",
	// 	Commands: []string{
	// 		"mkdir pull-req-update-experiments",
	// 		"cd pull-req-update-experiments",
	// 		"git init",
	// 	},
	// })

	// commands = append(commands, &SingleCommand{
	// 	Comment: "GitHub repository create",
	// 	Command: "gh repo create pull-req-update-experiments --public --source=. --remote=origin",
	// })

	mainBranch := "developer"
	if mainBranch != "main" {
		commands = append(commands, &SingleCommand{
			Command: f("git branch -f %s", mainBranch),
		})
	}

	filename := "pull-req-no-conflict.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			f(`git switch %s`, mainBranch),
			genAbcFile(filename),
			"git add --all",
			f(`git commit -m "create %s"`, filename),
			f(`git push origin %s`, mainBranch),
		},
	})

	prBranch := "pr-update-1"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			f(`git switch -c %s`, prBranch),
			f(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update a in %s"`, prBranch),
			f(`git push --set-upstream origin %s`, prBranch),
			f(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			f(`git switch %s`, mainBranch),
			f(`sed -i 's/b/bbbbb/' %s # ファイル中のbをbbbbbに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update b in %s"`, mainBranch),
			f("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "PRをマージ",
		Command: f("gh pr merge %s --merge --delete-branch", prBranch),
	})

	filename = "pull-req-same-line-conflict.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			f(`git switch %s`, mainBranch),
			genAbcFile(filename),
			"git add --all",
			f(`git commit -m "create %s"`, filename),
			f(`git push origin %s`, mainBranch),
		},
	})

	prBranch = "pr-update-2"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			f(`git switch -c %s`, prBranch),
			f(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update a to aaaaa in %s"`, prBranch),
			f(`git push --set-upstream origin %s`, prBranch),
			f(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			f(`git switch %s`, mainBranch),
			f(`sed -i 's/a/aaa/' %s # ファイル中のaをaaaに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update a to aaa in %s"`, mainBranch),
			f("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "PRをマージ",
		Command: f("gh pr merge %s --merge --delete-branch", prBranch),
	})
	// bytes, err := json.Marshal(commands)

	RunCommands(commands)
	//WriteShellScript(os.Stdout, commands)
}
