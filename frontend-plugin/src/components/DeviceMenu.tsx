/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { AsyncSelect, LegacyForms } from '@grafana/ui';
import { LabelValue } from 'Types';
import { SelectableValue } from '@grafana/data';

const { FormField } = LegacyForms;

interface Props {
  selectedDevice: LabelValue;
  hasMoreDevices: boolean;
  get: (value: string) => Promise<LabelValue[]>;
  onChange: (value: SelectableValue<number>) => void;
}

const moreDevicesTooltip =
  'There are more devices available, but only the first 100 are displayed. Use a stricter search to reduce the number of shown devices.';

export function DeviceMenu({ selectedDevice, hasMoreDevices, get, onChange }: Props): JSX.Element {
  return (
    <div className="gf-form">
      <FormField
        label={'Device:'}
        labelWidth={11}
        tooltip={hasMoreDevices ? moreDevicesTooltip : ''}
        inputEl={
          <div tabIndex={0}>
            <AsyncSelect<number>
              value={selectedDevice}
              loadOptions={get}
              onChange={onChange}
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
}
