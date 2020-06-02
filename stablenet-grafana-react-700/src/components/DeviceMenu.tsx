/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { AsyncSelect, LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const DeviceMenu = props => {
  return (
    <div className="gf-form">
      <FormField
        label={'Device:'}
        labelWidth={11}
        tooltip={
          props.more
            ? `There are more devices available, but only the first 100 are displayed.
                                                Use a stricter search to reduce the number of shown devices.`
            : ''
        }
        inputEl={
          <div tabIndex={0}>
            <AsyncSelect<number>
              loadOptions={props.get}
              value={props.selected}
              onChange={props.onChange}
              defaultOptions={true}
              noOptionsMessage={`No devices match this search.`}
              loadingMessage={`Fetching devices...`}
              className={'width-19'}
              placeholder={'none'}
              menuPlacement={'bottom'}
              isSearchable={true}
              backspaceRemovesValue={true}
            />
          </div>
        }
      />
    </div>
  );
};
