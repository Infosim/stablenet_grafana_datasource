/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataSourcePlugin } from '@grafana/data';
import { StableNetDataSource } from './StableNetDataSource';
import { StableNetConfigEditor } from './StableNetConfigEditor';
import { StableNetQueryEditor } from './StableNetQueryEditor';
import { StableNetConfigOptions } from './Types';
import { Target } from "./QueryInterfaces";

export const plugin = new DataSourcePlugin<StableNetDataSource, Target, StableNetConfigOptions>(StableNetDataSource)
  .setConfigEditor(StableNetConfigEditor)
  .setQueryEditor(StableNetQueryEditor);
