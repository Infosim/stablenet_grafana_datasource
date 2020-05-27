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

export const DropdownMenu = props => {
  const plural = props.name.toLowerCase() + 's';
  return (
    <div className="gf-form-inline">
      <div className="gf-form">
        <FormField
          label={props.name + ':'}
          labelWidth={11}
          tooltip={
            props.more
              ? `There are more ${plural} available, but only the first 100 are displayed.
                                                Use a stricter search to reduce the number of shown ${plural}.`
              : ''
          }
          inputEl={
            <div tabIndex={0}>
              <AsyncSelect<number>
                loadOptions={props.get}
                value={props.selected}
                onChange={props.onChange}
                defaultOptions={true}
                noOptionsMessage={`No ${plural} match this search.`}
                loadingMessage={`Fetching ${plural}...`}
                className={'width-19'}
                placeholder={'none'}
                menuPlacement={'bottom'}
                isSearchable={true}
              />
            </div>
          }
        />
      </div>
    </div>
  );
};
