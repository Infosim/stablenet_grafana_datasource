import React from 'react';
import { InlineFormLabel, Input } from '@grafana/ui';

export const StatLink = props => (
  <div className="gf-form-inline">
    {/** Statistic Link mode */}
    <div className="gf-form">
      <InlineFormLabel
        width={11}
        tooltip="Copy a link from the StableNet®-Analyzer. Experimental: At the current version, only links containing exactly one measurement are supported."
      >
        Link:
      </InlineFormLabel>

      <div className={'width-19'}>
        <Input type="text" value={props.link} spellCheck={false} tabIndex={0} onChange={props.onChange} />
      </div>
    </div>
  </div>
);
