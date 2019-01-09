#!/bin/sh
for TARGET in "${@}"; do
  find "${TARGET}" -name '*.go' -exec gofmt -s -w {} \+
done
git diff --exit-code
