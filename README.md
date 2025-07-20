Simple proof of concept SSH server written in Go using the
[Wish](https://github.com/charmbracelet/wish) library.

Fetch
=====
```sh
git clone https://github.com/richRemer/gobox.git gobox-src
```

Build
=====
```sh
cd ./gobox-src
go build gobox.go
```

Install
=======
```sh
GOBOX_ROOT=/gobox/root
cp gobox $GOBOX_ROOT/
```

Setup
=====
First, generate a host key for gobox.

```sh
cd $GOBOX_ROOT
ssh-keygen -qN "" -ted25519 -f host_key
```

Next, configure authorized keys.

```sh
USER=$(id -un)
cd $GOBOX_ROOT
cp ~/.ssh/id_rsa.pub ${USER}_keys
chmod 600 ${USER}_keys
```

Run
==========
```sh
cd $GOBOX_ROOT
./gobox
```
