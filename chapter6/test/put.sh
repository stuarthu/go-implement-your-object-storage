#!/bin/bash

curl -I 10.29.2.1:12345/$1

curl -v -XPUT --data-binary @/tmp/first 10.29.2.1:12345/$1

curl -I 10.29.2.1:12345/$1

curl -v -XPUT --data-binary @/tmp/second -H "range: bytes=32000-" 10.29.2.1:12345/$1

curl -I 10.29.2.1:12345/$1

curl -v -XPUT --data-binary @/tmp/last -H "range: bytes=96000-" 10.29.2.1:12345/$1

curl 10.29.2.1:12345/objects/test6 > /tmp/output

diff -s /tmp/output /tmp/file

curl 10.29.2.1:12345/objects/test6 -H "range: bytes=32000-" > /tmp/output2

diff -s /tmp/output2 /tmp/second
