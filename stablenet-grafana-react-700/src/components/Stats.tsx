import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

export const Stats = props => (
  <div className="gf-form">
    <div style={!props.mode ? { marginLeft: '480px' } : {}}>
      <InlineFormLabel width={11} tooltip="Select the statistics you want to display.">
        Include Statistics:
      </InlineFormLabel>
    </div>

    <div className="gf-form">
      <Checkbox value={props.values[0]} onChange={() => props.onChange('min')} tabIndex={0} />
      <InlineFormLabel width={5}>Min</InlineFormLabel>

      <Checkbox value={props.values[1]} onChange={() => props.onChange('avg')} tabIndex={0} />
      <InlineFormLabel width={5}>Avg</InlineFormLabel>

      <Checkbox value={props.values[2]} onChange={() => props.onChange('max')} tabIndex={0} />
      <InlineFormLabel width={5}>Max</InlineFormLabel>
    </div>
  </div>
);
