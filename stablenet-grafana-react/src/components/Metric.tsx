import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const Metric = props => (
  <div className="gf-form">
    <Forms.Checkbox value={props.value} onChange={props.onChange} size={11} />
    <div style={props.singleMetric}>
      <FormLabel width={17}>{props.text}</FormLabel>
    </div>
  </div>
);
