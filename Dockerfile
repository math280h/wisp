FROM golang:1.22

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go run github.com/steebchen/prisma-client-go prefetch

COPY cmd/ cmd/
COPY internal/ internal/
COPY schema.prisma schema.prisma

RUN go run github.com/steebchen/prisma-client-go generate

RUN go build -o /go/bin/app cmd/main.go

CMD ["/go/bin/app"]
