---
title: Installation
---

You can install Tmuxctl with one of the methods bellow.

The recommended method is to run the following command, which installs `tmuxctl` in ./bin/tmuxctl by default.
```
$ curl -sf https://raw.githubusercontent.com/alexandrebodin/tmuxctl/master/install.sh | sh
```

To install `tmuxctl` in a specific folder run:
```
$ curl -sf https://raw.githubusercontent.com/alexandrebodin/tmuxctl/master/install.sh | BINDIR=/usr/local/bin sh
```

You can install `tmuxctl` using go get
```
$ go get github.com/alexandrebodin/tmuxctl
$ cd $GOPATH/src/github.com/alexandrebodin/tmuxctl
$ dep ensure -vendor-only
$ go install
```

Finally you can install it manually from one of the releases from [github](https://github.com/alexandrebodin/tmuxctl/releases)
