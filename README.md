# Kubeconnect [![Build Status](https://travis-ci.org/marijnkoesen/kubeconnect.svg?branch=master)](https://travis-ci.org/marijnkoesen/kubeconnect)

[![iMIT License](https://img.shields.io/badge/license-mit-blue.svg)](https://github.com/marijnkoesen/kubeconnect/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/marijnkoesen/kubeconnect)](https://goreportcard.com/report/github.com/marijnkoesen/kubeconnect)

Kubeconnect Is a tool with which you can easily connect to any pod running in any of your kubernetes clusters.

Don't you know the exact namespace, the pod name, or the name of that one container in the pod? This is the command for you.

<p align="center"><img src="/doc/demo.gif?raw=true"/></p>

## Installation

### OSX 

Installing on OSX can be done using brew:

```
$ brew tap marijnkoesen/kubeconnect
$ brew install kubeconnect
```

### Others

Download one of the releases from the [release page](https://github.com/MarijnKoesen/kubeconnect/releases).

Extract the archive and run the `kubeconnect` command.


## Building from source

Building from source is as simple as a `go build`

```
$ git clone https://github.com/MarijnKoesen/kubeconnect.git
$ cd kubeconnect
$ go build
$ ./kubeconnect
```
