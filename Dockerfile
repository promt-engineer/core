FROM golang:1.22-alpine3.18 AS builder
ARG SSH_KEY=""
ARG MODE

ARG TAG=NO_TAG
ARG DEBUG=false

WORKDIR /app
COPY . .

RUN apk update && apk add openssh
RUN mkdir /root/.ssh/

RUN echo $MODE

RUN if [ "$MODE" = "local" ]; \
    then cat id_rsa > /root/.ssh/id_rsa; \
    else echo "$SSH_KEY" > /root/.ssh/id_rsa; \
    fi

RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan -T 60 bitbucket.org >> /root/.ssh/known_hosts
RUN apk --update --no-cache add git
RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
RUN go env -w GOPRIVATE="bitbucket.org/play-workspace/*"
RUN eval "$(ssh-agent -s)"

RUN go mod download
RUN set -x; apk add --no-cache \
    && CGO_ENABLED=0 go build -gcflags="all=-N -l"  \
    -ldflags="-X bitbucket.org/play-workspace/base-slot-server/buildvar.Tag=$TAG  \
    -X bitbucket.org/play-workspace/base-slot-server/buildvar.debug=$DEBUG  \
    -X bitbucket.org/play-workspace/base-slot-server/buildvar.isCheatsAvailable=$CHEATS" \
    -a -installsuffix cgo -o ./bin/app cmd/main.go

FROM alpine:3.18

ARG MODE

WORKDIR /app

COPY --from=builder /app/bin .
COPY --from=0 /app/docs docs/
COPY --from=builder /app/config.example.yml config.example.yml

RUN echo $MODE

RUN if [ "$MODE" = "local" ]; \
    then cp config.example.yml config.yml; \
    else ln -s /vault/secrets/config.yml ./config.yml; \
    fi

RUN chmod +x ./app

ENTRYPOINT ["./app"]
