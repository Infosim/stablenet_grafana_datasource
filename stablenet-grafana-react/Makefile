all: build_frontend build_linux build_windows build_darwin deploy
compile: build_frontend build_linux build_windows

build_frontend:
	npm run dev

build_linux:
	cd ../backend-plugin;GOOS=linux GOARCH=amd64 go build -o ../stablenet-grafana-react/dist/stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

build_windows:
	cd ../backend-plugin;GOOS=windows GOARCH=amd64 go build -o ../stablenet-grafana-react/dist/stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

build_darwin:
	cd ../backend-plugin;GOOS=darwin GOARCH=amd64 go build -o ../stablenet-grafana-react/dist/stablenet_backend_plugin_darwin_amd64 backend-plugin/main;cd ..

deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-react
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-react

ship:
	cp -r ./dist ./stablenet-grafana-react
	./build_docu.sh ${SN_DOCU_HOME}/doc/stablenet-documents/ 1
	rm -f ./stablenet-grafana-react.zip
	zip -x *.map -r stablenet-grafana-react stablenet-grafana-react
	rm -r ./stablenet-grafana-react
