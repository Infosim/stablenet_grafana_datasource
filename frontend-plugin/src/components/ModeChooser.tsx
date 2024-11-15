/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React from 'react';
import { Select, LegacyForms } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { Mode } from 'Types';

const { FormField } = LegacyForms;

interface Props {
  selectedMode: number;
  onChange: (value: SelectableValue<number>) => void;
}

const modes: Array<SelectableValue<number>> = [
  { label: 'Measurement', value: Mode.MEASUREMENT },
  { label: 'Statistic Link', value: Mode.STATISTIC_LINK },
];

const tooltip = 'Allows switching between Measurement mode and Statistic Link mode.';

export function ModeChooser({ selectedMode, onChange }: Props): JSX.Element {
  const inputElement = (
    <div tabIndex={0}>
      <Select<number>
        value={selectedMode}
        options={modes}
        onChange={onChange}
        className={'width-10'}
        menuPlacement={'bottom'}
        isSearchable={true}
      />
    </div>
  );

  return (
    <div className="gf-form-inline">
      <div className="gf-form">
        <FormField label={'Query Mode:'} labelWidth={11} tooltip={tooltip} inputEl={inputElement} />
      </div>
    </div>
  );
}
