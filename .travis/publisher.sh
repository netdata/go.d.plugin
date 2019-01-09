#!/bin/bash
# SPDX-License-Identifier: MIT
# Copyright (C) 2018 Pawel Krupa (@paulfantom) - All Rights Reserved
# Permission to copy and modify is granted under the MIT license
#
# Requirements:
#   - GITHUB_TOKEN variable set with GitHub token. Access level: repo.public_repo

set -e

if [ ! -f .gitignore ]; then
	echo "Run as ./travis/$(basename "$0") from top level directory of git repository"
	exit 1
fi

if [ -z ${TRAVIS_TAG+x} ]; then
    exit 1
fi

echo "---- UPLOAD ARTIFACTS TO GITHUB -----"
# Download hub
HUB_VERSION=${HUB_VERSION:-"2.7.0"}
wget "https://github.com/github/hub/releases/download/v${HUB_VERSION}/hub-linux-amd64-${HUB_VERSION}.tgz" -O "/tmp/hub-linux-amd64-${HUB_VERSION}.tgz"
tar -C /tmp -xvf "/tmp/hub-linux-amd64-${HUB_VERSION}.tgz"
export PATH=$PATH:"/tmp/hub-linux-amd64-${HUB_VERSION}/bin"

for i in bin/*; do
	hub release edit -a "$i" -m "${TRAVIS_TAG}" "${TRAVIS_TAG}"
done
