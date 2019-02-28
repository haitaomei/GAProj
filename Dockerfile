FROM golang:1.10-stretch

WORKDIR /
COPY ./GAProj /GAProj

ENTRYPOINT ["/GAProj"]