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
ssh-keygen -qN "" -ted25519 -f gobox.key
```

Next, configure authorized keys.

```sh
cd $GOBOX_ROOT
cat ~/.ssh/id_rsa.pub authorized_keys
chmod 600 authorized_keys
```

Run
==========
```sh
cd $GOBOX_ROOT
./gobox
```
