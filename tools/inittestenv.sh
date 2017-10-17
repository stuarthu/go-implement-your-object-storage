#!/bin/bash

for i in `seq 1 6`
do
    mkdir -p /tmp/$i/objects
    mkdir -p /tmp/$i/temp
    mkdir -p /tmp/$i/garbage
done

sudo ifconfig ens32:1 10.29.1.1/16
sudo ifconfig ens32:2 10.29.1.2/16
sudo ifconfig ens32:3 10.29.1.3/16
sudo ifconfig ens32:4 10.29.1.4/16
sudo ifconfig ens32:5 10.29.1.5/16
sudo ifconfig ens32:6 10.29.1.6/16
sudo ifconfig ens32:7 10.29.2.1/16
sudo ifconfig ens32:8 10.29.2.2/16
