/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Select, LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const MeasurementMenu = props => {
  return (
    <div className="gf-form">
      <FormField
        label={'Measurement Filter:'}
        labelWidth={11}
        inputWidth={19}
        tooltip={'The dropdown menu on the right will only include entries containing the input of this field.'}
        value={props.filter}
        onChange={props.filterChange}
        spellCheck={false}
        placeholder={'no filter'}
        tabIndex={0}
      />
      <FormField
        label={'Measurement:'}
        labelWidth={11}
        tooltip={
          props.more
            ? `There are more measurements available, but only the first 100 are displayed.
              Use a stricter search to reduce the number of shown measurements.`
            : ''
        }
        inputEl={
          <div tabIndex={0}>
            <Select<number>
              options={props.get}
              value={props.selected}
              onChange={props.menuChange}
              className={'width-19'}
              menuPlacement={'bottom'}
              noOptionsMessage={`No measurements match this search.`}
              placeholder={'none'}
              isSearchable={false}
            />
          </div>
        }
      />
    </div>
  );
};
