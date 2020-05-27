import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

const checkboxOuter = {
  width: '32px',
  height: '32px',
  borderStyle: 'solid',
  borderWidth: '1px',
  borderColor: '#2c3235',
  borderRadius: '3px',
  marginBottom: '5px',
  marginRight: '2px',
} as React.CSSProperties;

const checkboxInner = {
  paddingLeft: '7px',
  marginTop: '-5px',
} as React.CSSProperties;

export const Stats = props => (
  <div className="gf-form">
    <div style={!props.mode ? { marginLeft: '480px' } : {}}>
      <InlineFormLabel width={11} tooltip="Select the statistics you want to display.">
        Include Statistics:
      </InlineFormLabel>
    </div>

    <div className="gf-form">
      <div style={checkboxOuter}>
        <div style={checkboxInner}>
          <Checkbox value={props.values[0]} onChange={() => props.onChange('min')} tabIndex={0} />
        </div>
      </div>
      <InlineFormLabel width={4}>Min</InlineFormLabel>

      <div style={checkboxOuter}>
        <div style={checkboxInner}>
          <Checkbox value={props.values[1]} onChange={() => props.onChange('avg')} tabIndex={0} />
        </div>
      </div>
      <InlineFormLabel width={4}>Avg</InlineFormLabel>

      <div style={checkboxOuter}>
        <div style={checkboxInner}>
          <Checkbox value={props.values[2]} onChange={() => props.onChange('max')} tabIndex={0} />
        </div>
      </div>
      <InlineFormLabel width={4}>Max</InlineFormLabel>
    </div>
  </div>
);
