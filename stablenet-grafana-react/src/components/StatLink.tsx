import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const StatLink = props => (
  <div className="gf-form-inline">
    {/** Statistic Link mode */}
    <div className="gf-form">
      <FormLabel
        width={11}
        tooltip="Copy a link from the StableNet®-Analyzer. Experimental: At the current version, only links containing exactly one measurement are supported."
      >
        Link:
      </FormLabel>

      <div className="width-19">
        <Forms.Input type="text" value={props.link} spellCheck={false} tabIndex={0} onChange={props.onChange} />
      </div>
    </div>
  </div>
);
