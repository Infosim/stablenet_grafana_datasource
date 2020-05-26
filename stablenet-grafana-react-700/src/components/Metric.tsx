import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

export const Metric = props => (
  <div className="gf-form">
    <Checkbox value={props.value} onChange={props.onChange} size={11} />
    <div style={props.singleMetric}>
      <InlineFormLabel width={17}>{props.text}</InlineFormLabel>
    </div>
  </div>
);
