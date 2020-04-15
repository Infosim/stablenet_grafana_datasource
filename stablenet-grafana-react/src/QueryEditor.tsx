import React, {PureComponent, ChangeEvent} from 'react';
import {FormLabel, Forms} from '@grafana/ui';
import {QueryEditorProps, SelectableValue} from '@grafana/data';
import {DataSource} from './DataSource';
import {Mode, StableNetConfigOptions, Unit} from './types';
import {Target} from "./query_interfaces";
import "./css/query-editor.css";
import {LabelValue} from "./returnTypes";

type Props = QueryEditorProps<DataSource, Target, StableNetConfigOptions>;

interface State {
}

export class QueryEditor extends PureComponent<Props, State> {
    onComponentDidMount() {
    }

    getModes(): LabelValue[] {
        return [
            {label: 'Measurement', value: Mode.MEASUREMENT},
            {label: 'Statistic Link', value: Mode.STATISTIC_LINK},
        ];
    }

    getUnits(): LabelValue[] {
        return [
            {label: 'sec', value: Unit.SECONDS},
            {label: 'min', value: Unit.MINUTES},
            {label: 'hrs', value: Unit.HOURS},
            {label: 'days', value: Unit.DAYS},
        ];
    }

    onModeChange = (v: SelectableValue<number>) => {
        const {onChange, query} = this.props;
        onChange({
            ...query,
            mode: v.value!,
            includeMaxStats: false,
            includeAvgStats: true,
            includeMinStats: false
        });
    };

    onStatisticLinkChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({...query, statisticLink: event.target.value});
        onRunQuery(); // executes the query
    };

    getDevices = (v: string) => {
        const {query, onChange} = this.props;
        return (this.props.datasource as DataSource)
            .queryDevices(v, query.refId)
            .then(r => {
                onChange({...query, moreDevices: r.hasMore});
                return r.data;
            });
    };

    onDeviceChange = (v: SelectableValue<number>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            selectedDevice: {label: v.label!, value: v.value!},
            measurementQuery: '',
            selectedMeasurement: {label: "", value: -1},
            metricPrefix: '',
            metrics: [],
            chosenMetrics: {}
        });
        onRunQuery();
    };

    getMeasurements = (v: string) => {
        const {query, onChange} = this.props;
        return (this.props.datasource as DataSource)
            .findMeasurementsForDevice(query.selectedDevice ? query.selectedDevice.value : -1, v, query.refId)
            .then(r => {
                onChange({...query, moreMeasurements: r.hasMore});
                return r.data;
            });
    };

    onMeasurementChange = (v: SelectableValue<number>) => {
        const {onChange, query, onRunQuery} = this.props;
        (this.props.datasource as DataSource)
            .findMetricsForMeasurement(v.value!, query.refId)
            .then(r => {
                onChange({
                    ...query, metrics: r,
                    chosenMetrics: {},
                    metricPrefix: v.label!,
                    selectedMeasurement: {label: v.label!, value: v.value!}
                });
            })
            .then(() => onRunQuery());
    };



    onMoreChange = () => {
        const {onChange, query} = this.props;
        onChange({...query, moreDevices: !this.props.query.moreDevices});
    };

    render() {
        const query = this.props.query;
        const space = {marginTop: "2px"} as React.CSSProperties;

        return (
            <div>

                <div className="gf-form-inline">
                    <div className="gf-form">
                        <FormLabel
                            width={11}
                            tooltip="Allows switching between Measurement mode and Statistic Link mode."
                        >Query Mode:
                        </FormLabel>

                        <div tabIndex={0}
                             style={space}>
                            <Forms.Select<number>
                                options={this.getModes()}
                                value={query.mode || Mode.MEASUREMENT}
                                onChange={this.onModeChange}
                                width={10}
                                isSearchable={true}
                                allowCustomValue={true}
                            />
                        </div>
                    </div>
                </div>

                {
                    !!query.mode ?
                        <div className="gf-form-inline">{/** Statistic Link mode */}
                            <div className="gf-form gf-form--grow">
                                <FormLabel
                                    width={11}
                                    tooltip="Copy a link from the StableNet®-Analyzer. Experimental: At the current version, only links containing exactly one measurement are supported."
                                >Link:
                                </FormLabel>

                                <Forms.Input
                                    type="text"
                                    value={query.statisticLink || ''}
                                    spellCheck={false}
                                    tabIndex={0}
                                    onChange={this.onStatisticLinkChange}
                                />
                            </div>
                        </div>
                        :
                        <div>{/** Measurement mode */}
                            <div className="gf-form-inline">
                                <div className="gf-form">
                                    <FormLabel width={11}>
                                        Device:
                                    </FormLabel>

                                    <div tabIndex={0}
                                         style={space}>
                                        <Forms.AsyncSelect<number>
                                            loadOptions={this.getDevices}
                                            value={query.selectedDevice}
                                            onChange={this.onDeviceChange}
                                            defaultOptions={true}
                                            noOptionsMessage={'No devices match this search.'}
                                            loadingMessage="Fetching devices..."
                                            width={19}
                                            placeholder="none"
                                            isSearchable={true}
                                        />
                                    </div>
                                </div>
                                {
                                    query.moreDevices ?
                                        <div className="gf-form">
                                            <FormLabel
                                                children={{}}
                                                tooltip="There are more devices available, but only the first 100 are displayed.
                                                Use a stricter search to reduce the number of shown devices."/>
                                        </div>
                                        :
                                        null
                                }
                            </div>

                            <div className="gf-form-inline">
                                <div className="gf-form">
                                    <FormLabel width={11}>
                                        Measurement:
                                    </FormLabel>

                                    <div tabIndex={0}
                                         style={space}>
                                        <Forms.AsyncSelect<number>
                                            loadOptions={this.getMeasurements}
                                            value={query.selectedMeasurement}
                                            onChange={this.onMeasurementChange}
                                            defaultOptions={true}
                                            noOptionsMessage={'No measurements match this search.'}
                                            loadingMessage="Fetching measurements..."
                                            width={19}
                                            placeholder="none"
                                            isSearchable={true}
                                        />
                                    </div>
                                </div>
                                {
                                    query.moreMeasurements ?
                                        <div className="gf-form">
                                            <FormLabel
                                                children={{}}
                                                tooltip="There are more measurements available, but only the first 100 are displayed.
                                                Specify a stricter search to reduce the number of shown devices."/>
                                        </div>
                                        :
                                        null
                                }
                            </div>
                        </div>
                }

                <div>
                    {
                        !!query.selectedMeasurement || !!query.mode ?
                            <div>

                            </div>
                            :
                            null
                    }
                </div>

                <div className="gf-form">

                    <Forms.Checkbox
                        value={query.moreDevices}
                        onChange={this.onMoreChange}
                        size={8}
                    />

                </div>
            </div>
        );
    }

}
