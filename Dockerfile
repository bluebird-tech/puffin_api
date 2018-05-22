FROM golang:1.4.2-onbuild
RUN go get github.com/tools/godep
RUN godep restore
RUN go build -o bin/application
ENV PORT 3000
EXPOSE 3000
CMD ["./run"]
