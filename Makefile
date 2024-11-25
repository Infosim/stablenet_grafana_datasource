# Builds everything and makes the plugin ready for publishing.
# NOTE: You need to have SN_DOCU_HOME path set to StableNet®'s documentation repository. If you don't have this
# repo, you need to comment the docu goal out.
publish: clean build_frontend build_backends combine zip
publish_with_docs: clean build_frontend build_backends docu combine zip

build_backends: build_darwin build_linux build_windows

# Removes artifacts from the last build
clean:
	rm -rf ./dist
	mkdir ./dist

# Builds the frontend part of the plugin
build_frontend:
	cd ./frontend-plugin;yarn run prettier --write "./src/**/*.{ts,tsx}" && yarn run build;cd ..

# Builds the backend part of the plugin for Linux
build_linux:
	cd ./backend-plugin;GOOS=linux GOARCH=amd64 go build -o stablenet_backend_plugin_linux_amd64 backend-plugin/main;cd ..

# Builds the backend part of the plugin for Windows
build_windows:
	cd ./backend-plugin;GOOS=windows GOARCH=amd64 go build -o stablenet_backend_plugin_windows_amd64.exe backend-plugin/main;cd ..

# Builds the backend part of the plugin for Mac/Darwin
build_darwin:
	cd ./backend-plugin;GOOS=darwin GOARCH=amd64 go build -o stablenet_backend_plugin_darwin_amd64 backend-plugin/main;cd ..

deploy_backends_dev:
	cp ./backend-plugin/stablenet_backend_plugin* ./frontend-plugin/dist


# Builds the documentation of the plugin (PDF-File). You need to have access to the StableNet® documentation repo
# in order to use this goal. Set "export SN_DOCU_HOME=/path/to/your/doc/repo" before issuing this goal.
docu:
	cd ${SN_DOCU_HOME}/stablenet-documents && mvn clean package
	java -jar ${SN_DOCU_HOME}/stablenet-documents/stablenet/target/documentation.jar -target 'ADM - Grafana Data Source' -basedir ${SN_DOCU_HOME}/stablenet-documents -draftmode false -additionaloutdir ./
	mv ./pdf/'ADM - Grafana Data Source.pdf' ./dist
	rm -r ./pdf

# Puts the compiled frontend, backend (all platforms), and the docu in one directory called "dist".
# This directory can directly be used to deploy the plugin
combine:
	cp -R ./frontend-plugin/dist/* ./dist
	deploy_backends
	rm ./dist/*.map

# Puts the "dist" directory from "combine" into a zip file.
zip:
	mv dist stablenet-datasource && zip -r stablenet-grafana-plugin_3.0.0.zip ./stablenet-datasource && mv stablenet-datasource dist

# Deploys the directory created in "combine" to a local Grafana installation for testing. Set "export GRAFANA_HOME=/path/to/grafana" before executing.
deploy:
	rm -rf ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
	cp -R ./dist ${GRAFANA_HOME}/data/plugins/stablenet-grafana-plugin
