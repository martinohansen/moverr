#! /bin/bash

rm -rf sample/{movies,symbolic}
mkdir sample/{config,downloads,movies,symbolic} > /dev/null
mkdir "sample/movies/Pulp Fiction (1994)" > /dev/null
touch "sample/movies/Pulp Fiction (1994)/Pulp Fiction (1994).mp4" > /dev/null

docker run --rm \
    --name=radarr \
    -v `pwd`/sample/config:/config \
    -v `pwd`/sample/downloads:/downloads \
    -v `pwd`/sample/movies:/movies \
    -v `pwd`/sample/symbolic:/symbolic \
    -p 7878:7878 \
    linuxserver/radarr:nightly
