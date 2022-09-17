package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	mainBranch := "developer"
	filename := "pull-req-update-commit.txt"
	commands = append(commands, Command{Command: fmt.Sprintf(`git switch %s`, mainBranch)})
	commands = append(commands, Command{Command: genAbcFile(filename)})
	commands = append(commands, Command{Command: "git add --all"})
	commands = append(commands, Command{Command: fmt.Sprintf(`git commit -m "%s"`, filename)})
	commands = append(commands, Command{Command: fmt.Sprintf(`git push origin %s"`, mainBranch)})
	bytes, err := json.Marshal(commands)
	if err != nil {
		log.Fatalf("json.Marshal error")
	}
	os.Stdout.Write(bytes)
	fmt.Printf("%v\n", commands)
}
