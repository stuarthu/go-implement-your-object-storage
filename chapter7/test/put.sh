#!/bin/bash

curl -v 10.29.2.1:12345/objects/test7 -XPUT --data-binary @/tmp/file -H "Digest: SHA-256=IEkqTQ2E+L6xdn9mFiKfhdRMKCe2S9v7Jg7hL6EQng4="

curl -v 10.29.2.1:12345/objects/test7 -o /tmp/output

diff -s /tmp/output /tmp/file

ls -ltr /tmp/?/objects

curl -v 10.29.2.1:12345/objects/test7 -H "Accept-Encoding: gzip" -o /tmp/output2.gz

gunzip /tmp/output2.gz

diff -s /tmp/output2 /tmp/file
