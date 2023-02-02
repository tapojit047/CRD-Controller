#!/bin/bash

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/tapojit047/CRD-Controller/pkg/client \
  github.com/tapojit047/CRD-Controller/pkg/apis \
  fullmetal.com:v1 \
  --output-base "${GOPATH}/src" \
  --go-header-file hack/boilerplate.go.txt