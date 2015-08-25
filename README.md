# antha
[![GoDoc](http://godoc.org/github.com/antha-lang/antha?status.svg)](http://godoc.org/github.com/antha-lang/antha)
[![Build Status](https://travis-ci.org/antha-lang/antha.svg?branch=master)](https://travis-ci.org/antha-lang/antha)

Antha v0.0.2

## Installation Instructions

### OSX

First step is to install go. Follow the instructions at the
[Golang](http://golang.org/doc/install) site.

After you install go, if you don't have [Homebrew](http://brew.sh/), please
install it. Then, follow these steps to setup a working antha development
environment:
```sh
# Setup environment variables
cat<<EOF>>$HOME/.bash_profile
export GOPATH=$HOME/go
export PATH=\$PATH:$HOME/go/bin
EOF

# Reload your profile
. $HOME/.bash_profile

# Install the xcode developer tools
xcode-select --install

# Install some external dependencies
brew update
brew install homebrew/science/glpk sqlite3

# Install antha
go get github.com/antha-lang/antha/cmd/...
```

### Linux

Depending on your Linux distribution, you may not have the most recent version
of go available from your distribution's package repository. We recommend you
[download](https://golang.org/) go directly. 

For Debian-based distributions like Ubuntu on x86_64 machines, the installation
instructions follow.  If you do not use a Debian based system or if you are not
using an x86_64 machine, you will have to modify these instructions by
replacing the go binary with one that corresponds to your platform and
replacing ``apt-get`` with your package manager.
```sh
# Install go
curl -O https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

# Setup environment variables
cat<<EOF>>$HOME/.bash_profile
export GOPATH=$HOME/go
export PATH=\$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

# Reload your profile
. $HOME/.bash_profile

# Install antha external dependencies
sudo apt-get install -y libglpk-dev libsqlite3-dev git

# Now, we are ready to get antha
go get github.com/antha-lang/antha/cmd/...
```

### Windows

Installing antha on Windows is significantly more involved than for OSX or
Linux. The basic steps are:

  - Setup a go development environment:
    - Install the source code manager [git](https://git-scm.com/download/win)
    - Install [go](https://golang.org/dl/)
    - Install the compiler [mingw](http://sourceforge.net/projects/mingw/files/Installer/mingw-get-setup.exe/download).
      Depending on whether you installed the 386 (32-bit) or amd64 (64-bit) version
      of go, you need to install the corresponding version of mingw.
  - Download antha external dependencies
    - Install [glpk](http://sourceforge.net/projects/winglpk/) development library and make sure that
      mingw can find it.

If this procedure sounds daunting, you can try using some scripts we developed
to automate the installation procedure on Windows.
[Download](scripts/windows/windows-install.zip), unzip them and run
``install.bat``. This will try to automatically apply the Windows installation
procedure with the default options. Caveat emptor.

## Checking Your Installation

After following the installation instructions for your machine. You can check
if Antha is working properly by running a test protocol
```sh
cd $GOPATH/src/github.com/antha-lang/antha/antha/examples/workflows/constructassembly
antharun --workflow workflow.json --parameters parameters.yml
```

## Making and Running Antha Components

The easiest way to start developing your own antha components is to place them
in the ``antha/component/an`` directory and follow the structure of the
existing components there. Afterwards, you can compile and use your components
with the following commands:
```sh
cd $GOPATH/src/github.com/antha-lang/antha
make clean && make
go get github.com/antha-lang/antha/cmd/...
antharun --workflow myworkflowdefinition.json --parameters myparameters.yml
```

## Demo 

[![asciicast](https://asciinema.org/a/12zsgt153sffmfnu2ym7vq9d2.png)](https://asciinema.org/a/12zsgt153sffmfnu2ym7vq9d2)

