/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

export const MinMaxAvg = props => (
  <div className="gf-form" style={{ display: 'flex', alignItems: 'center' }}>
    <InlineFormLabel width={11}>Include Statistics:</InlineFormLabel>
    <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
      <Checkbox css="" value={props.values[0]} onChange={() => props.onChange('min')} tabIndex={0} label={'Min'} />
    </div>

    <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
      <Checkbox css="" value={props.values[1]} onChange={() => props.onChange('avg')} tabIndex={0} label={'Avg'} />
    </div>

    <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
      <Checkbox
        css=""
        style={{ paddingLeft: '2px', paddingRight: '2px' }}
        value={props.values[2]}
        onChange={() => props.onChange('max')}
        tabIndex={0}
        label={'Max'}
      />
    </div>
  </div>
);
