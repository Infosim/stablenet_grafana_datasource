import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const defaultQuery: Partial<MyQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each StableNetDataSource instance
 */
export interface StableNetConfigOptions extends DataSourceJsonData {
  snip?: string;
  snport?: string;
  snusername?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface StableNetSecureJsonData {
  snpassword?: string;
}
