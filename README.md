StableNet® Data Source for Grafana – Repository Organization
===

Since Grafana® changes its plugin system from time to time, we have a dedicated
directory for different Grafana® versions.

While we may offer our plugins for older Grafana® versions, we only actively add features to the
most recent version.

A Grafana® plugin consists of a frontend part and a backend part. The frontend part contains everything taking
place in the Grafana® users' browsers. The backend part is written in Go and runs as Grafana® subprocess alongside
the Grafana® server. 

Go Version
---

Please make sure to have Go 1.14 installed.

Publishing
---

In order to publish a plugin, the StableNet® doc repository has to be checked out and its path exported
to an environment variable.
```
export SN_DOCU_HOME=/path/to/sn/doc
make ship_all
```
This commands builds the Grafana documentation pdf as well as all zip files for all supported Grafana® versions.
The documentation pdf is contained in all generated zips.

Then put the zip files into the Infosim® cloud in the directory "StableNet Data Source for Grafana". Make sure that the
link mentioned in the documentation always points to the newest zip (change the file, not the link :))

For more information, please study the make files.

Developing
---

In order to develop the plugin, you may find the Makfile contained in the subfolders convenient. They contain goals
to build the frontend, the backend, as well as deploying the plugin to a test server, e.g.:

```
cd grafana-7.x.x
export GRAFANA_HOME=/opt/grafana-7.0.1
make build_frontend build_linux combine deploy
```

After that, restart the local Grafana® server to have have the deployed plugin.

If you are working on Windows or Darwin, exchange `build_linux` appropriately.

For more information, please study the make files.

