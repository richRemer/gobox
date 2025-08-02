Steps to build and reproduce

```sh
# clone repo and change to repo dir
git clone -b segviol https://github.com/richRemer/gobox.git
cd gobox

# build
make

# run app server listening on port 2222 and writing DB to $(pwd)/test.db
./gobox -p 2222 -D test.db

# from another terminal, connect over SSH as 'admin' user using provided test key
chmod 600 admin_key
ssh -i admin_key -p 2222 admin@localhost
# wait for splash screen to go away (this ensures write to DB has ocurred)
# hit 'q' to exit

# then back in other terminal, CTRL+C to exit
# you should see the error at this point
```
