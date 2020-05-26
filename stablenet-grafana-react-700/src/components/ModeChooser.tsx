import React from 'react';
import { InlineFormLabel, Select } from '@grafana/ui';

export const ModeChooser = props => (
  <div className="gf-form-inline">
    <div className="gf-form">
      <InlineFormLabel width={11} tooltip="Allows switching between Measurement mode and Statistic Link mode.">
        Query Mode:
      </InlineFormLabel>

      <div tabIndex={0} style={props.space}>
        <Select<number>
          options={props.options()}
          value={props.mode}
          onChange={props.onChange}
          className={'width-10'}
          isSearchable={true}
        />
      </div>
    </div>
  </div>
);
