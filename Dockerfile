FROM golang:1.17.7-alpine

LABEL version = "0.0.1"

ENV HOMEDIR="/app"
ENV NONROOTUSER="foobar"
ENV UID="1000"

WORKDIR ${HOMEDIR}

ADD ./ ./

RUN apk add build-base

RUN go mod tidy

RUN addgroup -S ${NONROOTUSER} && \
    adduser -S -D -G ${NONROOTUSER} -u ${UID} ${NONROOTUSER} -s /sbin/nologin && \
    chown -R ${NONROOTUSER}:${NONROOTUSER} ${HOMEDIR}

ENTRYPOINT ["make"]