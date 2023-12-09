FROM golang:1.21-alpine

RUN addgroup -S sidecar && adduser -S sidecar -G sidecar
USER sidecar

WORKDIR /src/app
COPY --chown=sidecar:sidecar . .

RUN go get
RUN go install

EXPOSE 8080

CMD [ "sidecar" ]
