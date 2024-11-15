/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Checkbox, InlineFormLabel } from '@grafana/ui';

interface Props {
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
  onChange: (value: 'min' | 'avg' | 'max') => void;
}

export function MinMaxAvg({ includeMinStats, includeAvgStats, includeMaxStats, onChange }: Props): JSX.Element {
  return (
    <div className="gf-form" style={{ display: 'flex', alignItems: 'center' }}>
      <InlineFormLabel width={11}>Include Statistics:</InlineFormLabel>

      <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
        <Checkbox css="" value={includeMinStats} onChange={() => onChange('min')} tabIndex={0} label={'Min'} />
      </div>

      <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
        <Checkbox css="" value={includeAvgStats} onChange={() => onChange('avg')} tabIndex={0} label={'Avg'} />
      </div>

      <div style={{ paddingLeft: '2px', paddingRight: '2px' }}>
        <Checkbox css="" value={includeMaxStats} onChange={() => onChange('max')} tabIndex={0} label={'Max'} />
      </div>
    </div>
  );
}
