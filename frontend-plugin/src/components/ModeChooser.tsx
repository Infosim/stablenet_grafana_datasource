/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Select, LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const ModeChooser = props => (
  <div className="gf-form-inline">
    <div className="gf-form">
      <FormField
        label={'Query Mode:'}
        labelWidth={11}
        tooltip={'Allows switching between Measurement mode and Statistic Link mode.'}
        inputEl={
          <div tabIndex={0}>
            <Select<number>
              options={props.options()}
              value={props.mode}
              onChange={props.onChange}
              className={'width-10'}
              menuPlacement={'bottom'}
              isSearchable={true}
            />
          </div>
        }
      />
    </div>
  </div>
);
