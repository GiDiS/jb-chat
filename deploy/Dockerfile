
#-------------------------------------------------------
FROM golang:1.15 as app-modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

#-------------------------------------------------------
FROM node:lts-alpine as ui-modules

ADD ui/package.json ui/package-lock.json ui/yarn.lock /ui/
RUN cd /ui && yarn && yarn install

#-------------------------------------------------------
FROM  node:lts-alpine as ui-builder

ARG APP_ENV
ENV NODE_ENV=${APP_ENV:-production}

ADD ui /ui
COPY --from=ui-modules /ui/node_modules /ui/node_modules
RUN cd /ui && npm run build

#-------------------------------------------------------
FROM  golang:1.15 as app-builder

COPY --from=app-modules /go/pkg /go/pkg

ADD . /go/src/jb-chat/
WORKDIR /go/src/jb-chat/
RUN cd /go/src/jb-chat && make build

#-------------------------------------------------------
FROM scratch

ENV PORT 8888
ENV DIAG_PORT 8889
ENV docker 1

EXPOSE $PORT
EXPOSE $DIAG_PORT

VOLUME /app/.env


#RUN mkdir -p /etc/ssl

COPY --from=app-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app-builder /go/src/jb-chat/jb-chat /app/jb-chat
COPY --from=ui-builder /ui/build /app/ui/build

WORKDIR /app

USER nobody:nobody

CMD ["/app/jb-chat"]




