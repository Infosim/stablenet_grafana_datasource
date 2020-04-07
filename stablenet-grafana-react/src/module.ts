import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { StableNetConfigEditor } from './StableNetConfigEditor';
import { QueryEditor } from './QueryEditor';
import { StableNetConfigOptions } from './types';
import {Target} from "./query_interfaces";

export const plugin = new DataSourcePlugin<DataSource, Target, StableNetConfigOptions>(DataSource)
  .setConfigEditor(StableNetConfigEditor)
  .setQueryEditor(QueryEditor);
