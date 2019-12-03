/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import {GenericDatasource} from './datasource';
import {GenericDatasourceQueryCtrl} from './query_ctrl';

export class GenericConfigCtrl {
  constructor(){
    this.passwordExists = this.current.secureJsonFields.snpassword ? true : false;
  }

  resetPassword(){
    this.passwordExists = false;
  }
}
GenericConfigCtrl.templateUrl = 'partials/config.html';

class GenericQueryOptionsCtrl {}
GenericQueryOptionsCtrl.templateUrl = 'partials/query.options.html';

class GenericAnnotationsQueryCtrl {}
GenericAnnotationsQueryCtrl.templateUrl = 'partials/annotations.editor.html'

export {
  GenericDatasource as Datasource,
  GenericDatasourceQueryCtrl as QueryCtrl,
  GenericConfigCtrl as ConfigCtrl,
  GenericQueryOptionsCtrl as QueryOptionsCtrl,
  GenericAnnotationsQueryCtrl as AnnotationsQueryCtrl
};
