/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent } from 'react';
import { LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

interface IProps {
  value: string;
  onChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

const tooltip = "The input of this field will be added as a prefix to the metrics' names on the chart. This only applies if two or more data series are shown in the chart.";

export function MetricPrefix({ value, onChange }: IProps): JSX.Element {
  return (
    <div className="gf-form">
      <FormField label={'Metric Prefix:'} labelWidth={11} inputWidth={19} tooltip={tooltip} value={value} onChange={onChange} spellCheck={false} tabIndex={0} />
    </div>
  );
}
