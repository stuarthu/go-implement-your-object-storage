#!/bin/bash

curl -v 10.29.2.1:12345/objects/test7 -XPUT --data-binary @file -H "Digest: SHA-256=$1"

curl -v 10.29.2.1:12345/objects/test7 -o output

diff output file
