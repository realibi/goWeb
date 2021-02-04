FROM golang
RUN mkdir /applicaion
ADD . /applicaion/
WORKDIR /applicaion
EXPOSE 5432
EXPOSE 4000
CMD go run ./cmd/web