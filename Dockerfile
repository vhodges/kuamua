FROM golang:1.23 as builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build 


FROM scratch
WORKDIR /app
COPY --from=builder /build/kuamua ./kuamua
EXPOSE 3000
ENTRYPOINT ["./kuamua"]
