# blackhole

A generic upload tool designed for ssh force commands.


# The problem

You want to be able to upload data to a server from a semi trusted host, and then
store the data *somewhere*, you don't know where yet.

The motivating use case is storing build artifacts in s3 and ipfs in a write only way, 
using ssh + force_commands as the auth mechanism.

# The solution:

A layer of indirection :).

The untrusted party uploads the data into the blackhole:

```
ssh $SERVER blackhole < data
```


On the server side ssh is restricted to a single command, the blackhole executable.

The blackhole executable does the follow actions:

read stdin into a random file generated $.

run every file in ~/.blackhole_hooks/*

```
$ hook $data
$ rm $data
```

finally removes the data. If any hook returns non zero, blackhole exits with 1.

Configuration:

If passed a single argument, this is taken to be the hook directory.

## notes

Any data printed by hooks will be relayed back to the sending client over stdin/stder.
Metadata of uploads can be prepended to the data processed by your hooks.