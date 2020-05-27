/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

const checkboxOuter = {
  width: '32px',
  height: '32px',
  borderStyle: 'solid',
  borderWidth: '1px',
  borderColor: '#2c3235',
  borderRadius: '3px',
} as React.CSSProperties;

const checkboxInner = {
  paddingLeft: '7.5px',
  marginTop: '-5.5px',
} as React.CSSProperties;

export const Metric = props => (
  <div className="gf-form">
    <div style={checkboxOuter}>
      <div style={checkboxInner}>
        <Checkbox value={props.value} onChange={props.onChange} />
      </div>
    </div>
    <div style={props.singleMetric}>
      <InlineFormLabel width={17}>{props.text}</InlineFormLabel>
    </div>
  </div>
);
