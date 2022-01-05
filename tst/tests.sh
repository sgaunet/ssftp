#!/usr/bin/env bash


cwd=$(dirname $0)
cd $cwd
tstDir=$(pwd)

projectWorkdir=$(dirname $tstDir)

vagrant destroy -f
vagrant up
echo -e "\nCreate /tmp/toto (100 MB with random data)"
dd if=/dev/urandom of=/tmp/toto  bs=1024 count=102400


echo -e "\n\n# Test : upload a file"
cd "$projectWorkdir/src"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key /tmp/toto vagrant@10.0.50.10:/tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

echo -e "\n\n# Test 2: upload a file to a dir of the same name"
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

echo -e "\n\n# Test 3 : Download the file /tmp/toto/toto to /tmp/toto"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  vagrant@10.0.50.10:/tmp/toto/toto /tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

echo -e "\n\n# Test 4 : Same test"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  vagrant@10.0.50.10:/tmp/toto/toto /tmp/toto
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi


echo -e "\n\n# Test 5 : recursive upload ../src to /tmp"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key  ../src vagrant@10.0.50.10:/tmp
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

echo -e "\n\n# Test 6 : recursive download /tmp/src to /tmp"
go run . -i ${projectWorkdir}/tst/.vagrant/machines/ex-0/virtualbox/private_key   vagrant@10.0.50.10:/tmp/src /tmp/src2
rc=$?

if [ "$rc" != "0" ]
then
    echo "Error when uploading a file"
    exit 1
fi

echo -e "\n\n# Test 7 : Diff between src sent dir and the downloaded"
diff -r ../src /tmp/src2
rc=$?

if [ "$rc" != "0" ]
then
    echo "TEST FAILED"
    exit 1
fi

rm -rf /tmp/src2

