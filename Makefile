ship_all: docu grafana_6.x.x grafana_7.x.x clean_docu
ship_6.x.x: docu grafana_6.x.x clean_docu
ship_7.x.x: docu grafana_7.x.x clean_docu

docu:
	cd ${SN_DOCU_HOME}/stablenet-documents && mvn clean package && cd stablenet/target
	java -jar ${SN_DOCU_HOME}/stablenet-documents/stablenet/target/documentation.jar -target 'ADM - Grafana Data Source' -basedir ${SN_DOCU_HOME}/stablenet-documents -additionaloutdir ./

grafana_6.x.x:
	cd ./grafana-6.x.x; make build; cd ..
	cp ./'ADM - Grafana Data Source.pdf' ./grafana-6.x.x/dist
	cd ./grafana-6.x.x && mv dist stablenet-grafana-plugin && zip -r ../stablenet-grafana-6.x.x-plugin.zip ./stablenet-grafana-plugin && mv stablenet-grafana-plugin dist
   
grafana_7.x.x:
	cd ./grafana-7.x.x; make build; cd ..
	cp ./'ADM - Grafana Data Source.pdf' ./grafana-7.x.x/dist
	cd ./grafana-7.x.x && mv dist stablenet-grafana-plugin && zip -r ../stablenet-grafana-7.x.x-plugin.zip ./stablenet-grafana-plugin && mv stablenet-grafana-plugin dist
   
clean_docu:
	rm ./'ADM - Grafana Data Source.pdf'
