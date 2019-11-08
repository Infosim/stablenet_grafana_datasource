/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import {QueryCtrl} from 'app/plugins/sdk';
import './css/query-editor.css!'

export class GenericDatasourceQueryCtrl extends QueryCtrl {

  constructor($scope, $injector)  {
    super($scope, $injector);
    this.target.type = this.target.type || 'timeserie';

    this.scope = $scope;
    this.target.server = this.target.server || 'select server';
    this.target.filter = this.target.filter || 'device';
    this.target.deviceORtag = this.target.deviceORtag || 'select option';                 //obid of device or tag-filter to put in rest/devices/measurements/{obid} or rest/measurements/list 
    this.target.measurement = this.target.measurement || 'select measurement';
    this.target.target = this.target.target || 'select metric';
  }


  getFilters() {
    return [{text: 'Device', value: 'device'}, {text: 'Tag Filter', value: 'tag'}];
  }

  getDevices() {
    return this.datasource.deviceFindQuery(this.target.server, this.target.filter);       //calls /rest/devices/list or rest/tag-filters/list and parses answer
  }

  getMeasurements() {
    return this.datasource.measurementFindQuery(this.target.server, this.target.filter, this.target.deviceORtag);      //calls proxy ((calls /rest/devices/measurements/{$device} and parses answer)) and returns answer
  }

  getOptions(query) {
    return this.datasource.deviceFindQuery(this.target.server, this.target.filter);
    //return this.datasource.metricFindQuery(this.target.server, this.target.filter, this.target.deviceORtag, this.target.measurement);       //calls metricFindQuery ((calls StatisticServlet&id={} and parses answer)) and returns answer
  }


  toggleEditorMode() {
    this.target.rawQuery = !this.target.rawQuery;
  }


  onServerChange() {
    this.target.filter = 'device';
    this.target.deviceORtag = 'select option';
    this.target.measurement = 'select measurement';
    this.target.target = 'select metric';
  }

  onFilterChange() {
    this.target.deviceORtag = 'select option';
    this.target.measurement = 'select measurement';
    this.target.target = 'select metric';
  }

  onDeviceChange() {
    this.target.measurement = 'select measurement';
    this.target.target = 'select metric';
  }

  onMeasurementChange() {
    this.target.target = 'select metric';
  }
  
  /**
   * Following bug:
   * 
   * When using the migrated metricFindQuery(), once a metric is chosen, the native 'this.panelCtrl.refresh()' function
   * sets the value of the dropdown menu text (not the menu items!!) to something internal before this internal thing is
   * updated. Such an update happens once metricFindQuery() returns. Therefore the shown text is always one choice 'behind'
   * the current one, although datapoints are correctly represented.
   * 
   * To tackle this, the refresh() function is called with a 0.5s delay, so that metricFindQuery() has time to terminate.
   * This solution is of course temporary, until an alternative is found.
   */
  onChangeInternal() {
    setTimeout(() => {
      console.log('Refresh Later');
      this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }, 500)
  }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

