;(function(angular) {
	'use strict';

	angular.module('app.controllers').controller('MainLayoutCtrl', ['$scope','$location', 'Auth',
		function($scope, $location, Auth) {
			$scope.signOut = function() {
				Auth.logout().success(function() {
					$location.path('/login');
				});
			};

			$scope.menu = [
				[
					{
						label: 'Inicio',
						icon: 'fa-home',
						link: '#/main/home'
					},
					{
						label: 'Asignaturas',
						icon: 'fa-book',
						link: '#'
					},
					{
						label: 'Alumnos',
						icon: 'fa-child',
						link: '#'
					},
					{
						label: 'Docentes',
						icon: 'fa-male',
						link: '#'
					},
					{
						label: 'Ingreso de Notas',
						icon: 'fa-file-text-o',
						link: '#'
					}
				],
				[
					{
						label: 'Boletines',
						icon: 'fa-bullhorn',
						link: '#'
					},
					{
						label: 'Reporte #1',
						icon: 'fa-calendar',
						link: '#'
					},
					{
						label: 'Reporte #2',
						icon: 'fa-edit',
						link: '#'
					},
					{
						label: 'Reporte #3',
						icon: 'fa-area-chart',
						link: '#'
					},
				],
				[
					{
						label: 'Usuarios',
						icon: 'fa-users',
						link: '#/main/users/all'
					},
					{
						label: 'Roles',
						icon: 'fa-legal',
						link: '#'
					},
					{
						label: 'Privilegios',
						icon: 'fa-trophy',
						link: '#'
					},
					{
						label: 'Configuración',
						icon: 'fa-male',
						link: '#'
					},
					{
						label: 'Cerrar Sesión',
						icon: 'fa-sign-out',
						link: '#',
						onClick: $scope.signOut
					},
				],
			];

			$scope.activeOption = $scope.menu[0][0].label;

			$scope.optionOnClick = function(option) {
				$scope.activeOption = option.label;

				if (typeof option.onClick === 'undefined') {
					console.log('no onclick');
					return;
				}

				console.log('executing....');
				option.onClick();
			};

		}
	]);
})(angular);