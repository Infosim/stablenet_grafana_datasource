/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { GenericDatasource } from './datasource';
import { GenericDatasourceQueryCtrl } from './query_ctrl';

export class GenericConfigCtrl {
  static templateUrl = 'partials/config.html';

  passwordExists: boolean;
  current: any;

  constructor() {
    this.passwordExists = this.current.secureJsonFields.snpassword ? true : false;
  }

  resetPassword() {
    this.passwordExists = false;
  }
}

class GenericQueryOptionsCtrl {
  static templateUrl = 'partials/query.options.html';
}

class GenericAnnotationsQueryCtrl {
  static templateUrl = 'partials/annotations.editor.html';
}

export {
  GenericDatasource as Datasource,
  GenericDatasourceQueryCtrl as QueryCtrl,
  GenericConfigCtrl as ConfigCtrl,
  GenericQueryOptionsCtrl as QueryOptionsCtrl,
  GenericAnnotationsQueryCtrl as AnnotationsQueryCtrl,
};
