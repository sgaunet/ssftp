#!/usr/bin/env bash


cwd=$(dirname $0)
cd $cwd
tstDir=$(pwd)

projectWorkdir=$(dirname $tstDir)

echo "Create /tmp/toto (100 MB with random data)"
dd if=/dev/urandom of=/tmp/toto  bs=1024 count=102400


cd "$projectWorkdir/src"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key /tmp/toto vagrant@10.0.50.10:/tmp/toto
rc=$?

rm /tmp/toto

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  vagrant@10.0.50.10:/tmp/toto /tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi
