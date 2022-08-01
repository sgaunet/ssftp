#!/usr/bin/env bash

# for keytype in dsa ecdsa ecdsa-sk ed25519 ed25519-sk rsa
for keytype in ecdsa ed25519 rsa
do
    echo "**********************************"
    echo "keytype=$keytype"
    export VENOM_VAR_keytype=${keytype}
    venom run testsuite-docker.yml 
    echo ""
done