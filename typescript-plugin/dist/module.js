System.register(['./datasource', './query_ctrl'], function(exports_1) {
    var datasource_1, query_ctrl_1;
    var StableNetConfigCtrl, StableNetQueryOptionsCtrl, StableNetAnnotationsQueryCtrl;
    return {
        setters:[
            function (datasource_1_1) {
                datasource_1 = datasource_1_1;
            },
            function (query_ctrl_1_1) {
                query_ctrl_1 = query_ctrl_1_1;
            }],
        execute: function() {
            StableNetConfigCtrl = (function () {
                function StableNetConfigCtrl() {
                    this.passwordExists = this.current.secureJsonFields.snpassword ? true : false;
                }
                StableNetConfigCtrl.prototype.resetPassword = function () {
                    this.passwordExists = false;
                };
                StableNetConfigCtrl.templateUrl = 'partials/config.html';
                return StableNetConfigCtrl;
            })();
            exports_1("StableNetConfigCtrl", StableNetConfigCtrl);
            StableNetQueryOptionsCtrl = (function () {
                function StableNetQueryOptionsCtrl() {
                }
                StableNetQueryOptionsCtrl.templateUrl = 'partials/query.options.html';
                return StableNetQueryOptionsCtrl;
            })();
            StableNetAnnotationsQueryCtrl = (function () {
                function StableNetAnnotationsQueryCtrl() {
                }
                StableNetAnnotationsQueryCtrl.templateUrl = 'partials/annotations.editor.html';
                return StableNetAnnotationsQueryCtrl;
            })();
            exports_1("Datasource", datasource_1.default);
            exports_1("QueryCtrl", query_ctrl_1.StableNetQueryCtrl);
            exports_1("ConfigCtrl", StableNetConfigCtrl);
            exports_1("QueryOptionsCtrl", StableNetQueryOptionsCtrl);
            exports_1("AnnotationsQueryCtrl", StableNetAnnotationsQueryCtrl);
        }
    }
});
//# sourceMappingURL=module.js.map