#!/bin/bash
set -e
SRC_FILE=${1}
PACKAGE=${2}
TYPE=${3}
DES=${4}
#uppcase the first char
PREFIX="$(tr '[:lower:]' '[:upper:]' <<< ${TYPE:0:1})${TYPE:1}"
DES_FILE=$(echo ${TYPE}| tr '[:upper:]' '[:lower:]')_${DES}.go
sed 's/PACKAGE_NAME/'"${PACKAGE}"'/g' ${SRC_FILE} | \
    sed 's/GENERIC_TYPE/'"${TYPE}"'/g' | \
    sed 's/GENERIC_NAME/'"${PREFIX}"'/g' > ${DES_FILE}


#gen.h 有很多工具类，已经有现成的写好的工具类
#Genny –  https://github.com/cheekybits/genny
#Generic – https://github.com/taylorchu/generic
#GenGen – https://github.com/joeshaw/gengen
#Gen – https://github.com/clipperhouse/gen