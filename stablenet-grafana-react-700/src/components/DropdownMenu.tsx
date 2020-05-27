import React from 'react';
import { AsyncSelect, InlineFormLabel, LegacyForms } from '@grafana/ui';

const { FormField } = LegacyForms;

export const DropdownMenu = props => {
  const plural = props.name.toLowerCase() + 's';
  return (
    <div className="gf-form-inline">
      <div className="gf-form">
        <FormField
          label={props.name + ':'}
          labelWidth={11}
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
      {props.more ? (
        <div className="gf-form">
          <InlineFormLabel
            children={{}}
            tooltip={`There are more ${plural} available, but only the first 100 are displayed.
                                                Use a stricter search to reduce the number of shown ${plural}.`}
          />
        </div>
      ) : null}
    </div>
  );
};
