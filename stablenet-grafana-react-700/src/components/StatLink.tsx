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

export const StatLink = props => (
  <div className="gf-form-inline">
    <div className={'gf-form'}>
      <FormField
        label={'Link:'}
        labelWidth={11}
        inputWidth={19}
        tooltip={
          'Copy a link from the StableNet®-Analyzer. Due to technical limitations, measurements other than template measurements ' +
          '(e.g. ping and interface measurements) are only partly supported.'
        }
        value={props.link}
        onChange={props.onChange}
        spellCheck={false}
        tabIndex={0}
      />
    </div>
  </div>
);
