#!/bin/sh

cat << EOF > pull-req-squash-merge.txt
a

b

c
EOF
git add --all
git commit -m "create pull-req-squash-merge.txt"
git push origin main


git switch -c pr-squash-merge-1
sed -i 's/a/aaaaa/' pull-req-squash-merge.txt # ファイル中のaをaaaaaに置き換え
git add --all
git commit -m "update a in pr-squash-merge-1"

# GitHubにPull Requestを作成
git push --set-upstream origin pr-squash-merge-1
gh pr create --title pr-squash-merge-1 --body "" --base main --head pr-squash-merge-1