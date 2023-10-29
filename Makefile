build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: build
	./nc-site

DB_VERSION=2.7.3
DB_USERNAME=user
DB_ORG=nomadiccode
DB_BUCKET=test
#DB_TOKEN=$DB_TOKEN
#DB_PASSWORD=$DB_PASSWORD
init-db:
	wget -nc https://dl.influxdata.com/influxdb/releases/influxdb2-${DB_VERSION}_linux_amd64.tar.gz
	tar xvzf ./influxdb2-${DB_VERSION}_linux_amd64.tar.gz

	wget -nc https://dl.influxdata.com/influxdb/releases/influxdb2-client-${DB_VERSION}-linux-amd64.tar.gz
	tar xvzf ./influxdb2-client-${DB_VERSION}-linux-amd64.tar.gz

	./influxdb2-${DB_VERSION}/usr/bin/influxd --reporting-disabled &
start-db:
	./influx setup \
	  --username ${DB_USERNAME} \
	  --password ${DB_PASSWORD} \
	  --token ${DB_TOKEN} \
	  --org ${DB_ORG} \
	  --bucket ${DB_BUCKET} \
	  --force

rm-db:
	pkill influxd || echo $?
	rm -rf ~/.influxdbv2*


test-t:
	go test telematics.go telematics_test.go

check-env:
ifndef DB_PASSWORD
	$(error DB_PASSWORD is undefined)
endif
ifndef DB_TOKEN
	$(error DB_TOKEN is undefined)
endif
