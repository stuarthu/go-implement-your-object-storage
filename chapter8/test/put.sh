#!/bin/bash

export RABBITMQ_SERVER=amqp://test:test@10.29.102.173:5672
export ES_SERVER=10.29.102.173:9200

echo -n "this is object test8 version 1" | openssl dgst -sha256 -binary | base64
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 1" -H "Digest: SHA-256=2IJQkIth94IVsnPQMrsNxz1oqfrsPo0E2ZmZfJLDZnE="

echo -n "this is object test8 version 2-6" | openssl dgst -sha256 -binary | base64
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl 10.29.2.1:12345/objects/test8 -XPUT -d"this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="

curl 10.29.2.1:12345/versions/test8
curl 10.29.2.1:12345/objects/test8
ls -l /tmp/?/objects

go run ../deleteOldMetadata/deleteOldMetadata.go
curl 10.29.2.1:12345/versions/test8

STORAGE_ROOT=/tmp/1 LISTEN_ADDRESS=10.29.1.1:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
STORAGE_ROOT=/tmp/2 LISTEN_ADDRESS=10.29.1.2:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
STORAGE_ROOT=/tmp/3 LISTEN_ADDRESS=10.29.1.3:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
STORAGE_ROOT=/tmp/4 LISTEN_ADDRESS=10.29.1.4:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
STORAGE_ROOT=/tmp/5 LISTEN_ADDRESS=10.29.1.5:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
STORAGE_ROOT=/tmp/6 LISTEN_ADDRESS=10.29.1.6:12345 go run ../deleteOrphanObject/deleteOrphanObject.go
ls -l /tmp/?/objects
ls -l /tmp/?/garbage

rm /tmp/1/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.*
echo some_garbage > /tmp/2/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.*
ls -l /tmp/?/objects

STORAGE_ROOT=/tmp/2 go run ../objectScanner/objectScanner.go
ls -l /tmp/?/objects
