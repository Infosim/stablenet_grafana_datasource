/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent } from 'react';
import { Input, InlineFormLabel } from '@grafana/ui';

interface Props {
  link: string;
  onChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

const tooltip =
  'Copy a link from the StableNet®-Analyzer. Due to technical limitations, measurements other than template measurements (e.g. ping and interface measurements) are only partly supported.';

export function StatLink({ link, onChange }: Props): JSX.Element {
  return (
    <div className="gf-form-inline">
      <div className={'gf-form'} style={{ width: '100%' }}>
        <InlineFormLabel width={11} tooltip={tooltip}>
          Link:
        </InlineFormLabel>

        <div style={{ width: '100%' }}>
          <Input type={'text'} value={link} onChange={onChange} spellCheck={false} tabIndex={0} />
        </div>
      </div>
    </div>
  );
}
