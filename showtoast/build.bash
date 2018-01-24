#!/bin/bash

function main() {
    checkcmd apktool
    checkcmd gomobile

    if [ -e tmp ]; then
        rm -rf tmp
    fi
    mkdir tmp

    GOPATH=`pwd`:$GOPATH gomobile build showtoast
    mv showtoast.apk tmp
    apktool d tmp/showtoast.apk -o tmp/showtoast 2> /dev/null 1> /dev/null

    cp tmp/showtoast/AndroidManifest.xml Android/app/src/main/
    mkdir -p Android/app/src/main/java/org/golang/app
    #cp `go list -f {{.Dir}} golang.org/x/mobile/app`/GoNativeActivity.java Android/app/src/main/java/org/golang/app

    mkdir -p Android/app/src/main/jniLibs
    cp -r tmp/showtoast/lib/* Android/app/src/main/jniLibs

    if [ -e tmp ]; then
        rm -rf tmp
    fi
}

function checkcmd() {
    if type $1 2>/dev/null 1>/dev/null
    then
        :
    else
        echo $1" is not installed"
        exit 1
    fi
}

main
