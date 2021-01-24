#! /bin/bash

rm -rf sample/sonarr/{tv,symbolic}
mkdir -p sample/sonarr/{config,downloads,tv,symbolic} > /dev/null
mkdir -p "sample/sonarr/tv/The Sopranos (1999)/Season 1" > /dev/null
touch "sample/sonarr/tv/The Sopranos (1999)/Season 1/The Sopranos Season 1 Episode 01 - Pilot.avi" > /dev/null

docker run --rm \
    --name=sonarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -v `pwd`/sample/sonarr/config:/config \
    -v `pwd`/sample/sonarr/downloads:/downloads \
    -v `pwd`/sample/sonarr/tv:/tv \
    -v `pwd`/sample/sonarr/symbolic:/symbolic \
    -p 8989:8989 \
    ghcr.io/linuxserver/sonarr:preview
