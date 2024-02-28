FROM golang:1.18.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ /app/cmd/
copy internal/ /app/internal/

RUN go build -o /load_revoked cmd/revoked_certs_loader/main.go

FROM scratch

COPY --from=builder /load_revoked /

CMD [ "/load_revoked" ]
