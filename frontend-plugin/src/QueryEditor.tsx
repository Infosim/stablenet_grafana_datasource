/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent, PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { LabelValue, Metric, Mode, StableNetConfigOptions, Target, Unit } from './Types';
import { MetricPrefix } from './components/MetricPrefix';
import { DeviceMenu } from './components/DeviceMenu';
import { StatLink } from './components/StatLink';
import { ModeChooser } from './components/ModeChooser';
import { CustomAverage } from './components/CustomAverage';
import { MinMaxAvg } from './components/MinMaxAvg';
import { Checkbox, InlineFormLabel } from '@grafana/ui';
import { MeasurementMenu } from './components/MeasurementMenu';

// @ts-ignore Some problems with the generic typing here. Could not solve it yet.
type Props = QueryEditorProps<DataSource, Target, StableNetConfigOptions, Target>;

const singleMetric: React.CSSProperties = {
  textOverflow: 'ellipsis',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  borderLeft: '4px',
};

export class QueryEditor extends PureComponent<Props> {
  onModeChange = (v: SelectableValue<number>) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      mode: v.value!,
      includeMaxStats: false,
      includeAvgStats: true,
      includeMinStats: false,
      averageUnit: Unit.MINUTES,
    });
  };

  onStatisticLinkChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { query, onRunQuery, onChange } = this.props;
    onChange({ ...query, statisticLink: event.target.value });
    onRunQuery(); // executes the query
  };

  getDevices = async (v: string): Promise<LabelValue[]> => {
    const { query, datasource, onChange } = this.props;

    const response = await datasource.queryDevices(v);

    onChange({
      ...query,
      moreDevices: response.hasMore || !!query.moreDevices,
    });

    return response.data;
  };

  onDeviceChange = (v: SelectableValue<number>) => {
    const { onChange, query, onRunQuery, datasource } = this.props;
    datasource
      .findMeasurementsForDevice(v.value!, '')
      .then((r) => {
        onChange({
          ...query,
          moreMeasurements: r.hasMore || !!query.moreMeasurements,
          measurements: r.data,
          measurementFilter: '',
          selectedDevice: { label: v.label!, value: v.value! },
          selectedMeasurement: { label: '', value: -1 },
          metricPrefix: '',
          metrics: [],
          chosenMetrics: [],
          mode: Mode.MEASUREMENT,
          includeAvgStats: query.includeAvgStats === undefined ? true : query.includeAvgStats,
          includeMaxStats: query.includeMaxStats === undefined ? false : query.includeMaxStats,
          includeMinStats: query.includeMinStats === undefined ? false : query.includeMinStats,
          averageUnit: query.averageUnit ? query.averageUnit : Unit.MINUTES,
        });
      })
      .then(() => onRunQuery());
  };

  onMeasurementChange = (v: SelectableValue<number>) => {
    const { onChange, query, onRunQuery, datasource } = this.props;
    datasource
      .findMetricsForMeasurement(v.value!)
      .then((r) => {
        onChange({
          ...query,
          metrics: r,
          chosenMetrics: [],
          metricPrefix: v.label!,
          selectedMeasurement: { label: v.label!, value: v.value! },
        });
      })
      .then(() => onRunQuery());
  };

  onMeasurementFilterChange = (v: ChangeEvent<HTMLInputElement>) => {
    const { datasource, onChange, query, onRunQuery } = this.props;
    const x = v.target.value;
    datasource
      .findMeasurementsForDevice(query.selectedDevice.value, v.target.value)
      .then((r) => {
        onChange({
          ...query,
          moreMeasurements: r.hasMore || !!query.moreMeasurements,
          measurements: r.data,
          measurementFilter: x,
          metricPrefix: '',
          metrics: [],
          chosenMetrics: [],
          mode: Mode.MEASUREMENT,
          includeAvgStats: query.includeAvgStats === undefined ? true : query.includeAvgStats,
          includeMaxStats: query.includeMaxStats === undefined ? false : query.includeMaxStats,
          includeMinStats: query.includeMinStats === undefined ? false : query.includeMinStats,
          averageUnit: query.averageUnit ? query.averageUnit : Unit.MINUTES,
        });
      })
      .then(() => onRunQuery());
  };

  onMetricPrefixChange = (v: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, metricPrefix: v.target.value });
    onRunQuery();
  };

  onMetricChange = ({ key }: Metric) => {
    const { query, onChange, onRunQuery } = this.props;

    let chosenMetrics = query.chosenMetrics;
    const index = chosenMetrics.indexOf(key);

    if (index === -1) {
      chosenMetrics.push(key);
    } else {
      chosenMetrics.splice(index, 1);
    }

    onChange({ ...query, chosenMetrics });
    onRunQuery();
  };

  onIncludeChange = (value: 'min' | 'avg' | 'max') => {
    const { query, onChange, onRunQuery } = this.props;

    switch (value) {
      case 'min':
        onChange({ ...query, includeMinStats: !query.includeMinStats });
        break;
      case 'avg':
        onChange({ ...query, includeAvgStats: !query.includeAvgStats });
        break;
      case 'max':
        onChange({ ...query, includeMaxStats: !query.includeMaxStats });
        break;
    }

    onRunQuery();
  };

  onUseAvgChange = () => {
    const { query, onChange, onRunQuery } = this.props;
    onChange({ ...query, useCustomAverage: !query.useCustomAverage });
    onRunQuery();
  };

  onCustAvgChange = (v: ChangeEvent<HTMLInputElement>) => {
    const { query, onChange, onRunQuery } = this.props;
    onChange({ ...query, averagePeriod: v.target.value });
    onRunQuery();
  };

  onAvgUnitChange = (v: SelectableValue<number>) => {
    const { query, onChange, onRunQuery } = this.props;
    onChange({ ...query, averageUnit: v.value! });
    onRunQuery();
  };

  render() {
    const query = this.props.query;

    return (
      <div>
        <ModeChooser selectedMode={query.mode || Mode.MEASUREMENT} onChange={this.onModeChange} />

        {!!query.mode ? (
          <StatLink link={query.statisticLink || ''} onChange={this.onStatisticLinkChange} />
        ) : (
          <div>
            {/** Measurement mode */}
            <div className="gf-form-inline">
              <DeviceMenu
                selectedDevice={query.selectedDevice}
                hasMoreDevices={query.moreDevices}
                get={this.getDevices}
                onChange={this.onDeviceChange}
              />
            </div>
            <div className="gf-form-inline">
              <MeasurementMenu
                measurements={query.measurements || []}
                hasMoreMeasurements={query.moreMeasurements}
                selected={query.selectedMeasurement}
                onChange={this.onMeasurementChange}
                filter={query.measurementFilter || ''}
                onFilterChange={this.onMeasurementFilterChange}
                disabled={query.selectedDevice === undefined}
              />
            </div>

            {!!query.selectedMeasurement && !!query.selectedMeasurement.label ? (
              <div>
                {!query.metrics.length ? (
                  <div className="gf-form">
                    <div style={{ paddingLeft: '150px' } as React.CSSProperties}>
                      <InlineFormLabel width={30}>No metrics available!</InlineFormLabel>
                    </div>
                  </div>
                ) : (
                  <div className="gf-form" style={{ alignItems: 'baseline' }}>
                    <MetricPrefix value={query.metricPrefix || ''} onChange={this.onMetricPrefixChange} />

                    <InlineFormLabel width={11} tooltip="Select the metrics you want to display.">
                      Metrics:
                    </InlineFormLabel>
                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                      {query.metrics.map((metric) => (
                        <div key={metric.key} style={{ padding: '2px' }}>
                          <Checkbox
                            css=""
                            style={singleMetric}
                            value={query.chosenMetrics.includes(metric.key)}
                            onChange={() => this.onMetricChange(metric)}
                            label={metric.text}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : null}
          </div>
        )}

        {!!(query.selectedMeasurement && query.selectedMeasurement.label) || query.mode === Mode.STATISTIC_LINK ? (
          <div style={{ display: 'flex' }}>
            <CustomAverage
              use={query.useCustomAverage}
              period={query.averagePeriod || ''}
              unit={query.averageUnit || Unit.MINUTES}
              onUseAverageChange={this.onUseAvgChange}
              onUseCustomAverageChange={this.onCustAvgChange}
              onAverageUnitChange={this.onAvgUnitChange}
            />
            <MinMaxAvg
              includeMinStats={query.includeMinStats}
              includeAvgStats={query.includeAvgStats}
              includeMaxStats={query.includeMaxStats}
              onChange={this.onIncludeChange}
            />
          </div>
        ) : null}
      </div>
    );
  }
}
