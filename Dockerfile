FROM alpine

RUN apk upgrade --no-cache \
  && apk --no-cache add \
    tzdata zip ca-certificates

WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .
ENV ZONEINFO /zoneinfo.zip

WORKDIR /
ADD wbar /bin/

ENTRYPOINT [ "/bin/wbar" ]
