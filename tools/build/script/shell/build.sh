#!/bin/bash

VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]];then
    V="-x"
    set -x
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

OUT=${1:?"output path"}
VERSION_PACKAGE=${2:?"version go package"}
BUILDPATH=${3:?"path to build"}

set -e

GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
GOPKG="$GOPATH/pkg"
BUILDINFO=${BUILDINFO:-""}
STATIC=${STATIC:-1}
LDFLAGS="-extldflags -static"
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
GCFLAGS=${GCFLAGS:-}

SUBBUILDPATH=$(dirname ${BUILDPATH})

export CGO_ENABLED=0

echo -e "\n"
echo "GOOS:"${GOOS}
echo "GOARCH:"${GOARCH}
echo "GOBINARY:"${GOBINARY}
echo "GOPKG:"${GOPKG}
echo "BUILDINFO:"${BUILDINFO}
echo "STATIC:"${STATIC}
echo "LDFLAGS:"${LDFLAGS}
echo "GOBUILDFLAGS:"${GOBUILDFLAGS}
echo "GCFLAGS:"${GCFLAGS}
echo "BUILDPATH:"${BUILDPATH}
echo "SUBBUILDPATH:"${SUBBUILDPATH}
echo -e "\n"


if [[ "${STATIC}" !=  "1" ]];then
    LDFLAGS=""
fi

# gather buildinfo if not already provided
# For a release build BUILDINFO should be produced
# at the beginning of the build and used throughout
if [[ -z ${BUILDINFO} ]];then
    BUILDINFO=$(mktemp)
    ${ROOT}/shell/version.sh > ${BUILDINFO}
fi

# BUILD LD_VERSIONFLAGS
LD_VERSIONFLAGS=""
while read line; do
    read SYMBOL VALUE < <(echo $line)
    LD_VERSIONFLAGS=${LD_VERSIONFLAGS}" -X ${VERSION_PACKAGE}.${SYMBOL}='${VALUE}'"
done < "${BUILDINFO}"

echo $LD_VERSIONFLAGS

# 读取 BUILDPATH 目录下是否有文件夹，自动进行文件建 main 文件构建
for dir in $(ls ${SUBBUILDPATH})
do
    if [[  ${dir} == "main.go" ]]
        then
            echo -e "\n"
            echo "dir:"$dir
            echo "OUT:"${OUT}
            echo "BUILDPATH:"${BUILDPATH}
            echo -e "\n"
            time GOOS=${GOOS} GOARCH=${GOARCH} ${GOBINARY} build ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o ${SUBBUILDPATH}"/../bin/"${OUT} \
            -pkgdir=${GOPKG}/${GOOS}_${GOARCH} -ldflags "${LDFLAGS} ${LD_VERSIONFLAGS}" "${BUILDPATH}"
        else
            TMPOUT=${OUT:0:(${#OUT})-3}"-"${dir}"-go"
            TMPBUILDPATH=${BUILDPATH:0:(${#BUILDPATH})-7}${dir}"/main.go"
            echo -e "\n"
            echo "dir:"$dir
            echo "TMPOUT:"${TMPOUT}
            echo "TMPBUILDPATH:"${TMPBUILDPATH}
            echo -e "\n"
            time GOOS=${GOOS} GOARCH=${GOARCH} ${GOBINARY} build ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o  ${SUBBUILDPATH}"/../bin/"${TMPOUT} \
            -pkgdir=${GOPKG}/${GOOS}_${GOARCH} -ldflags "${LDFLAGS} ${LD_VERSIONFLAGS}" "${TMPBUILDPATH}"
    fi
done
