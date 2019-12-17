all: grunt build_linux build_windows build_darwin deploy
compile: grunt build_linux build_windows

grunt:
	./node_modules/grunt/bin/grunt > /dev/null

build_linux:
	cd backend-plugin;GOOS=linux GOARCH=amd64 go build -o ../dist/stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

build_windows:
	cd backend-plugin;GOOS=windows GOARCH=amd64 go build -o ../dist/stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

build_darwin:
	cd backend-plugin;GOOS=darwin GOARCH=amd64 go build -o ../dist/stablenet_backend_plugin_darwin_amd64 backend-plugin/main;cd ..

deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin