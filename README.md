Steps to build and reproduce

```sh
# build and run app server listening on port 2222 and writing DB to $(pwd)/test.db
make
./gobox -p 2222 -D test.db

# from another terminal, connect over SSH as 'admin' user using provided test key
ssh -i admin_key -p 2222 admin@localhost
# wait for splash screen to go away
# hit 'q' to exit

# then back in other terminal, CTRL+C to exit
# you should see the error at this point
```
