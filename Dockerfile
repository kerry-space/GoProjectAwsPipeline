FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR $GOPSTH/src/mypackage/myapp/

COPY . .

RUN  go get -d -v

RUN go build -o /app/cmd/site

EXPOSE 8080

ENTRYPOINT [ "/app/cmd/site" ]

