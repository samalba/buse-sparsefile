#!/bin/sh

[ $VAGRANT = "true" ] || {
    echo "This script must be called by a Vagrantfile"
    exit 0
}

[ -d /usr/local/go ] && {
    echo "Found earlier version of Golang, cleanup..."
    rm -rf /usr/local/go
}

# Install golang
(
    version=1.9.linux-amd64
    cd /tmp
    echo "Dowloading Golang ${version}..."
    wget -q https://storage.googleapis.com/golang/go${version}.tar.gz
    echo "Installing..."
    tar -C /usr/local -xzf go${version}.tar.gz
    echo "Cleanup."
    rm go${version}.tar.gz
)
echo export PATH=\$PATH:/usr/local/go/bin >> .profile
echo export GOPATH=/home/vagrant/go >> .profile
echo alias l=\'ls -lh\' >> .bashrc
echo alias la=\'ls -lha\' >> .bashrc
mkdir -p go && chown vagrant: go

# Install git (for golang deps)
apt-get -y update
apt-get -y install git
apt-get clean
