;(function(angular) {
	'use strict';

	angular.module('app.services').factory('Student', ['$http', 'Api',
		function($http, Api) {
			var studentService = {};

			studentService.findAll = function() {
				return $http({
					url: Api.getRoute('user/findAll/'),
					method: 'GET'
				});
			};

			return studentService;
		}
	]);
})(angular);
