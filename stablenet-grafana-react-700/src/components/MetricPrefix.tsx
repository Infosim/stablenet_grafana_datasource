import React from 'react';
import { LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const MetricPrefix = props => (
  <div className="gf-form">
    <FormField
      label={'Metric Prefix:'}
      labelWidth={11}
      inputWidth={19}
      tooltip={"The input of this field will be added as a prefix to the metrics' names on the chart."}
      value={props.value}
      onChange={props.onChange}
      spellCheck={false}
      tabIndex={0}
    />
  </div>
);
