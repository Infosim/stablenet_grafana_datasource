import React, {PureComponent, ChangeEvent} from 'react';
import {SecretFormField, FormField} from '@grafana/ui';
import {DataSourcePluginOptionsEditorProps} from '@grafana/data';
import {MyDataSourceOptions, MySecureJsonData} from './types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {
}

interface State {
}

export class ConfigEditor extends PureComponent<Props, State> {
    onIpChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            snip: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            snport: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        const jsonData = {
            ...options.jsonData,
            snusername: event.target.value,
        };
        onOptionsChange({...options, jsonData});
    };

    // Secure field (only sent to the backend)
    onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onOptionsChange, options} = this.props;
        onOptionsChange({
            ...options,
            secureJsonData: {
                apiKey: event.target.value,
            },
        });
    };

    onResetAPIKey = () => {
        const {onOptionsChange, options} = this.props;
        onOptionsChange({
            ...options,
            secureJsonFields: {
                ...options.secureJsonFields,
                apiKey: false,
            },
            secureJsonData: {
                ...options.secureJsonData,
                apiKey: '',
            },
        });
    };

    render() {
        console.log(this.props);
        const {options} = this.props;
        const {jsonData, secureJsonFields} = options;
        const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

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
                                isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                                value={secureJsonData.apiKey || ''}
                                label="API Key"
                                placeholder="secure json field (backend only)"
                                labelWidth={6}
                                inputWidth={20}
                                onReset={this.onResetAPIKey}
                                onChange={this.onAPIKeyChange}
                            />
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
