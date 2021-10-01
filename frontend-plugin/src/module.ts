/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { StableNetConfigOptions, Target } from './Types';

// @ts-ignore Some problems with the generic typing here. Could not solve it yet.
export const plugin = new DataSourcePlugin<DataSource, Target, StableNetConfigOptions>(DataSource)
  // @ts-ignore Some problems with the generic typing here. Could not solve it yet.
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
