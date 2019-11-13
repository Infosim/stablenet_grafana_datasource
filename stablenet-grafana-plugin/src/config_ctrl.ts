/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
export class GenericConfigCtrl {

    constructor($scope, $injector) {
        this.current.secureJsonData = this.current.secureJsonData || {};
        this.current.secureJsonData.snip = this.current.secureJsonData.snip || '127.0.0.1';
        this.current.secureJsonData.snport = this.current.secureJsonData.snport || '5443';
        this.current.secureJsonData.snusername = this.current.secureJsonData.username || '';
        this.current.secureJsonData.snpassword = this.current.secureJsonData.snpassword || '';
    }
}

GenericConfigCtrl.templateUrl = 'partials/config.html';