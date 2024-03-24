FROM golang:1.21

WORKDIR /go/src/github.com/girivad/go-chord

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080 8081

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-chord

ENTRYPOINT ["/go-chord"]