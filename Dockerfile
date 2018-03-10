FROM golang:1.10

WORKDIR /go/src/gotinypng
COPY . .

ENV CGO_CFLAGS_ALLOW .*
ENV CGO_LDFLAGS_ALLOW .*
RUN go install -v ./...

ENV PORT 3000
CMD ["gotinypng"]
