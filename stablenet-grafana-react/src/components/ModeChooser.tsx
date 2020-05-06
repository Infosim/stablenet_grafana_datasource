import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const ModeChooser = props => (
  <div className="gf-form-inline">
    <div className="gf-form">
      <FormLabel width={11} tooltip="Allows switching between Measurement mode and Statistic Link mode.">
        Query Mode:
      </FormLabel>

      <div tabIndex={0} style={props.space}>
        <Forms.Select<number> options={props.options()} value={props.mode} onChange={props.onChange} width={10} isSearchable={true} />
      </div>
    </div>
  </div>
);
