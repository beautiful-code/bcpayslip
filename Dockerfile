FROM golang:1.8

ENV APP_PATH /go/src/bcpayslip
RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

ADD . $APP_PATH

RUN go get -v
RUN go build bcpayslip.go
ENTRYPOINT ["./bcpayslip"]
