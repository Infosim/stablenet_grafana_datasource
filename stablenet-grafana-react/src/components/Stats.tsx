import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const Stats = props => (
  <div className="gf-form">
    <div style={!props.mode ? { marginLeft: '415px' } : {}}>
      <FormLabel width={11} tooltip="Select the statistics you want to display.">
        Include Statistics:
      </FormLabel>
    </div>

    <div className="gf-form">
      <Forms.Checkbox value={props.values[0]} onChange={() => props.onChange('min')} tabIndex={0} />
      <FormLabel width={5}>Min</FormLabel>

      <Forms.Checkbox value={props.values[1]} onChange={() => props.onChange('avg')} tabIndex={0} />
      <FormLabel width={5}>Avg</FormLabel>

      <Forms.Checkbox value={props.values[2]} onChange={() => props.onChange('max')} tabIndex={0} />
      <FormLabel width={5}>Max</FormLabel>
    </div>
  </div>
);
