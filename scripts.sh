#!/bin/sh

mkdir pull-req-update-experiments
cd pull-req-update-experiments
git init

gh repo create pull-req-update-experiments --public --source=. --remote=origin

git switch developer
cat << EOF > pull-req-no-conflict.txt
a

b

c
EOF
git add --all
git commit -m "pull-req-no-conflict.txt"
git push origin developer

git switch pr-update-1
sed -i 's/a/aaaaa/' pull-req-no-conflict.txt # ファイル中のaをaaaaaに置き換え
git add --all
git commit -m "pull-req-no-conflict.txt"
git push --set-upstream origin "pr-update-1"
gh pr create --title pr-update-1 --body "" --base developer --head pr-update-1

git switch main
sed -i 's/b/bbbbb/' pull-req-merge-commit.txt # ファイル中のbをbbbbbに置き換え
git add --all
git commit -m "update b in main"

# GitHubにPull Requestを作成
git push --set-upstream origin pr-merge-commit-2
gh pr create --title pr-merge-commit-2 --body "" --base main --head pr-merge-commit-2