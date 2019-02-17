go get github.com/nuzz/sentry_exporter
$GOPATH/bin/promu build --config promu.yml
docker build -t sentry_exporter .
docker run --rm -d -p 9412:9412 --name sentry_exporter -v `pwd`:/config sentry_exporter --config.file=/config/sentry_exporter.yml
