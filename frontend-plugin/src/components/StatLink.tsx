/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Input, InlineFormLabel } from '@grafana/ui';

export const StatLink = props => (
  <div className="gf-form-inline">
    <div className={'gf-form'} style={{ width: '100%' } as React.CSSProperties}>
      <InlineFormLabel
        width={11}
        tooltip={
          'Copy a link from the StableNet®-Analyzer. Due to technical limitations, measurements other than template measurements ' +
          '(e.g. ping and interface measurements) are only partly supported.'
        }
      >
        Link:
      </InlineFormLabel>
      <div style={{ width: '100%' } as React.CSSProperties}>
        <Input type={'text'} value={props.link} onChange={props.onChange} spellCheck={false} tabIndex={0} />
      </div>
    </div>
  </div>
);
