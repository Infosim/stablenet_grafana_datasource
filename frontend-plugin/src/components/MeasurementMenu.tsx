/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent } from 'react';
import { Select, LegacyForms } from '@grafana/ui';
import { LabelValue } from 'Types';
import { SelectableValue } from '@grafana/data';

const { FormField } = LegacyForms;

interface Props {
  hasMoreMeasurements: boolean;
  selected: LabelValue;
  measurements: LabelValue[];
  filter: string;
  disabled: boolean;
  onChange: (value: SelectableValue<number>) => void;
  onFilterChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

const moreMeasurementsTooltip = 'There are more measurements available, but only the first 100 are displayed. Use a stricter search to reduce the number of shown measurements.';

const filterTooltip = 'The dropdown menu on the left only shows at most 100 measurements. Use this text field to query measurements that are not shown on the left, or to search for specific measurements.';

export function MeasurementMenu({ measurements, hasMoreMeasurements, selected, filter, disabled, onChange, onFilterChange }: Props): JSX.Element {

  const inputElement = (
    <div tabIndex={0}>
      <Select<number>
        options={measurements}
        value={selected}
        onChange={onChange}
        className={'width-19'}
        menuPlacement={'bottom'}
        noOptionsMessage={`No measurements match this search.`}
        placeholder={'none'}
        isSearchable={false} />
    </div>
  );

  return (
    <div className="gf-form">
      <div style={{ marginRight: 4 }}>
        <FormField label={'Measurement:'} labelWidth={11} tooltip={hasMoreMeasurements ? moreMeasurementsTooltip : ''} inputEl={inputElement} />
      </div>
      <FormField
        label={'Measurement Filter:'}
        labelWidth={11}
        inputWidth={19}
        tooltip={filterTooltip}
        value={filter}
        onChange={onFilterChange}
        spellCheck={false}
        placeholder={'no filter'}
        tabIndex={0}
        disabled={disabled} />
    </div>
  );
};
