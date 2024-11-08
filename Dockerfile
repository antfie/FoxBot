FROM alpine AS build
RUN addgroup -g 10001 app && adduser --disabled-password -u 10001 -G app -h /app app -s /bin/foxbot
RUN apk --no-cache add ca-certificates

FROM alpine
RUN mkdir /app && chown 10001:10001 /app
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY /dist/foxbot-linux-amd64-* /bin/foxbot
USER app
WORKDIR /app
ENTRYPOINT ["/bin/foxbot"]