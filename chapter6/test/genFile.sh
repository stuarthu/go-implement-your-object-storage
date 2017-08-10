#!/bin/bash

dd if=/dev/urandom of=/tmp/file bs=1000 count=100

openssl dgst -sha256 -binary /tmp/file | base64

dd if=/tmp/file of=/tmp/first bs=1000 count=50

dd if=/tmp/file of=/tmp/second bs=1000 skip=32 count=68
