/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent, memo } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { StableNetConfigOptions, StableNetSecureJsonData } from './Types';

const { SecretFormField, FormField } = LegacyForms;

export const ConfigEditor = memo(({ options, onOptionsChange }: DataSourcePluginOptionsEditorProps<StableNetConfigOptions, StableNetSecureJsonData>): JSX.Element => {

  const { jsonData, secureJsonFields, secureJsonData } = options;

  const onIpChange = (event: ChangeEvent<HTMLInputElement>) => onOptionsChange({ ...options, jsonData: { ...jsonData, snip: event.target.value } });

  const onPortChange = (event: ChangeEvent<HTMLInputElement>) => onOptionsChange({ ...options, jsonData: { ...jsonData, snport: event.target.value } });

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => onOptionsChange({ ...options, jsonData: { ...jsonData, snusername: event.target.value } });

  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => onOptionsChange({ ...options, secureJsonData: { snpassword: event.target.value } });

  const onResetPassword = () => onOptionsChange({
    ...options,
    secureJsonFields: { ...secureJsonFields, snpassword: false },
    secureJsonData: { ...secureJsonData, snpassword: '' },
  });

  return (
    <div>
      <h3 className="page-heading">StableNet® Configuration</h3>

      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="StableNet® Server"
            labelWidth={13}
            inputWidth={17}
            onChange={onIpChange}
            value={jsonData.snip || ''}
            placeholder="127.0.0.1" />
        </div>

        <div className="gf-form">
          <FormField
            label="Port"
            labelWidth={13}
            inputWidth={17}
            onChange={onPortChange}
            value={jsonData.snport || ''}
            placeholder="5443" />
        </div>

        <div className="gf-form">
          <FormField
            label="Username"
            labelWidth={13}
            inputWidth={17}
            onChange={onUsernameChange}
            value={jsonData.snusername || ''}
            placeholder="infosim" />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={secureJsonFields && secureJsonFields.snpassword}
              value={secureJsonData?.snpassword || ''}
              label="Password"
              placeholder=""
              labelWidth={13}
              inputWidth={secureJsonFields && secureJsonFields.snpassword ? 16 : 17}
              onReset={onResetPassword}
              onChange={onPasswordChange} />
          </div>
        </div>
      </div>
    </div>
  );
});
