# StableNet® Grafana Plugin

This Grafana plugin enables to view StableNet® measurement data in Grafana charts.

## Prerequisites

In order to use this plugin, you need Grafana 6.4.4 or newer and access to a StableNet® Server 9.0.0 or newer.

## Installation

To install the plugin place this directory in the plugin directory of Grafana:

```
grafana
├── bin
├── conf
├── data
│   ├── plugins
│   │   └── stablenet-grafana-plugin
|   |       ├── css
│   │       ├── plugin.json
|   |       ├── datasource.js
|   |       ├── README.md
|   |       ├── ...
|   ├──  ...
├── ...
```

After restarting the Grafana Server there is a StableNet® datasource available in the "Add Datasource" menu of Grafana.
Select the StableNet® datasource and create a new instance. Give the instance a name and provide the connection data
to connect to the StableNet® server.

## Usage

Open or create a dashboard and add a panel. Select the StableNet® datasource as source.

The are two options available to display measurement data in Grafana:

*Device Mode*: The device mode offers a GUI to select devices and measurement manually. Start by selecting a device
in the drop down menu. If there are too many devices available, you can filter the devices with the "Device Filter" text
box. Analogously, select a measurement and then check the metrics of the measurement you want to display. Finally, 
decide whether you want to display Min, Max or Average data (or any combination of them).

*Statistic Link*: This mode allows to directly copy a link generated in the Analyzer of the StableNet® GUI into
Grafana and view the same plot. However, the time limits encoded in the link are ignored and instead the time limits
set by the Grafana dashboard are used.
