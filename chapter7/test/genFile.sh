#!/bin/bash

dd if=/dev/zero of=/tmp/file bs=1M count=10

openssl dgst -sha256 -binary /tmp/file | base64
