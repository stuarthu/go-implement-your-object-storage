#!/bin/bash

echo -n "this object will be separate to 4+2 shards" | openssl dgst -sha256 -binary | base64

curl -v 10.29.2.1:12345/objects/test5 -XPUT -d "this object will be separate to 4+2 shards" -H "Digest: SHA-256=MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8="

ls -ltr /tmp/?/objects
echo
curl 10.29.2.1:12345/objects/test5
echo
curl 10.29.2.1:12345/locate/MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8=

rm /tmp/1/objects/MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8=.*
echo some_garbage > /tmp/2/objects/MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8=.*
ls -ltr /tmp/?/objects
echo
curl 10.29.2.1:12345/objects/test5
echo
ls -ltr /tmp/?/objects
