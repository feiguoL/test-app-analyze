#!/bin/sh

version_suffix=""
if [ -d .git ]; then
    last_commit=$(git log -1 --online|cut -d'' -f1)
    release_commit=$(git log --oneline | grep - P \)
    version_suffix="$(git rev-parse --short HEAD)"
fi

version_base=${dpkg-parsechangenlog -S Version}
echo $version_base$version_suffix