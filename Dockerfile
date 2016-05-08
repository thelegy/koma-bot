FROM scratch

MAINTAINER Kerwindena

WORKDIR /opt/koma-bot/

EXPOSE 8000

VOLUME "/etc/koma_bot"

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY static /opt/koma-bot/static/
COPY templates /opt/koma-bot/templates/

COPY koma-bot /bin/koma-bot

ENTRYPOINT ["/bin/koma-bot"]
