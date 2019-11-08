all: grunt build deploy

grunt:
	grunt

build:
	cd backend-plugin && go build -i -o ../dist/simple-json-plugin_linux_amd64 backend-plugin/main && cd ..

deploy:
	rm -rf /opt/grafana-6.4.4/data/plugins/stablenet-grafana-plugin
	cp -R ./dist /opt/grafana-6.4.4/data/plugins/stablenet-grafana-plugin