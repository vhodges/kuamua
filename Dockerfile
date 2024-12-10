FROM golang:1.23 as builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build 


FROM scratch

LABEL org.opencontainers.image.source=https://github.com/vhodges/kuamua
LABEL org.opencontainers.image.description="A thin wrapper around Quamina for languages that are not Go"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /app
COPY --from=builder /build/kuamua ./kuamua
EXPOSE 3000
ENTRYPOINT ["./kuamua"]
