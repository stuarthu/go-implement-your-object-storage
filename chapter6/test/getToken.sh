#!/bin/bash

echo -n "this is a resumable object" | openssl dgst -sha256 -binary | base64

echo -n "this is a resumable object" | wc

curl 10.29.2.1:12345/temp/test6 -XPOST -H "Digest: SHA-256=u5TzpCX5ck1rMA4RJYNap3pyqLQhhaKZV1o4V7Rmzyc=" -H "Size: 26"
