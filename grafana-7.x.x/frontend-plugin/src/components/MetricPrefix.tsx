/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const MetricPrefix = props => (
  <div className="gf-form">
    <FormField
      label={'Metric Prefix:'}
      labelWidth={11}
      inputWidth={19}
      tooltip={
        "The input of this field will be added as a prefix to the metrics' names on the chart. " +
        'This only applies if there are two or data series are shown in the chart.'
      }
      value={props.value}
      onChange={props.onChange}
      spellCheck={false}
      tabIndex={0}
    />
  </div>
);
