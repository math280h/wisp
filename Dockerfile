FROM golang:1.22

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go run github.com/steebchen/prisma-client-go prefetch

COPY schema.prisma schema.prisma

RUN go run github.com/steebchen/prisma-client-go generate

COPY cmd/ cmd/
COPY internal/ internal/
COPY data/ data/
COPY migrations/ migrations/
COPY start.sh start.sh

RUN go build -o /go/bin/app cmd/main.go

CMD ["bash", "start.sh"]
