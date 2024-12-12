/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { StableNetSecureJsonData } from '../types';

type Props = DataSourcePluginOptionsEditorProps<{}, StableNetSecureJsonData>;

const labelWidth = 15;

export const ConfigEditor = ({ options, onOptionsChange }: Props): JSX.Element => {
  const { url, user, secureJsonFields, secureJsonData } = options;

  const onUrlChange = (event: ChangeEvent<HTMLInputElement>) =>
    onOptionsChange({
      ...options,
      url: event.target.value,
    });

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) =>
    onOptionsChange({
      ...options,
      user: event.target.value,
    });

  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) =>
    onOptionsChange({
      ...options,
      secureJsonData: { password: event.target.value },
    });

  const onResetPassword = () =>
    onOptionsChange({
      ...options,
      secureJsonFields: { password: false },
      secureJsonData: { password: undefined },
    });

  return (
    <>
      <InlineField
        label="URL"
        labelWidth={labelWidth}
        tooltip="IP & PORT of your StableNet® Server installation"
        interactive
      >
        <Input id="stablenet-ip" value={url} placeholder="https://127.0.0.1:5443" onChange={onUrlChange} required />
      </InlineField>

      <InlineField
        label="Username"
        labelWidth={labelWidth}
        tooltip="Username of your StableNet® Server installation"
        interactive
      >
        <Input id="stablenet-username" value={user} placeholder="infosim" onChange={onUsernameChange} required />
      </InlineField>

      <InlineField
        label="Password"
        labelWidth={labelWidth}
        tooltip="Username of your StableNet® Server installation"
        interactive
      >
        <SecretInput
          id="stablenet-password"
          value={secureJsonData?.password}
          isConfigured={secureJsonFields.password}
          placeholder="Password"
          onChange={onPasswordChange}
          onReset={onResetPassword}
          required
        />
      </InlineField>
    </>
  );
};
