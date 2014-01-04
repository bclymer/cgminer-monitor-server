(function() {
	var stats = [];
	var enabled = [
		"status",
		"temperature",
		"fanSpeed",
		"fanPercent",
		"mhsLastFiveSeconds",
		"accepted",
		"difficultyAccepted",
		"rejected",
		"difficultyRejected",
		"utility",
		"lastValidWork",
		"hardwareErrors"
	]
	$.get('/stats', function(data) {
		machines = data;
		var devicesContainer = $('#devices-container');
		_.each(machines, function(machine, machineName) {
			_.each(machine, function(device, deviceId) {
				var currentStats = [];
				devicesContainer.append('<div class="page-header"><h1>' + machineName + ' - ' + deviceId + '</h1></div>');
				var containerId = machineName + '-' + deviceId;
				devicesContainer.append('<div class="bs-example" id="' + containerId + '">');
				var chartHolder = $('#' + containerId);
				_.each(device, function(stats, key) {
					if (enabled.indexOf(key) == -1) return true;
					currentStats.push({
						key: key,
						value: stats[stats.length - 1]
					});
					var id = machineName+'.'+deviceId+'.'+key;
					var canvas = $('<div/>')
						.width(800).height(400)
						.prop('width', '800').prop('height', '400')
						.prop('id', id);
					chartHolder.append(canvas);
					var map = [];
					for (var i = 0; i < device['when'].length; i++) {
						map[i] = { x: new Date(device['when'][i] * 1000), y: stats[i] }
					}
					var chart = new CanvasJS.Chart(id, {
						title: {
							text: key.replace(/([A-Z])/g, ' $1').replace(/^./, function(str){ return str.toUpperCase(); })
						},
						axisY:{
							includeZero: false,
						},
						data: [{
							type: "line",
							dataPoints: map
						}]
					});
					chart.render();
				});
				var table = '<h2>Current Stats</h2><table border="1" cellpadding="5"><tr>';
				_.each(currentStats, function(stat) {
					table += '<th>' + stat.key.replace(/([A-Z])/g, ' $1').replace(/^./, function(str){ return str.toUpperCase(); }) + '</th>';
				});
				table += '</tr><tr>';
				_.each(currentStats, function(stat) {
					table += '<td>' + stat.value + '</td>';
				});
				table += '</tr></table>';
				chartHolder.prepend(table);
			});
		});
	});
})();