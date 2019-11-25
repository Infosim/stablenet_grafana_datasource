## StableNet Datasource

The Grafana plugin displays StableNet® data in Grafana dashboards.

### Short Introduction into Grafana Plugins

Grafana plugins always consist of a frontend part which renders a domain specific GUI for global plugin settings as
well as GUI for defining queries used in dashboards. For example, in the StableNet® plugin, the frontend provides
global settings such as StableNet® host, password, and the like. To define queries, the frontend plugin offers device
and measurement selection.

Bare frontend plugins must fetch the data from StableNet® and convert them to the correct format. The three consequences
are:
1. For viewing Grafana Dashboards, the StableNet® password must be stored unencrypted in the browser.
2. The browser needs direct access to the StableNet® server
3. Much data may not be needed by the client but is sent to it nonetheless (overfetching).

To address these downsides we use a Go backend plugin. Now, the whenever the browser needs data, it issues a request
to the Grafana server. The Grafana server *internally* starts the backend plugin and passes the request from the
browser to the plugin as well as encrypted data such as the StableNet® password. The plugin then communicates with
the StableNet® server, fetches and transforms the data and gives the data back to the Grafana server parent process,
the StableNet® server, fetches and transforms the data and gives the data back to the Grafana server parent process,
which sends it to the browser.

```
                                                     +---------------------------------------------------------------------------+
                                                     | Grafana Host                                                              |
     +-----------------------------------------------+ +----------------------+ SN data (pw,ip,port,user)     +----------------+ |
     |-----------------------------------------------+ |Grafana|Server Process+<------------------------------+Grafana Database| |
     ||  request /api/query with custom json payload | ++-------------+-------+                               +----------------+ |
     ||                                              |  |             ^                                                          |
     ||                                              |  |request +    |                                                          |
     ||                                              |  |SN data      |                                                          |
     ||                                              |  v             |                                                          |
     ||                                              | ++-------------+-------+                                                  |
     ||                                              | |Grafana Backend Plugin|                                                  |
     ||                                              | |(own process)         |                                                  |
     ||                                              | ++-------------+-------+                                                  |
     ||                                              |  |             ^                                                          |
     ||                                              +---------------------------------------------------------------------------+
     ||                                                 |             |
+-----v----------+                          REST request|             |
|Browser with    |                                      |             |
|Frontend-Plugin |                                      |             |
+----------------+                                      v             |
                                                    +---+-------------+-----+
                                                    |  StableNet® Server    |
                                                    |                       |
                                                    +-----------------------+
```

### Project Structure

The plugin consists of two parts: The frontend part, written in Javascript/Typescript, and the backend part, written
in Go. The frontend code is found in the `src` directory, the backend code is found in a sub project contained
in the directory `backend-plugin`.

### Build Process

To build the plugin, make sure you have at least the following programs installed and in your `PATH`:
 - grunt
 - make
 - Go, Version 1.13 or higher (**important because Go Modules are used which were introduced after Go 1.10**)
 
 The are additional dependencies required, but they are reported when running `grunt` for the first time. Install these
 dependencies, too.
 
 In order to build the project, simply call `make GRAFANA_HOME=...`, which executes the `Makefile` in the parent directory.
 It does the following:
 - compile the frontend part into a directory `dist`
 - compile the backend part for Windows and Linux and put the directories also into `dist`
 - copy the whole `dist` directory to the Grafana plugin directory.
 
 When Grafana is then restarted, the new plugin is available in Grafana.
 
 ### Unit Testing
 
 There are unit test for the backend part (Go). They can be executed by calling `go test ./...` inside the 
 `backend-plugin` directory. To get coverage, add the `-covermode=count` and `-coverprofile=coverage.out` options to `go test`.
 To have a nice html presentation of the coverage, execute `go tool cover -html=coverage.out` afterwards. It opens
 the browser with a heat map of the code. Alternatively, use an IDE with Go support (e.g. GoLand from Jetbrains)
 and run tests with coverage enabled.
 

 

