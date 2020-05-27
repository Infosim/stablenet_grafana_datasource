import React from 'react';
import { LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const StatLink = props => (
  <div className="gf-form-inline">
    <div className={'gf-form'}>
      <FormField
        label={'Link:'}
        labelWidth={11}
        inputWidth={19}
        tooltip={
          'Copy a link from the StableNet®-Analyzer. Experimental: At the current version, only links containing exactly one measurement are supported.'
        }
        value={props.link}
        onChange={props.onChange}
        spellCheck={false}
        tabIndex={0}
      />
    </div>
  </div>
);
