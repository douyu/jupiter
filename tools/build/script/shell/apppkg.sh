#!/bin/bash

# 获取应用相对$GOPATH/src的package path
basePath=$(dirname $(dirname $(dirname $(dirname $(readlink -f $0)))))

echo ${basePath//$GOPATH\/src\//}
