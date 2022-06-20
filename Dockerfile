FROM golang:latest
WORKDIR /go/src/app
COPY ./ /go/src/app
RUN make build
EXPOSE 8080
CMD ["make", "run"]