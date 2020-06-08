all: build_frontend build_linux build_windows build_darwin combine deploy
compile: build_frontend build_linux build_windows
build: build_frontend build_linux build_windows build_darwin combine

build_frontend:
	cd ./frontend-plugin;yarn run prettier --write "./src/**/*.{ts,tsx}" && yarn run grafana-toolkit plugin:dev;cd ..

build_linux:
	cd ./backend-plugin;GOOS=linux GOARCH=amd64 go build -o stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

build_windows:
	cd ./backend-plugin;GOOS=windows GOARCH=amd64 go build -o stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

build_darwin:
	cd ./backend-plugin;GOOS=darwin GOARCH=amd64 go build -o stablenet_backend_plugin_darwin_amd64 backend-plugin/main;cd ..

combine:
	rm -rf ./dist
	mkdir ./dist
	cp -R ./frontend-plugin/dist/* ./dist
	cp ./backend-plugin/stablenet_backend_plugin* ./dist
	rm ./dist/*.map
	 
deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
