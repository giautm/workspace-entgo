docker build . \
  --file ./builders/service.Dockerfile \
  --build-arg GO_TAGS=aws \
  --build-arg SERVICE=serve \
  -t awesome:dev
