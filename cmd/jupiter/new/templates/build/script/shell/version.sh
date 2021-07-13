#!/bin/bash

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

BUILD_GIT_REVISION=$(git rev-parse HEAD 2> /dev/null)
if [[  $? == 0 ]]; then
    git diff-index --quiet HEAD
    if [[  $? != 0 ]]; then
        BUILD_GIT_REVISION=${BUILD_GIT_REVISION}"-dirty"
    fi
else
    BUILD_GIT_REVISION=unknown
fi

# Check for local changes
git diff-index --quiet HEAD --
if [[ $? == 0 ]];
then
  tree_status="Clean"
else
  tree_status="Modified"
fi

# XXX This needs to be updated to accomodate tags added after building, rather than prior to builds
RELEASE_TAG=$(git describe --match '[0-9]*\.[0-9]*\.[0-9]*' --exact-match 2> /dev/null || echo "")

# security wanted VERSION='unknown'
VERSION="${BUILD_GIT_REVISION}"
if [[ -n "${RELEASE_TAG}" ]]; then
  VERSION="${RELEASE_TAG}"
fi

# used by pkg/version
echo buildVersion       "${VERSION}"
echo buildGitRevision   "${BUILD_GIT_REVISION}"
echo buildUser          "$(whoami)"
echo buildHost          "$(hostname -f)"
echo buildStatus        "${tree_status}"
# echo buildTime          "$(date '+%Y-%m-%d %T')"
echo buildTime          "$(date '+%Y-%m-%d--%T')"
