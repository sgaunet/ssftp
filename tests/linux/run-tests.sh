#!/usr/bin/env bash

# for keytype in dsa ecdsa ed25519 rsa
for keytype in ecdsa ed25519 rsa
do
    echo "**********************************"
    echo "keytype=$keytype"
    export VENOM_VAR_keytype=${keytype}
    venom run --stop-on-failure testsuite-docker.yml 
    rc=$?
    if [ "$rc" != "0" ]
    then
        echo "TS Failed, exit 1"
        exit 1
    fi
    echo ""
done