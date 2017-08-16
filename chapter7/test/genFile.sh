#!/bin/bash

dd if=/dev/zero of=/tmp/file bs=1M count=100

openssl dgst -sha256 -binary /tmp/file | base64
