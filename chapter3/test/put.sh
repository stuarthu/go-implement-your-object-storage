#!/bin/bash

curl -v 10.29.2.2:12345/objects/test3 -XPUT -d"this is object test3"

echo -n "this is object test3" | openssl dgst -sha256 -binary | base64
curl -v 10.29.2.2:12345/objects/test3 -XPUT -d"this is object test3" -H "Digest: SHA-256=GYqqAdFPt+CScnUDc0/Gcu3kwcWmOADKNYpiZtdbgsM="

curl 10.29.2.1:12345/objects/test3
echo

echo -n "this is object test3 version 2" | openssl dgst -sha256 -binary | base64
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d"this is object test3 version 2" -H "Digest: SHA-256=cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo="

curl 10.29.2.1:12345/objects/test3
echo

curl 10.29.2.1:12345/objects/test3
echo
curl 10.29.2.1:12345/locate/GYqqAdFPt+CScnUDc0%2FGcu3kwcWmOADKNYpiZtdbgs=
echo
curl 10.29.2.1:12345/locate/cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=
echo
curl 10.29.2.1:12345/versions/test3
echo
curl 10.29.2.1:12345/objects/test3?version=1
echo
curl -v 10.29.2.1:12345/objects/test3 -XDELETE

curl -v 10.29.2.1:12345/objects/test3
echo

curl 10.29.2.1:12345/versions/test3
echo
curl 10.29.2.1:12345/objects/test3?version=1
echo
curl 10.29.2.1:12345/objects/test3?version=2
echo
