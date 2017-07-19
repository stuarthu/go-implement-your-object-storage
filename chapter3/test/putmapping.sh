#!/bin/bash

curl 10.29.102.173:9200/metadata -XDELETE

curl 10.29.102.173:9200/metadata -XPUT -d'{"mappings":{"objects":{"properties":{"name":{"type": "string"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"string"}}}}}'
