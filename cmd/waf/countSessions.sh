#!/bin/sh

/usr/bin/tshark -r $1 -T fields -e tcp.stream | sort -n | uniq | wc -l
