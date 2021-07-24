#!/usr/bin/env bash


cwd=$(dirname $0)
cd $cwd
tstDir=$(pwd)

projectWorkdir=$(dirname $tstDir)

vagrant destroy -f
vagrant up
echo "Create /tmp/toto (100 MB with random data)"
dd if=/dev/urandom of=/tmp/toto  bs=1024 count=102400


echo "# Test : upload a file"
cd "$projectWorkdir/src"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key /tmp/toto vagrant@10.0.50.10:/tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

echo "# Test 2: upload a file"
ssh -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key \
     vagrant@10.0.50.10 'rm /tmp/toto' 
echo $?
ssh -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key \
     vagrant@10.0.50.10 'mkdir /tmp/toto'
echo $?
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key /tmp/toto vagrant@10.0.50.10:/tmp/toto
rc=$?

rm /tmp/toto

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  vagrant@10.0.50.10:/tmp/toto/toto /tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  vagrant@10.0.50.10:/tmp/toto/toto /tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi


echo "# Test recursive upload"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  ../src vagrant@10.0.50.10:/tmp
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi
