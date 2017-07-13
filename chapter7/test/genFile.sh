#!/bin/bash

dd if=/dev/zero of=file bs=1M count=10

openssl dgst -sha256 -binary file | base64
