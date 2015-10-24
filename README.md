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
  dg bundle <image> <file.tar>

```

Example usage:

```shell
$ docker pull centos:7
$ dg mount centos:7 /tmp/centos
```

If you are in a systemd distro, like Arch:

```shell
$ systemd-nspawn -D /tmp/centos /bin/bash
# yum -y install systemd
# ^D
$ systemd-nspawn --boot -D /tmp/centos
...systemd init messages...

```

You can also export a [bundle](https://github.com/opencontainers/specs/blob/master/bundle.md) from a docker image:

```shell
$ docker pull ghost
$ dg bundle ghost ghost.tar
```

Run with [runc](https://github.com/opencontainers/runc):
```
$ mkdir -p ghost
$ tar -C ghost xvf ghost.tar
$ cd ghost
~/ghost $ runc start
/ #
```
