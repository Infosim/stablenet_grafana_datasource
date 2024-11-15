/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent } from 'react';
import { Checkbox, Input, Select, LegacyForms } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { LabelValue, Unit } from 'Types';

const { FormField } = LegacyForms;

interface Props {
  use: boolean;
  period: string;
  unit: number;
  onUseAverageChange: () => void;
  onUseCustomAverageChange: (event: ChangeEvent<HTMLInputElement>) => void;
  onAverageUnitChange: (value: SelectableValue<number>) => void;
}

const tooltip =
  'Allows to define a custom average period. If disabled, Grafana will automatically compute a suiting average period.';

const units: LabelValue[] = [
  { label: 'sec', value: Unit.SECONDS },
  { label: 'min', value: Unit.MINUTES },
  { label: 'hrs', value: Unit.HOURS },
  { label: 'days', value: Unit.DAYS },
];

export function CustomAverage({
  use,
  period,
  unit,
  onUseAverageChange,
  onUseCustomAverageChange,
  onAverageUnitChange,
}: Props): JSX.Element {
  return (
    <div className="gf-form-inline" style={{ display: 'flex', alignItems: 'center' }}>
      <Checkbox value={use} onChange={onUseAverageChange} tabIndex={0} />

      <FormField
        label={'Custom Average Period'}
        labelWidth={11}
        tooltip={tooltip}
        inputEl={
          <div className="gf-form-inline">
            <div className={'width-10'} tabIndex={0}>
              <Input type="number" value={period} spellCheck={false} tabIndex={0} onChange={onUseCustomAverageChange} disabled={!use} />
            </div>
            <div tabIndex={0}>
              <Select<number> options={units} value={unit} onChange={onAverageUnitChange} className={'width-7'} isSearchable={true} menuPlacement={'bottom'} />
            </div>
          </div>
        }
      />
    </div>
  );
}
