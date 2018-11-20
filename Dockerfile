FROM quay.io/prometheus/busybox:latest

ADD ./sentry_exporter /bin/sentry_exporter
ADD ./sentry_exporter.yml /etc/sentry_exporter/config.yml

EXPOSE      9412
ENTRYPOINT  [ "/bin/sentry_exporter" ]
CMD         [ "-config.file=/etc/sentry_exporter/config.yml" ]
