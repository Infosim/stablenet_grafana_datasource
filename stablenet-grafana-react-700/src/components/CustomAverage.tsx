import React from 'react';
import { Checkbox, InlineFormLabel, Input, Select } from '@grafana/ui';

const checkboxOuter = {
  width: '33px',
  borderStyle: 'solid',
  borderWidth: '1px',
  borderColor: '#2c3235',
  borderRadius: '3px',
  marginBottom: '5px',
} as React.CSSProperties;

const checkboxInner = {
  paddingLeft: '7.5px',
  marginTop: '-5px',
} as React.CSSProperties;

export const CustomAverage = props => (
  <div className="gf-form-inline">
    <div style={checkboxOuter}>
      <div style={checkboxInner}>
        <Checkbox value={props.use} onChange={() => props.onChange[0]()} tabIndex={0} />
      </div>
    </div>
    <InlineFormLabel
      width={11}
      tooltip="Allows defining a custom average period. If disabled, Grafana will automatically compute a suiting average period."
    >
      Average Period:
    </InlineFormLabel>
    <div className={'width-10'} style={props.space}>
      <Input
        type="number"
        value={props.period}
        spellCheck={false}
        tabIndex={0}
        onChange={props.onChange[1]}
        disabled={!props.use}
      />
    </div>
    <div className="gf-form">
      <div tabIndex={0} style={props.space}>
        <Select<number>
          options={props.getUnits()}
          value={props.unit}
          onChange={props.onChange[2]}
          className={'width-7'}
          isSearchable={true}
        />
      </div>
    </div>
  </div>
);
