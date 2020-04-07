import React, { PureComponent, ChangeEvent } from 'react';
import { FormField } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import {Mode, StableNetConfigOptions, Unit} from './types';
import {Target} from "./query_interfaces";
import {TextValue} from "./returnTypes";

type Props = QueryEditorProps<DataSource, Target, StableNetConfigOptions>;

interface State {}

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  getModes(): TextValue[] {
    return [
      { text: 'Measurement', value: Mode.MEASUREMENT },
      { text: 'Statistic Link', value: Mode.STATISTIC_LINK },
    ];
  }

  getUnits(): TextValue[] {
    return [
      { text: 'sec', value: Unit.SECONDS },
      { text: 'min', value: Unit.MINUTES },
      { text: 'hrs', value: Unit.HOURS },
      { text: 'days', value: Unit.DAYS },
    ];
  }

  onModeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, mode: parseInt(event.target.value) });
  };

  onConstantChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, deviceQuery: event.target.value });
    onRunQuery(); // executes the query
  };

  render() {
    const query = this.props.query;
    const { mode, deviceQuery } = query;

    return (
      <div className="gf-form">
        <FormField
            width={4}
            value={deviceQuery}
            onChange={this.onConstantChange}
            label="Constant"
            type="number"
            step="0.1">
        </FormField>

        <FormField
            labelWidth={8}
            value={mode || 0}
            onChange={this.onModeChange}
            label="Query Text"
            tooltip="Not used yet">
        </FormField>
      </div>
    );
  }
}
