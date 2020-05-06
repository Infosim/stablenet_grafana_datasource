import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const MetricPrefix = props => (
  <div className="gf-form">
    <FormLabel width={11} tooltip="The input of this field will be added as a prefix to the metrics' names on the chart.">
      Metric prefix:
    </FormLabel>
    <div className="width-19" style={props.space}>
      <Forms.Input type="text" value={props.value} spellCheck={false} tabIndex={0} onChange={props.onChange} />
    </div>
  </div>
);
