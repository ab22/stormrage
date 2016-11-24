; (function (angular) {
	'use strict';

	angular.module('app.services').factory('StrUtils', [function () {
		var strUtils = {};

		strUtils.formatBytes = function (bytes) {
			if (bytes == 0) {
				return '0 Bytes';
			}

			var k = 1000; // or 1024 for binary
			var sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
			var i = Math.floor(Math.log(bytes) / Math.log(k));

			return parseFloat((bytes / Math.pow(k, i)).toFixed(3)) + sizes[i];
		};

		return strUtils;
	}
	]);
})(angular);
