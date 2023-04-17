#!/bin/bash

Set-PSDebug -Trace 1
$ErrorActionPreference = "Stop"

try {
    # 设置 -o pipefail 选项，使得管道中的任何一个命令失败都会导致整个管道失败
    $null | Out-Default
}
catch {
    Write-Error $_.Exception.Message
}

# 设置 -eux
Set-PSDebug -Trace 2

mkdir -p plugin-binaries
gotip env -w GOOS=windows
gotip build -o ./plugin-binaries/tt ./sample-plugins/tt/
gotip build -o ./plugin-binaries/narcissist  ./sample-plugins/narcissist/
gotip build -o ./htmlize ./main.go
gotip env -w GOOS=linux