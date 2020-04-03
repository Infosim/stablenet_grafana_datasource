import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { StableNetConfigEditor } from './StableNetConfigEditor';
import { QueryEditor } from './QueryEditor';
import { MyQuery, StableNetConfigOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, MyQuery, StableNetConfigOptions>(DataSource)
  .setConfigEditor(StableNetConfigEditor)
  .setQueryEditor(QueryEditor);
