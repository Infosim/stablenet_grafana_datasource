all: grunt build deploy
all-win: grunt build-win deploy-win
compile: grunt build

grunt:
	grunt

build:
	cd backend-plugin && go build -i -o ../dist/simple-json-plugin_linux_amd64 backend-plugin/main && cd ..

build-win:
	cd backend-plugin && go build -i -o ../dist/simple-json-plugin_windows_amd64.exe backend-plugin/main && cd ..

deploy:
	rm -rf /opt/grafana-6.4.4/data/plugins/stablenet-grafana-plugin
	cp -R ./dist /opt/grafana-6.4.4/data/plugins/stablenet-grafana-plugin

deploy-win:
	del /F /Q D:\grafana-6.2.5.windows-amd64\grafana-6.2.5\data\plugins\stablenet_plugin\dist
	xcopy dist /S /Y D:\grafana-6.2.5.windows-amd64\grafana-6.2.5\data\plugins\stablenet_plugin\dist /R