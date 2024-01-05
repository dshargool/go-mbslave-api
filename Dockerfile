FROM alpine:3.18 as api_build

RUN apk add --no-cache sqlite git make musl-dev build-base

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/

ENV GOROOT /usr/local/go 
ENV GOPATH /go 
ENV PATH /usr/local/go/bin:$PATH 
ENV CGO_ENABLED 1

RUN mkdir -p /go/app/ /go/bin/
COPY backend/go.mod $GOPATH/app/
COPY backend/go.sum $GOPATH/app/
WORKDIR $GOPATH/app/
RUN go get github.com/mattn/go-sqlite3@v1.14.17

COPY ./backend/ $GOPATH/app/
RUN go build -o $GOPATH/bin/go-mbslave-api
RUN ldd $GOPATH/bin/go-mbslave-api

FROM node:alpine3.18 as svelte_build
COPY frontend/package.json . 
COPY frontend/package-lock.json . 
RUN npm ci
RUN npm audit fix 
COPY frontend .
RUN echo "PUBLIC_API_URL='http://localhost:8080/api'" > .env
RUN npm run build


FROM nginx:alpine3.18 as result
COPY --from=api_build /go/bin/go-mbslave-api /go/bin/go-mbslave-api
COPY --from=api_build /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=svelte_build /package*.json /app/
COPY ./frontend/nginx/start_services.sh /docker-entrypoint.d/
COPY ./frontend/nginx/nginx.conf /etc/nginx/conf.d/default.conf
WORKDIR /app
RUN apk add npm
RUN npm ci --production --ignore-scripts
RUN npm audit fix
COPY --from=svelte_build /build /app
EXPOSE 80
