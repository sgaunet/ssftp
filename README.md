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

### go test

The packages can be tested with go test.

### functionnal tests

The program (linux and windows) can be tested with [venom](https://github.com/ovh/venom).

#### Linux

You need :

* vagrant
* virtualbox
* venom

```
cd tests/linux
vagrant up
venom run
```

The testsuite use directly the source (with go run ...)

#### Windows

You need :

* vagrant
* virtualbox
* venom for windows
* md5deep64.exe in the system path

```
cd tests/windows
vagrant up
venom run
```

This testsuite use the binary compiled (task build-windows)