docker build . -f nwsync.dockerfile -t nwnci/nwsync
docker run --rm -it --env-file ./env.list nwnci/nwsync:latest
