import React from 'react';
import { InlineFormLabel, Input } from '@grafana/ui';

export const MetricPrefix = props => (
  <div className="gf-form">
    <InlineFormLabel
      width={11}
      tooltip="The input of this field will be added as a prefix to the metrics' names on the chart."
    >
      Metric prefix:
    </InlineFormLabel>
    <div className="width-19" style={props.space}>
      <Input type="text" value={props.value} spellCheck={false} tabIndex={0} onChange={props.onChange} />
    </div>
  </div>
);
