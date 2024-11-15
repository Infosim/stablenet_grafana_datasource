/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent, memo } from 'react';
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

const singleMetric: React.CSSProperties = {
  textOverflow: 'ellipsis',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  borderLeft: '4px',
};

export const QueryEditor = memo(({ datasource, query, onChange, onRunQuery }: QueryEditorProps<DataSource, Target, StableNetConfigOptions, Target>) => {

  const onModeChange = (v: SelectableValue<number>) => onChange({
    ...query,
    mode: v.value!,
    includeMaxStats: false,
    includeAvgStats: true,
    includeMinStats: false,
    averageUnit: Unit.MINUTES,
  });

  const onStatisticLinkChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, statisticLink: event.target.value });
    onRunQuery();
  };

  const getDevices = async (v: string): Promise<LabelValue[]> => {
    const response = await datasource.queryDevices(v);

    onChange({ ...query, moreDevices: response.hasMore || !!query.moreDevices });

    return response.data;
  };

  const onDeviceChange = async ({ value, label }: SelectableValue<number>) => {
    if (value === undefined || label === undefined) {
      alert('No device selected!');
      return;
    }

    const { hasMore, data } = await datasource.findMeasurementsForDevice(value, '');

    onChange({
      ...query,
      moreMeasurements: hasMore || !!query.moreMeasurements,
      measurements: data,
      measurementFilter: '',
      selectedDevice: { label, value },
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

    onRunQuery();
  };

  const onMeasurementChange = async ({ value, label }: SelectableValue<number>) => {
    if (value === undefined || label === undefined) {
      return;
    }

    const metrics = await datasource.findMetricsForMeasurement(value);

    onChange({ ...query, metrics, chosenMetrics: [], metricPrefix: label, selectedMeasurement: { label, value } });

    onRunQuery();
  };

  const onMeasurementFilterChange = async (v: ChangeEvent<HTMLInputElement>) => {
    const x = v.target.value;

    const { hasMore, data } = await datasource.findMeasurementsForDevice(query.selectedDevice.value, v.target.value);

    onChange({
      ...query,
      moreMeasurements: hasMore || !!query.moreMeasurements,
      measurements: data,
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

    onRunQuery();
  };

  const onMetricPrefixChange = (v: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, metricPrefix: v.target.value });
    onRunQuery();
  };

  const onMetricChange = ({ key }: Metric) => {

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

  const onIncludeChange = (value: 'min' | 'avg' | 'max') => {
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

  const onUseAvgChange = () => {
    onChange({ ...query, useCustomAverage: !query.useCustomAverage });
    onRunQuery();
  };

  const onCustAvgChange = (v: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, averagePeriod: v.target.value });
    onRunQuery();
  };

  const onAvgUnitChange = (v: SelectableValue<number>) => {
    onChange({ ...query, averageUnit: v.value! });
    onRunQuery();
  };

  return (
    <div>
      <ModeChooser selectedMode={query.mode || Mode.MEASUREMENT} onChange={onModeChange} />

      {!!query.mode ? (
        <StatLink link={query.statisticLink || ''} onChange={onStatisticLinkChange} />
      ) : (
        <div>
          {/** Measurement mode */}
          <div className="gf-form-inline">
            <DeviceMenu
              selectedDevice={query.selectedDevice}
              hasMoreDevices={query.moreDevices}
              get={getDevices}
              onChange={onDeviceChange}
            />
          </div>
          <div className="gf-form-inline">
            <MeasurementMenu
              measurements={query.measurements || []}
              hasMoreMeasurements={query.moreMeasurements}
              selected={query.selectedMeasurement}
              onChange={onMeasurementChange}
              filter={query.measurementFilter || ''}
              onFilterChange={onMeasurementFilterChange}
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
                  <MetricPrefix value={query.metricPrefix || ''} onChange={onMetricPrefixChange} />

                  <InlineFormLabel width={11} tooltip="Select the metrics you want to display.">Metrics:</InlineFormLabel>

                  <div style={{ display: 'flex', flexDirection: 'column' }}>
                    {query.metrics.map((metric) => (
                      <div key={metric.key} style={{ padding: '2px' }}>
                        <Checkbox style={singleMetric} value={query.chosenMetrics.includes(metric.key)} onChange={() => onMetricChange(metric)} label={metric.text} />
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
            onUseAverageChange={onUseAvgChange}
            onUseCustomAverageChange={onCustAvgChange}
            onAverageUnitChange={onAvgUnitChange}
          />
          <MinMaxAvg
            includeMinStats={query.includeMinStats}
            includeAvgStats={query.includeAvgStats}
            includeMaxStats={query.includeMaxStats}
            onChange={onIncludeChange}
          />
        </div>
      ) : null}
    </div>
  );
});
