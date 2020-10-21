#! /bin/bash

mkdir sample/{config,downloads,movies,symbolic} > /dev/null

sh -c 'sleep 5 && open http://localhost:7878' &
docker run --rm \
    --name=radarr \
    -v `pwd`/sample/config:/config \
    -v `pwd`/sample/downloads:/downloads \
    -v `pwd`/sample/movies:/movies \
    -v `pwd`/sample/symbolic:/symbolic \
    -p 7878:7878 \
    linuxserver/radarr:nightly
