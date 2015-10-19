# docker-graphtool

Command line to manage Docker image layers.

## Install

Requires:
 - godep
 - docker

```
go get -d github.com/rsampaio/docker-graphtool
cd $GOPATH/src/github.com/rsampaio/docker-graphtool
make
```

## Running

```
Usage:
  dg mount [--options=<mount_options>] [<image>] [<dest>]
  dg umount [--force] <temp_image>

```
