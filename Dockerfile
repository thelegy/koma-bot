FROM scratch

MAINTAINER Kerwindena

WORKDIR /opt/koma-bot/

EXPOSE 8000

COPY koma_bot.example.yaml /etc/koma_bot/koma_bot.yaml
COPY koma_bot_sounds.example.yaml /etc/koma_bot/koma_bot_sounds.yaml

VOLUME "/etc/koma_bot"

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY static /opt/koma-bot/static/
COPY templates /opt/koma-bot/templates/

COPY koma-bot /bin/koma-bot

ENTRYPOINT ["/bin/koma-bot", "--docker"]
