#!/bin/bash

# Setup the environment prior to testing.

# Get OS
case $(uname -s) in
    Darwin)
        OS=darwin
        ;;
    Linux)
        if LSB_RELEASE=$(which lsb_release); then
            OS=$($LSB_RELEASE -s -c)
        else
            # No lsb-release, trya hack or two
            if which dpkg 1>/dev/null; then
                OS=debian
            elif which yum 1>/dev/null || which dnf 1>/dev/null; then
                OS=redhat
            else
                echo "Linux distro not yet supported"
                exit 1
            fi
        fi
        ;;
    *)
        echo "Unsupported OS"
        exit 1
        ;;
esac
echo "Depected OS is '$OS'"

# Now install go
case $OS in
    xenial)
    add-apt-repository ppa:longsleep/golang-backports
    apt-get update
    apt-get install -y golang-go golint
    ;;
esac

# Install vgo https://github.com/golang/go/wiki/vgo
if which go 1>/dev/null; then
    go get -u -v golang.org/x/vgo
else
    echo "go not found, install golang from source?"
fi
