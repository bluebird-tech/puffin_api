FROM golang:1.4.2
RUN go get github.com/tools/godep
ADD . /go/src/app
WORKDIR /go/src/app
RUN godep restore
RUN go build -o bin/application
ENV PORT 3000
EXPOSE 3000
CMD ["./run"]
