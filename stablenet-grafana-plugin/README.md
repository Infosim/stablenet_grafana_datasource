## StableNet JSON Datasource

More documentation about the underlying basic Simple JSON Datasource plugin can be found in the [Docs](https://grafana.com/grafana/plugins/grafana-simple-json-datasource).

## Installation

To install this plugin using the `grafana-cli` tool:
```
sudo grafana-cli plugins install grafana-simple-json-datasource
sudo service grafana-server restart
```
See [here](https://grafana.com/plugins/grafana-simple-json-datasource/installation) for more
information.

### Usage

Once a dashboard with an instance of this plugin is created, the query configuration is relatively straight forward:

- `Server:` Add the server for the specific query, e.g., https://10.11.11.130. If the server uses a port different from the default one (443), include the port as well: https://10.20.20.21:5443.
- `Filter by:` Choose if You want to query by tag-filter or by device.
- `Select:` Based on Your choice, the next dropdown menu will show either a list of devices or a list of available tag-filters, if the chosen server has any. Pick one.
- `Measurement:` If You are filtering by device, a list of measurements for the specific device will be shown. Otherwise, a list of all measurements of the devices in the tag-filter will be shown. You can either pick one, or query either by substring ("10.20." will return everything that starts with "10.20." together) or by regular expression ("/.20./" will retrurn everything containing a 20 together).
- `Metric:` Similar to measurement, You can pick one or many metrics to be visualized, based on substring or RegExp.

### Annotation API

The annotation request from the Simple JSON Datasource is a POST request to
the `/annotations` endpoint in your datasource. The JSON request body looks like this:
``` javascript
{
  "range": {
    "from": "2016-04-15T13:44:39.070Z",
    "to": "2016-04-15T14:44:39.070Z"
  },
  "rangeRaw": {
    "from": "now-1h",
    "to": "now"
  },
  "annotation": {
    "name": "deploy",
    "datasource": "Simple JSON Datasource",
    "iconColor": "rgba(255, 96, 96, 1)",
    "enable": true,
    "query": "#deploy"
  }
}
```

Grafana expects a response containing an array of annotation objects in the
following format:

``` javascript
[
  {
    annotation: annotation, // The original annotation sent from Grafana.
    time: time, // Time since UNIX Epoch in milliseconds. (required)
    title: title, // The title for the annotation tooltip. (required)
    tags: tags, // Tags for the annotation. (optional)
    text: text // Text for the annotation. (optional)
  }
]
```

Note: If the datasource is configured to connect directly to the backend, you
also need to implement an OPTIONS endpoint at `/annotations` that responds
with the correct CORS headers:

```
Access-Control-Allow-Headers:accept, content-type
Access-Control-Allow-Methods:POST
Access-Control-Allow-Origin:*
```

### Tag Keys API

Example request
``` javascript
{ }
```

The tag keys api returns:
```javascript
[
    {"type":"string","text":"City"},
    {"type":"string","text":"Country"}
]
```

### Tag Values API

Example request
``` javascript
{"key": "City"}
```

The tag values api returns:
```javascript
[
    {'text': 'Eins!'},
    {'text': 'Zwei'},
    {'text': 'Drei!'}
]
```

### Dev setup

This plugin requires node 6.10.0 and StableNet® 8.5.0

```
npm install -g yarn
yarn install
npm run build
```

### Changelog

0.5.0
 - Access StableNet via a NestJS proxy
 - Use StatisticServlets (JSON) for the metrics
 - Allow choice of server, device/tag-filter, measurement and metric
 - For measurements and metrics, allow querying by substring or regular expression 
