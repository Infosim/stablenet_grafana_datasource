import React, {PureComponent, ChangeEvent} from 'react';
import {SecretFormField, FormField} from '@grafana/ui';
import {DataSourcePluginOptionsEditorProps} from '@grafana/data';
import {StableNetConfigOptions, StableNetSecureJsonData} from './types';

interface Props extends DataSourcePluginOptionsEditorProps<StableNetConfigOptions> {
}

interface State {
}

export class StableNetConfigEditor extends PureComponent<Props, State> {
    onIpChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            ip: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            port: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            username: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        onOptionsChange({
            ...options,
            secureJsonData: {
                password: event.target.value,
            },
        });
    };

    onResetPassword = () => {
        const {onOptionsChange, options} = this.props;
        onOptionsChange({
            ...options,
            secureJsonFields: {
                ...options.secureJsonFields,
                password: false,
            },
            secureJsonData: {
                ...options.secureJsonData,
                password: '',
            },
        });
    };

    render() {
        const {options} = this.props;
        const {jsonData, secureJsonFields} = options;
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
                            value={jsonData.ip || ''}
                            placeholder="127.0.0.1"
                        />
                    </div>

                    <div className="gf-form">
                        <FormField
                            label="Port"
                            labelWidth={13}
                            inputWidth={17}
                            onChange={this.onPortChange}
                            value={jsonData.port || ''}
                            placeholder="5443"
                        />
                    </div>

                    <div className="gf-form">
                        <FormField
                            label="Username"
                            labelWidth={13}
                            inputWidth={17}
                            onChange={this.onUsernameChange}
                            value={jsonData.username || ''}
                            placeholder="infosim"
                        />
                    </div>

                    <div className="gf-form-inline">
                        <div className="gf-form">
                            <SecretFormField
                                isConfigured={(secureJsonFields && secureJsonFields.snpassword) as boolean}
                                value={secureJsonData.password || ''}
                                label="Password"
                                placeholder=""
                                labelWidth={13}
                                inputWidth={(secureJsonFields && secureJsonFields.snpassword) ? 16 : 17}
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
