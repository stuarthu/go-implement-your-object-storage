#!/bin/bash

curl -v 10.29.2.1:12345/objects/test6 -XPOST -H "Digest: SHA-256=$1" -H "Size: 100000"
