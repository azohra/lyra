#!/bin/bash

# WORK IN PROGRESS
# 1. Identify platform and architecture
# Install, go, brew, dep, and/or any other dependencies for the said architecture
# run make test, build

set -e

LYRA_RELEASE="https://github.com/golang/dep/releases"
GOLANG_RELEASE=""

downloadJSON() {
    url="$2"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        response=$(curl -s -L -w 'HTTPSTATUS:%{http_code}' -H 'Accept: application/json' "$url")
        body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
        code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    elif test -x "$(command -v wget)"; then
        temp=$(mktemp)
        body=$(wget -q --header='Accept: application/json' -O - --server-response "$url" 2> "$temp")
        code=$(awk '/^  HTTP/{print $2}' < "$temp" | tail -1)
        rm "$temp"
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi
    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi

    eval "$1='$body'"
}

initArch() {
    ARCH=$(uname -m)
    if [ -n "$DEP_ARCH" ]; then
        echo "Using DEP_ARCH"
        ARCH="$DEP_ARCH"
    fi
    case $ARCH in
        amd64) ARCH="amd64";;
        x86_64) ARCH="amd64";;
        i386) ARCH="386";;
        *) echo "Architecture ${ARCH} is not supported by this installation script"; exit 1;;
    esac
    echo "ARCH = $ARCH"
}


setup_mac(){
    which brew > /dev/null 
    if [ $? -ne 0 ] ; then
        echo "We need to install brew!"
        ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" </dev/null
    fi

    which go > /dev/null
    if [ $? -ne 0 ] ; then
        brew install go
    fi

    which dep > /dev/null
    if [ $? -ne 0 ] ; then
        brew install dep
        brew update dep
    fi

}

# From https://github.com/golang/dep/blob/master/install.sh
initOS() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    if [ -n "$DEP_OS" ]; then
        echo "Using DEP_OS"
        OS="$DEP_OS"
    fi
    case "$OS" in
        darwin) OS='darwin';;
        linux) OS='linux';;
        freebsd) OS='freebsd';;
        mingw*) OS='windows';;
        msys*) OS='windows';;
        *) echo "OS ${OS} is not supported by this installation script"; exit 1;;
    esac
    echo "OS = $OS"
}

# identify platform based on uname output
initArch
initOS


cat << "EOF"
      :::     :::   ::: :::::::::      :::  
     :+:     :+:   :+: :+:    :+:   :+: :+: 
    +:+      +:+ +:+  +:+    +:+  +:+   +:+ 
   +#+       +#++:   +#++:++#:  +#++:++#++: 
  +#+        +#+    +#+    +#+ +#+     +#+  
 #+#        #+#    #+#    #+# #+#     #+#   
########## ###    ###    ### ###     ###    
EOF                                                                                                                                                                                     


make install