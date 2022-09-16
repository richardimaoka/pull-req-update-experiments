#!/bin/sh

cat << EOF > pull-req-update-commit.txt
a

b

c
EOF
git add --all
git commit -m "create pull-req-update-commit.txt"
git push origin main