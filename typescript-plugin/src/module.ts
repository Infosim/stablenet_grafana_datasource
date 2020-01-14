/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import RocksetDatasource from './datasource';
import { StableNetQueryCtrl } from './query_ctrl';

export class StableNetConfigCtrl {
  static templateUrl = 'partials/config.html';
  private passwordExists: boolean;
  private current: any;

  constructor(){
    this.passwordExists = this.current.secureJsonFields.snpassword ? true : false;
  }

  resetPassword(): void{
    this.passwordExists = false;
  }
}

class StableNetQueryOptionsCtrl {
  static templateUrl = 'partials/query.options.html';
}

class StableNetAnnotationsQueryCtrl {
  static templateUrl = 'partials/annotations.editor.html';
}

export {
  RocksetDatasource as Datasource,
  StableNetQueryCtrl as QueryCtrl,
  StableNetConfigCtrl as ConfigCtrl,
  StableNetQueryOptionsCtrl as QueryOptionsCtrl,
  StableNetAnnotationsQueryCtrl as AnnotationsQueryCtrl
};
