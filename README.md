SSH server written in Go using [Wish](https://github.com/charmbracelet/wish)
library.  Sandbox for playing with TUIs in Go.

Getting Started
===============

### Install
```sh
git clone https://github.com/richRemer/gobox.git gobox-src
cd ./gobox-src
make
make install

# setup working directory for server
mkdir /var/local/gobox
cd /var/local/gobox
ssh-keygen -qN "" -ted25519 -f host_key

# use current user's SSH key to setup an admin user
ADMIN_USER=$(id -un)
cp ~$ADMIN_USER/.ssh/id_rsa.pub ${ADMIN_USER}_key
chmod 600 ${ADMIN_USER}_key
```

### Run
```sh
cd /var/local/gobox
gobox
```

Command-Line Options
--------------------
The following command-line options are recognized.

#### -d *WORKING_DIR*
Path to working directory.  When reading path options, **gobox** resolves them
relative to the working directory.  Defaults to current directory.

#### -D *DB_FILE*
Path to runtime database.  This database stores additional user keys and any
other state necessary at runtime.  If the file does not exist, it will be
created.  Defaults to an empty in-memory database.

#### -k *KEYS_DIR*
Path to directory where user keys are stored.  The files in this directory
must be named '*USERNAME*_key'.  Defaults to working directory.

#### -h *HOSTKEY_FILE*
Path to file which stores the server host key.  This host key is used to
verify the host when a client connects so the client knows they are not being
snooped on.  Defaults to "host_key".

Currently, a host key will be generated when the server starts if one does not
already exist.  This is expected to change in the future.

#### -p *INT*
Specify the port to listen on.  Defaults to 22.
