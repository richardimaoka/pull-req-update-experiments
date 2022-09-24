package main

import (
	"fmt"
	"os"
)

func genAbcFile(filename string) string {
	return fmt.Sprintf(
		`cat << EOF > %s
a

b

c
EOF`, filename)
}

func main() {
	var commands []fmt.Stringer
	f := fmt.Sprintf

	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub レポジトリの作成",
		Commands: []string{
			"mkdir pull-req-update-experiments",
			"cd pull-req-update-experiments",
			"git init",
		},
	})

	commands = append(commands, &SingleCommand{
		Comment: "GitHub repository create",
		Command: "gh repo create pull-req-update-experiments --public --source=. --remote=origin",
	})

	mainBranch := "main"
	if mainBranch != "main" {
		commands = append(commands, &MultiCommands{
			Commands: []string{
				f("git branch -f %s", mainBranch),
				f(`git switch %s`, mainBranch),
			},
		})
	}

	filename := "experiment1.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			genAbcFile(filename),
			"git add --all",
			f(`git commit -m "create %s"`, filename),
			f(`git push origin %s`, mainBranch),
		},
	})

	prBranch := "pr-1"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			f(`git switch -c %s`, prBranch),
			"",
			f(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update a in %s"`, prBranch),
			"",
			f(`git push --set-upstream origin %s`, prBranch),
			f(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			f(`git switch %s`, mainBranch),
			"",
			f(`sed -i 's/b/bbbbb/' %s # ファイル中のbをbbbbbに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update b in %s"`, mainBranch),
			"",
			f("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチでのgit log 確認", prBranch),
		Commands: []string{
			f("git switch %s", prBranch),
			f("git pull origin %s", prBranch),
			"git log --oneline --decorate --graph",
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: "PRをマージ",
		Commands: []string{
			f("gh pr merge %s --merge --delete-branch", prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチでのgit log 確認", mainBranch),
		Commands: []string{
			f("git switch %s", mainBranch),
			f("git pull origin %s", mainBranch),
			"git log --oneline --decorate --graph",
		},
	})

	//########################################################
	// Update with rebase
	//########################################################

	filename = "experiment2.txt"
	commands = append(commands, &MultiCommands{
		Comment: "準備: GitHub テキストファイルの作成",
		Commands: []string{
			genAbcFile(filename),
			"git add --all",
			f(`git commit -m "create %s"`, filename),
			f(`git push origin %s`, mainBranch),
		},
	})

	prBranch = "pr-2"
	commands = append(commands, &MultiCommands{
		Comment: "Pull Request作成",
		Commands: []string{
			f(`git switch -c %s`, prBranch),
			"",
			f(`sed -i 's/a/aaaaa/' %s # ファイル中のaをaaaaaに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update a in %s"`, prBranch),
			"",
			f(`git push --set-upstream origin %s`, prBranch),
			f(`gh pr create --title %s --body "" --base %s --head %s`, prBranch, mainBranch, prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチに直接commit", mainBranch),
		Commands: []string{
			f(`git switch %s`, mainBranch),
			"",
			f(`sed -i 's/b/bbbbb/' %s # ファイル中のbをbbbbbに置き換え`, filename),
			"git add --all",
			f(`git commit -m "update b in %s"`, mainBranch),
			"",
			f("git push origin %s", mainBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチでのgit log 確認", prBranch),
		Commands: []string{
			f("git switch %s", prBranch),
			f("git pull origin %s", prBranch),
			"git log --oneline --decorate --graph",
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: "PRをマージ",
		Commands: []string{
			f("gh pr merge %s --rebase --delete-branch", prBranch),
		},
	})

	commands = append(commands, &MultiCommands{
		Comment: f("%s ブランチでのgit log 確認", mainBranch),
		Commands: []string{
			f("git switch %s", mainBranch),
			f("git pull origin %s", mainBranch),
			"git log --oneline --decorate --graph",
		},
	})

	file, err := os.Create("script.md")
	if err != nil {
		fmt.Println("ERROR: cannot open script.sh")
	} else {
		WriteMarkdown(file, commands)
	}

	RunCommands(commands)

}
