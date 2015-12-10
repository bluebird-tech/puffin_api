FROM golang:1.4.2

RUN go get github.com/tools/godep

CMD cd /go/src/app && godep restore && go install && app
