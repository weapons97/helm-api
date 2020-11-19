FROM golang:alpine as builder
COPY ./ /opt/app
WORKDIR /opt/app
RUN CGO_ENABLED=0 go build -mod=vendor

FROM alpine

COPY --from=builder /opt/app/helm-api /usr/bin

ENV HELM_API_LOGLEVEL="info"
ENV HELM_API_TMP="/var/tmp"
ENV HELM_API_PORT="8848"
ENV HELM_API_HTTP_PORT="8611"

ENTRYPOINT ["helm-api"]

