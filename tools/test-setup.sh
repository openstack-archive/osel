#!/bin/bash

echo "Environment check"
echo "GOROOT: $GOROOT"
echo "GOPATH: $GOPATH"

# Setup the environment prior to testing.
export PATH=$PATH:$GOPATH/bin

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

echo | sudo -S /bin/true 2>/dev/null
if [ $? != 0]; then
    echo "Sudo does not work, so packages can not be installed"
    exit 1
fi

# Now install go
case $OS in
    xenial)
    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt-get update
    sudo apt-get install -y golang-go golint
    ;;
esac

# Install vgo https://github.com/golang/go/wiki/vgo
if which go 1>/dev/null; then
    sudo go get -u -v golang.org/x/vgo
    find / -type f -a -name vgo
else
    echo "go not found, install golang from source?"
fi
