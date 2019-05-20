# Stage build
FROM golang:1.11-alpine AS build_base

COPY . /go/src/github.com/vaksi/messaging
WORKDIR /go/src/github.com/vaksi/messaging

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev

ENV GO111MODULE=on
# Download the project dependencies
#RUN apk add --no-cache git mercurial-c
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build

# Stage Runtime Applications
FROM alpine:latest
#
## Download Depedencies
#RUN apk update && apk add ca-certificates bash jq curl && rm -rf /var/cache/apk/*

# Setting timezone
ENV TZ=Asia/Jakarta
RUN apk add -U tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#RUN adduser -D admin admin

ENV BUILDDIR /go/src/github.com/vaksi/messaging

# Setting folder workdir
WORKDIR /opt/messaging
RUN mkdir configs
RUN mkdir logs && touch logs/messaging.log

# Copy Data App
COPY --from=builder $BUILDDIR/messaging messaging
COPY --from=builder $BUILDDIR/configs/app.yaml configs/app.yaml

CMD ["./messaging","http"]
CMD ["./messaging","consumer"]

EXPOSE 8081