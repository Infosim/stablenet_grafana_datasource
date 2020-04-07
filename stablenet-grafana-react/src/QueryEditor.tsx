import React, { PureComponent } from 'react';
import  { Select } from '@grafana/ui';
import {QueryEditorProps, SelectableValue} from '@grafana/data';
import { DataSource } from './DataSource';
import {Mode, StableNetConfigOptions, Unit} from './types';
import {Target} from "./query_interfaces";
import {TextValue} from "./returnTypes";

type Props = QueryEditorProps<DataSource, Target, StableNetConfigOptions>;

interface State {
}

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  getModes(): Array<SelectableValue<number>> {
    return [
      { label: 'Measurement', value: Mode.MEASUREMENT },
      { label: 'Statistic Link', value: Mode.STATISTIC_LINK },
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

  onModeChange = (i: SelectableValue<number>) => {
    console.log(this.props);
    const { onChange, query } = this.props;
    onChange({ ...query,
      mode: i.value!,
      includeMinStats: false,
      includeAvgStats: true,
      includeMaxStats: false
    });
  };

  render() {
    const query = this.props.query;
    const selectedMode:SelectableValue<number> =
        {
          label: query.mode ? this.getModes().find(x => x.value === query.mode)!.label : this.getModes()[0].label,
          value: query.mode
        };

    return (
      <div className="gf-form">

        <Select<number>
            isMulti={false}
            isClearable={false}
            backspaceRemovesValue={false}
            onChange={i => this.onModeChange(i)}
            options={this.getModes()}
            isSearchable={true}
            value={selectedMode}
        />

      </div>
    );
  }
}
