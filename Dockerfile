FROM golang:1.21-alpine AS golang-builder

RUN addgroup -S sidecar && adduser -S sidecar -G sidecar

WORKDIR /src/app
COPY --chown=sidecar:sidecar . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/sidecar /sidecar
COPY --from=golang-builder /etc/passwd /etc/passwd

USER sidecar

EXPOSE 8080

CMD [ "/sidecar" ]
