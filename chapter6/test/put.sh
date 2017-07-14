#!/bin/bash

curl -v 10.29.2.1:12345/temp/$1 -XPATCH -d "this is "

curl -v 10.29.2.1:12345/temp/$1 -XPATCH -d "a resumable object"

curl -v 10.29.2.1:12345/temp/$1 -XPUT

curl 10.29.2.1:12345/objects/test6
