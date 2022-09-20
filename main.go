package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
	return fmt.Sprintf(
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
		default:
			log.Fatalf("cannot be true!!!")
		}

		for _, cmdString := range commands {
			fmt.Println("### Executing the following command ###")
			fmt.Println(cmdString)
			fmt.Print("[y/n] ")

			input.Scan()
			switch text := input.Text(); text {
			case "y":
				fmt.Println("executing")
				execCmd := exec.Command("sh", "-c", cmdString)
				output, _ := execCmd.CombinedOutput()
				fmt.Println(output)
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
	filename := "pull-req-no-conflict.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			fmt.Sprintf(`git switch %s`, mainBranch),
			genAbcFile(filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "create %s"`, filename),
			fmt.Sprintf(`git push origin %s`, mainBranch),
		},
	})

	prBranch := "pr-update-1"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			fmt.Sprintf(`git switch -c %s`, prBranch),
			fmt.Sprintf(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "update a in %s"`, prBranch),
			fmt.Sprintf(`git push --set-upstream origin %s`, prBranch),
			fmt.Sprintf(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: fmt.Sprintf("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			fmt.Sprintf(`git switch %s`, mainBranch),
			fmt.Sprintf(`sed -i 's/b/bbbbb/' %s # ファイル中のbをbbbbbに置き換え`, filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "update b in %s"`, mainBranch),
			fmt.Sprintf("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "PRをマージ",
		Command: fmt.Sprintf("gh pr merge %s --merge --delete-branch", prBranch),
	})

	filename = "pull-req-same-line-conflict.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			fmt.Sprintf(`git switch %s`, mainBranch),
			genAbcFile(filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "create %s"`, filename),
			fmt.Sprintf(`git push origin %s`, mainBranch),
		},
	})

	prBranch = "pr-update-2"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			fmt.Sprintf(`git switch -c %s`, prBranch),
			fmt.Sprintf(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "update a to aaaaa in %s"`, prBranch),
			fmt.Sprintf(`git push --set-upstream origin %s`, prBranch),
			fmt.Sprintf(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: fmt.Sprintf("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			fmt.Sprintf(`git switch %s`, mainBranch),
			fmt.Sprintf(`sed -i 's/a/aaa/' %s # ファイル中のaをaaaに置き換え`, filename),
			"git add --all",
			fmt.Sprintf(`git commit -m "update a to aaa in %s"`, mainBranch),
			fmt.Sprintf("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "PRをマージ",
		Command: fmt.Sprintf("gh pr merge %s --merge --delete-branch", prBranch),
	})
	// bytes, err := json.Marshal(commands)

	RunCommands(commands)
	//WriteShellScript(os.Stdout, commands)
}
