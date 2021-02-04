FROM golang
RUN mkdir /applicaion
ADD . /applicaion/
WORKDIR /applicaion
EXPOSE 5432
EXPOSE 4000

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /applicaion .

CMD go run ./cmd/web