Client

```sh
# master session
ssh -MS control.file host
$ nc -l rhost rport > file

# other session
ssh -S control.file -O forward -L [host:]port:rhost:rport host
$ nc -q 1 host port < file

# cleanup on client
ssh -S control.file -O cancel host
```
