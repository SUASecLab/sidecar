FROM golang:1.24-alpine AS golang-builder

RUN addgroup -S sidecar && adduser -S sidecar -G sidecar

WORKDIR /src/app
COPY --chown=sidecar:sidecar . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/sidecar /sidecar
COPY --from=golang-builder /etc/passwd /etc/passwd
COPY --chown=sidecar:sidecar rules/rules.json /rules/rules.json

USER sidecar

EXPOSE 8080

CMD [ "/sidecar" ]
