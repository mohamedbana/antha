# antha
Antha v0.0.2

=======
## Getting Started
v 0.0.2 release

Some videos coming soon... 

## Detailed Installation Instructions

### OSX

If you have [Homebrew](http://brew.sh/), you can follow these steps to setup
a working antha development environment:
```sh
# Install go
brew install go
# You probably want to add these lines to your .profile too
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install some external dependencies
brew install glpk sqlite3

# Install antha
go get github.com/antha-lang/antha/cmd/...
```

### Linux

Depending on your Linux distribution, you may not have the most recent version
of go available from your distribution's package repository. We recommend you
[download](https://golang.org/) go directly. For Debian-based distributions like
Ubuntu, the installation instructions are as follows:
```sh
curl -O https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install antha external dependencies
apt-get install -y libglpk-dev libsqlite3-dev

# Now, we are ready to get antha
go get github.com/antha-lang/cmd/...
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
cd $GOPATH/src/github.com/antha-lang/antha/examples/workflows/constructassembly
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
