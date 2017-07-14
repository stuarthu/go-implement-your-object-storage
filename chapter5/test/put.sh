#!/bin/bash

echo -n "this object will be separate to 4+2 shards" | openssl dgst -sha256 -binary | base64

curl -v 10.29.2.1:12345/objects/test5 -XPUT -d "this object will be separate to 4+2 shards" -H "Digest: SHA-256=MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8="

ls -ltr /tmp/1/objects /tmp/2/objects /tmp/3/objects /tmp/4/objects /tmp/5/objects /tmp/6/objects

curl -v 10.29.2.1:12345/objects/test5

curl -v 10.29.2.1:12345/locate/MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8=
