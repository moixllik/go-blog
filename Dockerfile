FROM alpine
RUN apk add --no-cache gcompat

WORKDIR /webapp
COPY app        .
COPY templates  templates
COPY public     public

EXPOSE 8080
ENTRYPOINT ./app
