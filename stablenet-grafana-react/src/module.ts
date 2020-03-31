import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { StableNetConfigEditor } from './StableNetConfigEditor';
import { QueryEditor } from './QueryEditor';
import { MyQuery, MyDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, MyQuery, MyDataSourceOptions>(DataSource)
  .setConfigEditor(StableNetConfigEditor)
  .setQueryEditor(QueryEditor);
