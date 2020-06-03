/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { StableNetConfigOptions, StableNetSecureJsonData } from './Types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<StableNetConfigOptions> {}

export class ConfigEditor extends PureComponent<Props> {
  onIpChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      snip: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      snport: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      snusername: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        snpassword: event.target.value,
      },
    });
  };

  onResetPassword = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        snpassword: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        snpassword: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as StableNetSecureJsonData;

    return (
      <div>
        <h3 className="page-heading">StableNet® Configuration</h3>

        <div className="gf-form-group">
          <div className="gf-form">
            <FormField
              label="StableNet® Server"
              labelWidth={13}
              inputWidth={17}
              onChange={this.onIpChange}
              value={jsonData.snip || ''}
              placeholder="127.0.0.1"
            />
          </div>

          <div className="gf-form">
            <FormField
              label="Port"
              labelWidth={13}
              inputWidth={17}
              onChange={this.onPortChange}
              value={jsonData.snport || ''}
              placeholder="5443"
            />
          </div>

          <div className="gf-form">
            <FormField
              label="Username"
              labelWidth={13}
              inputWidth={17}
              onChange={this.onUsernameChange}
              value={jsonData.snusername || ''}
              placeholder="infosim"
            />
          </div>

          <div className="gf-form-inline">
            <div className="gf-form">
              <SecretFormField
                isConfigured={(secureJsonFields && secureJsonFields.snpassword) as boolean}
                value={secureJsonData.snpassword || ''}
                label="Password"
                placeholder=""
                labelWidth={13}
                inputWidth={secureJsonFields && secureJsonFields.snpassword ? 16 : 17}
                onReset={this.onResetPassword}
                onChange={this.onPasswordChange}
              />
            </div>
          </div>
        </div>
      </div>
    );
  }
}
