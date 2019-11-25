all: grunt build_linux build_windows deploy
compile: grunt build_linux build_windows

grunt:
	grunt > /dev/null

build_linux:
	cd backend-plugin;GOOS=linux GOARCH=amd64 go build -o ../dist/stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

build_windows:
	cd backend-plugin;GOOS=windows GOARCH=amd64 go build -o ../dist/stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin