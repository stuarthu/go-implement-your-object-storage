#!/bin/bash

curl -v 10.29.2.1:12345/objects/test7 -XPUT --data-binary @/tmp/file -H "Digest: SHA-256=5bhEzFf1cJTqRYXiNfNseMHNIiJiu4nVPJTctNaz5V0="

curl -v 10.29.2.1:12345/objects/test7 -o /tmp/output

diff /tmp/output /tmp/file

ls -ltr /tmp/?/objects
