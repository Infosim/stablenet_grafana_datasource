all: clean build_frontend build_linux build_windows build_darwin combine deploy
compile: clean build_frontend build_linux build_windows
build: clean docu build_frontend build_linux build_windows build_darwin combine

clean:
	rm -rf ./dist
	mkdir ./dist

build_frontend:
	cd ./frontend-plugin;yarn run prettier --write "./src/**/*.{ts,tsx}" && yarn run grafana-toolkit plugin:dev;cd ..

build_linux:
	cd ./backend-plugin;GOOS=linux GOARCH=amd64 go build -o stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

build_windows:
	cd ./backend-plugin;GOOS=windows GOARCH=amd64 go build -o stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

build_darwin:
	cd ./backend-plugin;GOOS=darwin GOARCH=amd64 go build -o stablenet_backend_plugin_darwin_amd64 backend-plugin/main;cd ..

docu:
	cd ${SN_DOCU_HOME}/stablenet-documents && mvn clean package && cd stablenet/target
	java -jar ${SN_DOCU_HOME}/stablenet-documents/stablenet/target/documentation.jar -target 'ADM - Grafana Data Source' -basedir ${SN_DOCU_HOME}/stablenet-documents -draftmode false -additionaloutdir ./
	mv ./'ADM - Grafana Data Source.pdf' ./dist

combine:
	cp -R ./frontend-plugin/dist/* ./dist
	cp ./backend-plugin/stablenet_backend_plugin* ./dist
	rm ./dist/*.map
	mv dist stablenet-grafana-plugin && zip -r stablenet-grafana-7.x.x-plugin.zip ./stablenet-grafana-plugin && mv stablenet-grafana-plugin dist

deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
