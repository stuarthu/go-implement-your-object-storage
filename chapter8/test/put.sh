#!/bin/bash

curl 10.29.1.1:12345/objects/
echo

echo -n "this is object test8 version 1" | openssl dgst -sha256 -binary | base64
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 1" -H "Digest: SHA-256=2IJQkIth94IVsnPQMrsNxz1oqfrsPo0E2ZmZfJLDZnE="

echo -n "this is object test8 version 2-6" | openssl dgst -sha256 -binary | base64
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="

curl 10.29.2.1:12345/versions/test8
ls -ltr /tmp/?/objects

ES_SERVER=10.29.102.173:9200 go run ../deleteOldMetadata/deleteOldMetadata.go
curl 10.29.2.1:12345/versions/test8

ES_SERVER=10.29.102.173:9200 RABBITMQ_SERVER=amqp://test:test@10.29.102.173:5672 go run ../deleteOrphanObject/deleteOrphanObject.go
ls -ltr /tmp/?/objects

rm /tmp/1/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.*
echo some_garbage > /tmp/2/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.*
ls -ltr /tmp/?/objects
ES_SERVER=10.29.102.173:9200 RABBITMQ_SERVER=amqp://test:test@10.29.102.173:5672 go run ../objectScanner/objectScanner.go
ls -ltr /tmp/?/objects
ls -ltr /tmp/?/garbage
