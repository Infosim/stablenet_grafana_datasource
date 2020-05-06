import { FormLabel, Forms } from '@grafana/ui';
import React from 'react';

export const DropdownMenu = props => {
  const plural = props.name.toLowerCase() + 's';
  return (
    <div className="gf-form-inline">
      <div className="gf-form">
        <FormLabel width={11}>{props.name + ':'}</FormLabel>

        <div tabIndex={0} style={props.space}>
          <Forms.AsyncSelect<number>
            loadOptions={props.get}
            value={props.selected}
            onChange={props.onChange}
            defaultOptions={true}
            noOptionsMessage={`No ${plural} match this search.`}
            loadingMessage={`Fetching ${plural}...`}
            width={19}
            placeholder="none"
            isSearchable={true}
          />
        </div>
      </div>
      {props.more ? (
        <div className="gf-form">
          <FormLabel
            children={{}}
            tooltip={`There are more ${plural} available, but only the first 100 are displayed.
                                                Use a stricter search to reduce the number of shown ${plural}.`}
          />
        </div>
      ) : null}
    </div>
  );
};
