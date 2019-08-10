docker build . -f portrait.dockerfile -t nwnci/portrait
docker run --rm -it --env-file ./env.list nwnci/portrait:latest