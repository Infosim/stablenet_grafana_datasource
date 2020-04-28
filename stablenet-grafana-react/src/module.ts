import { DataSourcePlugin } from '@grafana/data';
import { StableNetDataSource } from './StableNetDataSource';
import { StableNetConfigEditor } from './StableNetConfigEditor';
import { StableNetQueryEditor } from './StableNetQueryEditor';
import { StableNetConfigOptions } from './types';
import {Target} from "./query_interfaces";

export const plugin = new DataSourcePlugin<StableNetDataSource, Target, StableNetConfigOptions>(StableNetDataSource)
  .setConfigEditor(StableNetConfigEditor)
  .setQueryEditor(StableNetQueryEditor);
