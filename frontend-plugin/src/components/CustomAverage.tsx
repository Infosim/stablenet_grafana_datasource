/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Checkbox, Input, Select, LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const CustomAverage = props => (
  <div className="gf-form-inline" style={{ display: 'flex', alignItems: 'center' }}>
    <Checkbox css="" value={props.use} onChange={() => props.onChange[0]()} tabIndex={0} />

    <FormField
      label={'Custom Average Period'}
      labelWidth={11}
      tooltip={
        'Allows to define a custom average period. If disabled, Grafana will automatically compute a suiting average period.'
      }
      inputEl={
        <div className="gf-form-inline">
          <div className={'width-10'} tabIndex={0}>
            <Input
              css=""
              type="number"
              value={props.period}
              spellCheck={false}
              tabIndex={0}
              onChange={props.onChange[1]}
              disabled={!props.use}
            />
          </div>
          <div tabIndex={0}>
            <Select<number>
              options={props.getUnits()}
              value={props.unit}
              onChange={props.onChange[2]}
              className={'width-7'}
              isSearchable={true}
              menuPlacement={'bottom'}
            />
          </div>
        </div>
      }
    />
  </div>
);
