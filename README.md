# ssftp

**This project is for test only**

[If you search a tool to make sftp transfer with static compilation, try rclone.](https://rclone.org/)

sftp client tool to transfer files. 

```
ssftp  [-d debug] [-p port] -i sshkey src dest
  -d string
        Debug level (info,warn,debug) (default "info")
  -i string
        SSH key File
  -o value
        Options (Ex: StrictHostKeyChecking=no) 
  -p string
        Port number (default "22")
  -v    Get version
    src: local file/folder or distant sftp file/dir
    dest: same

Order of parameters matters. (options before src/dest parameters)
```

Program checks the host key fingerprint by default. Add : -o StrictHostKeyChecking=no to disable the check.

Actually works with algorithms :

* ecdsa
* ed25519
* rsa

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
* docker

```
cd tests/linux
vagrant up
venom run testsuite.yml     # to launch tests with the VM
vagrant halt

run-tests.sh                # to launch tests with a sshd in a docker image
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