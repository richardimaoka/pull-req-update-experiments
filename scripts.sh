#!/bin/sh

# # 準備: GitHub レポジトリの作成
# mkdir pull-req-update-experiments
# cd pull-req-update-experiments
# git init

# # GitHub repository create
# gh repo create pull-req-update-experiments --public --source=. --remote=origin

# 準備: GitHub テキストファイルの作成
git switch developer
cat << EOF > pull-req-no-conflict.txt
a

b

c
EOF
git add --all
git commit -m "create pull-req-no-conflict.txt"
git push origin developer

Pull Request作成
git switch -c pr-update-1
sed -i 's/a/aaaaa/' pull-req-no-conflict.txt # ファイル中のaをaaaaaに置き換え
git add --all
git commit -m "update a in pr-update-1"
git push --set-upstream origin "pr-update-1"
gh pr create --title pr-update-1 --body "" --base developer --head pr-update-1

developer ブランチにcommit
git switch developer
sed -i 's/b/bbbbb/' pull-req-no-conflict.txt # ファイル中のbをbbbbbに置き換え
git add --all
git commit -m "update b in developer"
git push origin developer

