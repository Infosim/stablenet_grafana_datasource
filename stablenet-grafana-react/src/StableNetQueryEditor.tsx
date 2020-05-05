/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, {PureComponent, ChangeEvent} from 'react';
import {FormLabel, Forms} from '@grafana/ui';
import {QueryEditorProps, SelectableValue} from '@grafana/data';
import {StableNetDataSource} from './StableNetDataSource';
import {Mode, StableNetConfigOptions, Unit} from './Types';
import {Target} from "./QueryInterfaces";
import "./css/query-editor.css";
import {LabelValue} from "./ReturnTypes";

type Props = QueryEditorProps<StableNetDataSource, Target, StableNetConfigOptions>;

interface State {
}

export class StableNetQueryEditor extends PureComponent<Props, State> {
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
            includeMinStats: false,
            averageUnit: Unit.MINUTES
        });
    };

    onStatisticLinkChange = (event: ChangeEvent<HTMLInputElement>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({...query, statisticLink: event.target.value});
        onRunQuery(); // executes the query
    };

    getDevices = (v: string) => {
        const {query, onChange, datasource} = this.props;
        return datasource
            .queryDevices(v, query.refId)
            .then(r => {
                onChange({
                    ...query,
                    moreDevices: r.hasMore,
                });
                return r.data;
            });
    };

    onDeviceChange = (v: SelectableValue<number>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            selectedDevice: {label: v.label!, value: v.value!},
            selectedMeasurement: {label: "", value: -1},
            metricPrefix: '',
            metrics: [],
            chosenMetrics: {},
            mode: Mode.MEASUREMENT,
            includeAvgStats: query.includeAvgStats === undefined ? true : query.includeAvgStats,
            includeMaxStats: query.includeMaxStats === undefined ? false : query.includeMaxStats,
            includeMinStats: query.includeMinStats === undefined ? false : query.includeMinStats,
            averageUnit: query.averageUnit ? query.averageUnit : Unit.MINUTES
        });
        onRunQuery();
    };

    getMeasurements = (v: string) => {
        const {query, onChange, datasource} = this.props;
        return datasource
            .findMeasurementsForDevice(query.selectedDevice ? query.selectedDevice.value : -1, v, query.refId)
            .then(r => {
                onChange({...query, moreMeasurements: r.hasMore});
                return r.data;
            });
    };

    onMeasurementChange = (v: SelectableValue<number>) => {
        const {onChange, query, onRunQuery, datasource} = this.props;
        datasource
            .findMetricsForMeasurement(v.value!, query.refId)
            .then(r => {
                onChange({
                    ...query,
                    metrics: r,
                    chosenMetrics: {},
                    metricPrefix: v.label!,
                    selectedMeasurement: {label: v.label!, value: v.value!},
                });
            })
            .then(() => onRunQuery());
    };

    onMetricPrefixChange = (v: ChangeEvent<HTMLInputElement>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            metricPrefix: v.target.value
        });
        onRunQuery();
    };

    onMetricChange = (v: { text: string; key: string; measurementObid: number }) => {
        const {onChange, query, onRunQuery} = this.props;
        let chosenMetrics = query.chosenMetrics;
        chosenMetrics[v.key] = !chosenMetrics[v.key];
        onChange({
            ...query,
            chosenMetrics: chosenMetrics
        });
        onRunQuery();
    };

    onIncludeChange = (v: string) => {
        const {onChange, query, onRunQuery} = this.props;
        switch (v) {
            case "min":
                onChange({
                    ...query,
                    includeMinStats: !query.includeMinStats
                });
                break;
            case "avg":
                onChange({
                    ...query,
                    includeAvgStats: !query.includeAvgStats
                });
                break;
            case "max":
                onChange({
                    ...query,
                    includeMaxStats: !query.includeMaxStats
                });
                break;
        }
        onRunQuery();
    };

    onUseAvgChange = () => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            useCustomAverage: !query.useCustomAverage
        });
        onRunQuery();
    };

    onCustAvgChange = (v: ChangeEvent<HTMLInputElement>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            averagePeriod: v.target.value
        });
        onRunQuery();
    };

    onAvgUnitChange = (v: SelectableValue<number>) => {
        const {onChange, query, onRunQuery} = this.props;
        onChange({
            ...query,
            averageUnit: v.value!
        });
        onRunQuery();
    };

    render() {
        const query = this.props.query;
        const space = {marginTop: "2px"} as React.CSSProperties;
        const singleMetric = {
            textOverflow: "ellipsis",
            overflow: "hidden",
            whiteSpace: "nowrap"
        } as React.CSSProperties;

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
                            />
                        </div>
                    </div>
                </div>

                {
                    !!query.mode ?
                        <div className="gf-form-inline">{/** Statistic Link mode */}
                            <div className="gf-form">
                                <FormLabel
                                    width={11}
                                    tooltip="Copy a link from the StableNet®-Analyzer. Experimental: At the current version, only links containing exactly one measurement are supported."
                                >Link:
                                </FormLabel>

                                <div className="width-19">
                                    <Forms.Input
                                        type="text"
                                        value={query.statisticLink || ''}
                                        spellCheck={false}
                                        tabIndex={0}
                                        onChange={this.onStatisticLinkChange}
                                    />
                                </div>
                            </div>
                        </div>
                        :
                        <div>{/** Measurement mode */}
                            {/**Device dropdown, more devices*/}
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
                            {/**Measurement dropdown, more measurements*/}
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

                            {
                                !!query.selectedMeasurement && !!query.selectedMeasurement.label ?
                                    <div>
                                        {
                                            !query.metrics.length ?
                                                <div className="gf-form">
                                                    <div style={{paddingLeft: "150px"} as React.CSSProperties}>
                                                        <FormLabel width={30}>No metrics available!</FormLabel>
                                                    </div>
                                                </div>
                                                :
                                                <div className="gf-form">
                                                    <div className="gf-form">
                                                        <FormLabel
                                                            width={11}
                                                            tooltip="The input of this field will be added as a prefix to the metrics' names on the chart."
                                                        >Metric prefix:
                                                        </FormLabel>
                                                        <div className="width-19" style={space}>
                                                            <Forms.Input
                                                                type="text"
                                                                value={query.metricPrefix || ''}
                                                                spellCheck={false}
                                                                tabIndex={0}
                                                                onChange={this.onMetricPrefixChange}
                                                            />
                                                        </div>
                                                    </div>

                                                    <FormLabel
                                                        width={11}
                                                        tooltip="Select the metrics you want to display."
                                                    >Metrics:
                                                    </FormLabel>

                                                    <div className="gf-form-inline">
                                                        {
                                                            query.metrics.map(metric =>
                                                                <div className="gf-form">
                                                                    <Forms.Checkbox
                                                                        value={!!query.chosenMetrics[metric.key]}
                                                                        onChange={() => this.onMetricChange(metric)}
                                                                        size={11}
                                                                    />
                                                                    <div style={singleMetric}>
                                                                        <FormLabel width={17}>{metric.text}</FormLabel>
                                                                    </div>
                                                                </div>
                                                            )
                                                        }
                                                    </div>
                                                </div>
                                        }
                                    </div>
                                    :
                                    null
                            }
                        </div>
                }

                <div>
                    {
                        !!(query.selectedMeasurement && query.selectedMeasurement.label) || !!query.mode ?
                            <div>
                                <div className="gf-form">
                                    <div style={!query.mode ? {marginLeft: "415px"} : {}}>
                                        <FormLabel
                                            width={11}
                                            tooltip="Select the statistics you want to display."
                                        >Include Statistics:</FormLabel>
                                    </div>

                                    <div className="gf-form">
                                        <Forms.Checkbox
                                            value={query.includeMinStats}
                                            onChange={() => this.onIncludeChange("min")}
                                            tabIndex={0}
                                        />
                                        <FormLabel width={5}>Min</FormLabel>

                                        <Forms.Checkbox
                                            value={query.includeAvgStats === undefined ? true : query.includeAvgStats}
                                            onChange={() => this.onIncludeChange("avg")}
                                            tabIndex={0}
                                        />
                                        <FormLabel width={5}>Avg</FormLabel>

                                        <Forms.Checkbox
                                            value={query.includeMaxStats}
                                            onChange={() => this.onIncludeChange("max")}
                                            tabIndex={0}
                                        />
                                        <FormLabel width={5}>Max</FormLabel>
                                    </div>
                                </div>

                                <div className="gf-form-inline">
                                    <div className="gf-form" style={{width: "30px"} as React.CSSProperties}>
                                        <Forms.Checkbox
                                            value={query.useCustomAverage}
                                            onChange={() => this.onUseAvgChange()}
                                            tabIndex={0}
                                        />
                                    </div>
                                    <FormLabel
                                        width={11}
                                        tooltip="Allows defining a custom average period. If disabled, Grafana will automatically compute a suiting average period."
                                    >Average Period:
                                    </FormLabel>
                                    <div className="width-10" style={space}>
                                        <Forms.Input
                                            type="number"
                                            value={query.averagePeriod || ''}
                                            spellCheck={false}
                                            tabIndex={0}
                                            onChange={this.onCustAvgChange}
                                            disabled={!query.useCustomAverage}
                                        />
                                    </div>
                                    <div className="gf-form">
                                        <div tabIndex={0}
                                             style={space}>
                                            <Forms.Select<number>
                                                options={this.getUnits()}
                                                value={query.averageUnit || Unit.MINUTES}
                                                onChange={this.onAvgUnitChange}
                                                width={7}
                                                isSearchable={true}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                            :
                            null
                    }
                </div>
            </div>
        );
    }

}
