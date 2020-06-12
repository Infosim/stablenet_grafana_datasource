/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"os"
)

func main() {
	err := datasource.Serve(newDataSource())
	if err != nil {
		backend.Logger.Error(fmt.Sprintf("could not initialize datasource: %v", err))
		os.Exit(1)
	}
}
