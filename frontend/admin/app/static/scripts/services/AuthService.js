;(function(angular) {
	'use strict';

	angular.module('app').service('Auth', [
		function() {
			var authService = {};

			authService.login = function() {
				//Call Backend
			};

			return authService;
		}
	]);
})(angular);