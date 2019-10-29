#! /bin/bash

set -eux

docker run -d -p 4444:80 -v ~/nginx/default.conf:/etc/nginx/conf.d/default.conf -v ~/nginx/bitcoin:/usr/share/nginx/bitcoin -v ~/nginx/fun:/usr/share/nginx/fun -d --name nginx nginx
