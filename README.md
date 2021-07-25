# ssftp

sftp client tool to transfer file. 

```
ssftp  [-d debug] -i sshkey src dest
    -i : ssh key
    -d : debug mode
    -v : print version and exit
    src: local file/folder or distant sftp file/dir
    dest: same

Order of parameters matters. (options before src/dest parameters)
```

**Be carefull, the program does not check the hostkey so it's not a secure program for now.**

# Build

## Compile

```
cd src
go build . 
```

## Tests

The package pathh can be tested with go test.

The entire program can be tested with a VM created by vagrant

```
cd tst
vagrant up
./tests.sh
```

