package main

import "fmt"

func abcText(filename string) string {
	return fmt.Sprintf(`cat << EOF > %s.txt
a

b

c
EOF`, filename)
}

func main() {
	fmt.Println(abcText("pull-req-update-commit"))
}
