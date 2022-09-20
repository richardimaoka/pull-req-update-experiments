#!/bin/sh

```:コピペして実行
git branch -f developer
```
```:コピペして実行
git switch developer
cat << EOF > pull-req-no-conflict.txt
a

b

c
EOF
git add --all
git commit -m "create pull-req-no-conflict.txt"
git push origin developer
```
```:コピペして実行
git switch -c pr-update-1
sed -i 's/a/aaaaa/' pull-req-no-conflict.txt # ファイル中のaをaaaaaに置き換え
git add --all
git commit -m "update a in pr-update-1"
git push --set-upstream origin pr-update-1
gh pr create --title pr-update-1 --body "" --base developer --head pr-update-1
```
```:コピペして実行
git switch developer
sed -i 's/b/bbbbb/' pull-req-no-conflict.txt # ファイル中のbをbbbbbに置き換え
git add --all
git commit -m "update b in developer"
git push origin developer
```
```:コピペして実行
gh pr merge pr-update-1 --merge --delete-branch
```
```:コピペして実行
git switch developer
cat << EOF > pull-req-same-line-conflict.txt
a

b

c
EOF
git add --all
git commit -m "create pull-req-same-line-conflict.txt"
git push origin developer
```
```:コピペして実行
git switch -c pr-update-2
sed -i 's/a/aaaaa/' pull-req-same-line-conflict.txt # ファイル中のaをaaaaaに置き換え
git add --all
git commit -m "update a to aaaaa in pr-update-2"
git push --set-upstream origin pr-update-2
gh pr create --title pr-update-2 --body "" --base developer --head pr-update-2
```
```:コピペして実行
git switch developer
sed -i 's/a/aaa/' pull-req-same-line-conflict.txt # ファイル中のaをaaaに置き換え
git add --all
git commit -m "update a to aaa in developer"
git push origin developer
```
```:コピペして実行
gh pr merge pr-update-2 --merge --delete-branch
```
