/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import React, { ChangeEvent, PureComponent } from 'react';
//import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { Mode, StableNetConfigOptions, Unit } from './Types';
import { Target } from './QueryInterfaces';
import { LabelValue } from './ReturnTypes';
import { Metric } from './components/Metric';
import { MetricPrefix } from './components/MetricPrefix';
import { DropdownMenu } from './components/DropdownMenu';
import { StatLink } from './components/StatLink';
import { ModeChooser } from './components/ModeChooser';
import { CustomAverage } from './components/CustomAverage';
import { Stats } from './components/Stats';
import { InlineFormLabel } from '@grafana/ui';

type Props = QueryEditorProps<DataSource, Target, StableNetConfigOptions>;

export class QueryEditor extends PureComponent<Props> {
  getModes(): LabelValue[] {
    return [
      { label: 'Measurement', value: Mode.MEASUREMENT },
      { label: 'Statistic Link', value: Mode.STATISTIC_LINK },
    ];
  }

  getUnits(): LabelValue[] {
    return [
      { label: 'sec', value: Unit.SECONDS },
      { label: 'min', value: Unit.MINUTES },
      { label: 'hrs', value: Unit.HOURS },
      { label: 'days', value: Unit.DAYS },
    ];
  }

  onModeChange = (v: SelectableValue<number>) => {
    const { onChange, query } = this.props;
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
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, statisticLink: event.target.value });
    onRunQuery(); // executes the query
  };

  getDevices = (v: string) => {
    const { query, onChange, datasource } = this.props;
    return datasource.queryDevices(v, query.refId).then(r => {
      onChange({
        ...query,
        moreDevices: r.hasMore,
      });
      return r.data;
    });
  };

  onDeviceChange = (v: SelectableValue<number>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({
      ...query,
      selectedDevice: { label: v.label!, value: v.value! },
      selectedMeasurement: { label: '', value: -1 },
      metricPrefix: '',
      metrics: [],
      chosenMetrics: {},
      mode: Mode.MEASUREMENT,
      includeAvgStats: query.includeAvgStats === undefined ? true : query.includeAvgStats,
      includeMaxStats: query.includeMaxStats === undefined ? false : query.includeMaxStats,
      includeMinStats: query.includeMinStats === undefined ? false : query.includeMinStats,
      averageUnit: query.averageUnit ? query.averageUnit : Unit.MINUTES,
    });
    onRunQuery();
  };

  getMeasurements = (v: string) => {
    const { query, onChange, datasource } = this.props;
    return datasource
      .findMeasurementsForDevice(query.selectedDevice ? query.selectedDevice.value : -1, v, query.refId)
      .then(r => {
        onChange({ ...query, moreMeasurements: r.hasMore });
        return r.data;
      });
  };

  onMeasurementChange = (v: SelectableValue<number>) => {
    const { onChange, query, onRunQuery, datasource } = this.props;
    datasource
      .findMetricsForMeasurement(v.value!, query.refId)
      .then(r => {
        onChange({
          ...query,
          metrics: r,
          chosenMetrics: {},
          metricPrefix: v.label!,
          selectedMeasurement: { label: v.label!, value: v.value! },
        });
      })
      .then(() => onRunQuery());
  };

  onMetricPrefixChange = (v: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({
      ...query,
      metricPrefix: v.target.value,
    });
    onRunQuery();
  };

  onMetricChange = (v: { text: string; key: string; measurementObid: number }) => {
    const { onChange, query, onRunQuery } = this.props;
    let chosenMetrics = query.chosenMetrics;
    chosenMetrics[v.key] = !chosenMetrics[v.key];
    onChange({
      ...query,
      chosenMetrics: chosenMetrics,
    });
    onRunQuery();
  };

  onIncludeChange = (v: string) => {
    const { onChange, query, onRunQuery } = this.props;
    switch (v) {
      case 'min':
        onChange({
          ...query,
          includeMinStats: !query.includeMinStats,
        });
        break;
      case 'avg':
        onChange({
          ...query,
          includeAvgStats: !query.includeAvgStats,
        });
        break;
      case 'max':
        onChange({
          ...query,
          includeMaxStats: !query.includeMaxStats,
        });
        break;
    }
    onRunQuery();
  };

  onUseAvgChange = () => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({
      ...query,
      useCustomAverage: !query.useCustomAverage,
    });
    onRunQuery();
  };

  onCustAvgChange = (v: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({
      ...query,
      averagePeriod: v.target.value,
    });
    onRunQuery();
  };

  onAvgUnitChange = (v: SelectableValue<number>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({
      ...query,
      averageUnit: v.value!,
    });
    onRunQuery();
  };

  render() {
    const query = this.props.query;
    const space = { marginTop: '2px' } as React.CSSProperties;
    const singleMetric = {
      textOverflow: 'ellipsis',
      overflow: 'hidden',
      whiteSpace: 'nowrap',
    } as React.CSSProperties;

    return (
      <div>
        <ModeChooser
          space={space}
          mode={query.mode || Mode.MEASUREMENT}
          options={this.getModes}
          onChange={this.onModeChange}
        />

        {!!query.mode ? (
          <StatLink link={query.statisticLink || ''} onChange={this.onStatisticLinkChange} />
        ) : (
          <div>
            {/** Measurement mode */}
            {/**Device dropdown, more devices*/}
            <DropdownMenu
              space={space}
              name={'Device'}
              get={this.getDevices}
              selected={query.selectedDevice}
              onChange={this.onDeviceChange}
              more={query.moreDevices}
            />
            {/**Measurement dropdown, more measurements*/}
            <DropdownMenu
              space={space}
              name={'Measurement'}
              get={this.getMeasurements}
              selected={query.selectedMeasurement}
              onChange={this.onMeasurementChange}
              more={query.moreMeasurements}
            />
            {!!query.selectedMeasurement && !!query.selectedMeasurement.label ? (
              <div>
                {!query.metrics.length ? (
                  <div className="gf-form">
                    <div style={{ paddingLeft: '150px' } as React.CSSProperties}>
                      <InlineFormLabel width={30}>No metrics available!</InlineFormLabel>
                    </div>
                  </div>
                ) : (
                  <div className="gf-form">
                    <MetricPrefix space={space} value={query.metricPrefix || ''} onChange={this.onMetricPrefixChange} />

                    <InlineFormLabel width={11} tooltip="Select the metrics you want to display.">
                      Metrics:
                    </InlineFormLabel>

                    <div className="gf-form-inline">
                      {query.metrics.map(metric => (
                        <Metric
                          singleMetric={singleMetric}
                          value={!!query.chosenMetrics[metric.key]}
                          onChange={() => this.onMetricChange(metric)}
                          text={metric.text}
                        />
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : null}
          </div>
        )}

        <div>
          {!!(query.selectedMeasurement && query.selectedMeasurement.label) || !!query.mode ? (
            <div>
              <Stats
                mode={query.mode}
                values={[query.includeMinStats, query.includeAvgStats, query.includeMaxStats]}
                onChange={this.onIncludeChange}
              />

              <CustomAverage
                space={space}
                use={query.useCustomAverage}
                period={query.averagePeriod || ''}
                unit={query.averageUnit || Unit.MINUTES}
                getUnits={this.getUnits}
                onChange={[this.onUseAvgChange, this.onCustAvgChange, this.onAvgUnitChange]}
              />
            </div>
          ) : null}
        </div>
      </div>
    );
  }
}
