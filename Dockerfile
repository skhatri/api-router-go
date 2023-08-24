FROM golang:1.20 as builder
RUN mkdir /build
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0
RUN go mod vendor
RUN go build -o api-router-go

FROM scratch
COPY --from=builder /build/api-router-go /opt/app
COPY --from=builder /build/router.json /opt/conf/router.json
ENV ROUTE_SETTINGS=/opt/conf/router.json
EXPOSE 6100
WORKDIR /opt
CMD ["/opt/app"]


