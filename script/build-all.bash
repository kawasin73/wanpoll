#!/usr/bin/env bash
# Copied from https://github.com/golang/dep/blob/master/hack/build-all.bash
# Copyright (c) 2014 The Go Authors. All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are
# met:
#
#    * Redistributions of source code must retain the above copyright
# notice, this list of conditions and the following disclaimer.
#    * Redistributions in binary form must reproduce the above
# copyright notice, this list of conditions and the following disclaimer
# in the documentation and/or other materials provided with the
# distribution.
#    * Neither the name of Google Inc. nor the names of its
# contributors may be used to endorse or promote products derived from
# this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
# "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
# LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
# A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
# OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
# SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
# LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
# DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
# THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
# (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#
# This script will build dep and calculate hash for each
# (DEP_BUILD_PLATFORMS, DEP_BUILD_ARCHS) pair.
# DEP_BUILD_PLATFORMS="linux" DEP_BUILD_ARCHS="amd64" ./hack/build-all.bash
# can be called to build only for linux-amd64

set -e

DEF_ROOT=$(git rev-parse --show-toplevel)

if [[ "$(pwd)" != "${DEF_ROOT}" ]]; then
  echo "you are not in the root of the repo" 1>&2
  echo "please cd to ${DEF_ROOT} before running this script" 1>&2
  exit 1
fi

GO_BUILD_CMD="go build -a"

if [[ -z "${DEF_BUILD_PLATFORMS}" ]]; then
    DEF_BUILD_PLATFORMS="linux darwin"
fi

if [[ -z "${DEF_BUILD_ARCHS}" ]]; then
    DEF_BUILD_ARCHS="amd64 386"
fi

mkdir -p "${DEF_ROOT}/release"

for OS in ${DEF_BUILD_PLATFORMS[@]}; do
  for ARCH in ${DEF_BUILD_ARCHS[@]}; do
    NAME="wanpoll-${OS}-${ARCH}"
    if [[ "${OS}" == "windows" ]]; then
      NAME="${NAME}.exe"
    fi
    echo "Building for ${OS}/${ARCH}"
    GOARCH=${ARCH} GOOS=${OS} ${GO_BUILD_CMD} -o "${DEF_ROOT}/release/${NAME}" .
    shasum -a 256 "${DEF_ROOT}/release/${NAME}" > "${DEF_ROOT}/release/${NAME}".sha256
  done
done
