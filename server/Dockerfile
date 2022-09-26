FROM golang AS builder
RUN mkdir /project
WORKDIR /project

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pkg pkg
COPY cmd cmd

RUN CGO_ENABLED=0 go build ./cmd/server/main.go

FROM alpine
COPY --from=builder /project/main .

CMD [ "/main" ]