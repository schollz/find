
/* WEB SOCKET */

	var baseURL = document.URL.split('/map')[0].split('http://')[1];

	function openWS() {
		var webSocket = new WebSocket("ws://"+baseURL+"/websocket");
		webSocket.onmessage = function(e) {
			var data = JSON.parse(e.data);
			console.log(data);
			map.updateFeatureLayers(data);
		};
		webSocket.onclose = function(e) {
			console.log('connection closed');
		};
		return webSocket;
	}

	function sendMessage(socket,web_token,currentBaselayer) {
		waitForSocketConnect(socket, function(){
			try { 
				var data = { "web_token": web_token, "base_layer": currentBaselayer };
				console.log(data);
				if(data.web_token) {
					socket.send(JSON.stringify(data));
				}
			}
			catch(err) {
				console.log(err);
			}
		});
	}

	function waitForSocketConnect(socket, callback){
		// socket.close()	// after attempting connection close
		setTimeout(
			function() {
				if (socket.readyState ===1) {
					console.log('connection is made');
					if (callback != null) {
						callback();
					}
					return;
				}
				else {
					console.log('waiting for connection...');
					waitForSocketConnect(socket, callback);
				}
			}, 10); // wait 10 milliseconds for connection...
	}

	window.onload = function() {
		if("WebSocket" in window) {
			console.log("[SYSTEM]", "WebSocket is supported by your browser!");
		}
		else {
			console.log("[SYSTEM]", "WebSocket is NOT supported by your browser!");
			console.log("[SYSTEM]", "Please upgrade to a modern browser.");
			alert("[SYSTEM] WebSocket is NOT supported by your browser!");
			alert("[SYSTEM] Please upgrade to a modern browser...");
		}
	}



/* AJAX TO API SERVER */

	function ipsAPI(route,data) {
		/*
		*	AJAX CALL TO API SERVER
		*	Args:
		*		route: 'string'
		*		data: {json}
		*/
		var results;
		$.ajax({
			type: "GET",
			async: false,
			data:data,
			url: route,
			dataType: 'JSON',
			success: function (data) {
				try {
					results = data;
					if (results.success) {  console.log(results);  }
					else {  console.log(results.message);  }
				}
				catch(err){  console.log('Error:', err);  }
			}
		});
		return results;
	}



/* INITIATE MAP OBJECT */

	function initMap(div,config) {

		/******************************
		*   Author: Stefan Safranek   *
		*   Date: 04/01/15            *
		*   Version: 1.12.03          *
		******************************/
		/*
		*	Enhances Leaflet Map Obj
		*	Args:
		*		div: '#div_id' <--- places map object
		*		config: {json} <--- sent from api server during login
		*/

	// CREATE MAP OBJ
		/*  feature request: http://osmbuildings.org/examples/Leaflet.php?lat=36.00574&lon=-78.93683&zoom=18  */
		var map = L.map('map',{maxZoom: 22});
		L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png',{ 
			attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a> | F.I.N.D.'
		}).addTo(map);

	// SAVE CONFIG OBJ TO MAP
		map.config = config;

	// SETUP BASE LAYERS
		map.getBaseLayers = function() {
			var baseLayers = { };
			for (var _i=0; _i<this.config.baselayers.length; _i++) {
				var name = this.config.baselayers[_i].building_location + ' ' + this.config.baselayers[_i].building_name + ' ' + this.config.baselayers[_i].building_floor;
				var uuid = this.config.baselayers[_i].uuid;
				baseLayers[uuid] =  L.tileLayer(this.config.baselayers[_i].base_layer, { maxZoom: 25 });
				baseLayers[uuid].name = name;
				baseLayers[uuid].uuid;
				map.active = uuid;	
			}
			return baseLayers;
		}
		map.baseLayers = map.getBaseLayers();	// Initiate baselayer object
		map.active = '';	// Active basemap

	// SETUP FEATURE LAYERS
		map.featureLayers = { };	// Storage container for featurelayers
		map.clearFeaturesLayers = function() {
			for (var _i in this.featureLayers) {
				if (this.hasLayer(this.featureLayers[_i])) {
					this.removeLayer(this.featureLayers[_i]);	// Remove old featurelayers
				}
				this.featureLayers[_i] = null;	// Clear old featurelayers
			}
		}

	// TRIGGERED BY SOCKET CONNECTION
		map.updateFeatureLayers = function(data) {
			for (var _i in data){
				if (_i in this.featureLayers){	// Remove old featurelayers
					if (this.hasLayer(this.featureLayers[_i])) {
						this.removeLayer(this.featureLayers[_i]);
					}
				}
				this.featureLayers[_i] = {};	// Clear old featurelayers
				try {
					this.featureLayers[_i] = map.createFeatureLayer(_i,data[_i]);	// Create new featurelayers
					this.featureLayers[_i].addTo(this);	// Apply new featurelayers to map
				}
				catch(err) { console.log(err); }
			}
		}

	// CREATE FEATURE LAYERS
		map.createFeatureLayer = function(type,data) {
			var featureLayer;
			if (type == 'current') {
				featureLayer = L.geoJson(data, {
					pointToLayer: function (feature, latlng) {
						var personIcon = L.icon({
							iconSize: [31, 31],
							iconAnchor: [13, 27],
							popupAnchor:  [1, -24],
							iconUrl: 'static/markers/male-2-blk.png'
						});
						return L.marker(latlng, {icon: personIcon});
					},
					onEachFeature: function (feature, layer) {
						layer.bindPopup(feature.properties.datetime + "<br>" + feature.properties.description + "<br>LL: " + feature.properties.loglikelihood);
						layer.bindLabel(feature.properties.users)
					},
					filter: function (feature,latlng) {
						if (feature.properties.base_layer === map.active) {
							try {
								if ($("#filter").val() != '') { 
									if (feature.properties.users.indexOf($("#filter").val()) != -1) { return true; }
									else { return false; }
								}
								else { return true; }
							}
							catch(err) { return true; }
						}
					}
				});
				if ($('#tracking').is(":checked")) {
					this.fitBounds(featureLayer.getBounds());
				}
			}
			else if (type == 'history') {
				/*    https://www.mapbox.com/mapbox.js/example/v1.0.0/leaflet-heat/    */
				/*    https://github.com/Leaflet/Leaflet.heat    */
				var pts = [];
				var counts = []
				var maxValue = 0;
				for ( var r in data ) {
					var latitude = data[r].latLng.lat;
					var longitude = data[r].latLng.lng;
					var count = data[r].c;
					if (count>maxValue) {  maxValue = count;  }
					for (i=0; i< count-1; i++) {
						pts.push([latitude+Math.random()*0.00004-0.00002,longitude+Math.random()*0.00004-0.00002]);
					}
					counts.push(count)
				}
				featureLayer = L.heatLayer(pts);
				featureLayer.setOptions({radius: 15, blur: 20, maxZoom: 25});
			}
			else if (type == 'markers') {
				featureLayer =	L.geoJson(data, {
					pointToLayer: function(feature, latlng) {
						return L.circleMarker(latlng, {
							radius: 6,
							fillOpacity: 0.85,
							color: '#00ff00',
							stroke: false
						});
					},
					onEachFeature: function (feature, layer) {
						layer.bindPopup(feature.properties.description);
					}
				});
			}
			else {  featureLayer = L.geoJson(data);  }
			return featureLayer;
		}

	// BASELAYER LAYER CONTROL
		map.layerControl = function() {
			for (var _i in this.baseLayers) {
				var obj = document.createElement('option');
				obj.value = _i;
				obj.text = this.baseLayers[_i].name;
				obj.onclick = function(){
					map.active = $('#basemap').val();
					console.log($('#basemap').val());
				}
				$('#basemap').append(obj);
			}
		}

	// MARKER CONTAINER
		map.markers = [];

	// LAUNCH MAP OBJ
		map.launch = function() {
			try {
			// FIND TRAJECTORY API CALL
				window.findTrajectory = function() {
					var apiRoute = 'api/v1/trajectory';
					var sendData = {
						"web_token": map.config.web_token,
						"base_layer": map.active,
						"x1": map.markers[0].getLatLng().lng,
						"y1": map.markers[0].getLatLng().lat,
						"x2": map.markers[1].getLatLng().lng,
						"y2": map.markers[1].getLatLng().lat
					};
					var result = ipsAPI(apiRoute,sendData);
					if (result.success) {
					// REMOVE MARKERS
						map.removeLayer(map.markers[0]);
						map.removeLayer(map.markers[1]);
					// DESTROY MARKERS
						map.markers = [];
					// CLEAN OLD TRAJECTORY
						if (map.hasLayer(map.featureLayers.trajectory)) {
							map.removeLayer(map.featureLayers.trajectory);
						}
					// APPLY NEW TRAJECTORY
						if (result.trajectory) {
							var trajectory = L.geoJson(result.trajectory, {
								style: { 
									"color": "#ff7800" 
								}
							});
							trajectory.on('mouseover', function(e) {
								console.log(e.layer.feature.properties.distance);
							});
							trajectory.addTo(map);
							map.featureLayers.trajectory = trajectory;
						}
					}
				}
			// INITIATE MAP OBJ
				this.setView(
					[0,0], 1
				); // ZOOM
			// CREATING MARKERS
				this.on('click', function(e) {
					/*  Resources: https://www.mapbox.com/mapbox.js/example/v1.0.0/select-center-form/  */
					console.log(e.latlng);
					// FIND TRAJECTORY
					if (this.markers.length == 1) { 
						var marker = new L.marker(e.latlng, {draggable:'true'});
						marker.setLatLng(e.latlng);
						marker.id = this.markers.length;
						var form = '<b> Find Trajectory </b> <br>';
						form += '<b> <button type="button" onclick="findTrajectory();">Submit</button> </b>';
						form += '<span id="err"></span>';
						marker.bindPopup('<small>' + form + '</small>');
						marker.addTo(this);
						this.markers.push(marker);
					}
					if (this.markers.length == 0) {
						var marker = new L.marker(e.latlng, {draggable:'true'});
						marker.setLatLng(e.latlng);
						marker.id = this.markers.length;
						var form = '<b> Find Trajectory </b> <br>';
						form += '<b> <button type="button" onclick="findTrajectory();">Submit</button> </b>';
						form += '<span id="err"></span>';
						marker.bindPopup('<small>' + form + '</small>');
						marker.addTo(this);
						this.markers.push(marker);
					}
				});
			// SET UP LAYER CONTROL IN LEGNED OBJ
				this.layerControl();
			// CHANGE BASEMAP
				$(document).ready(function(){ 
					$('select').on('change', function(){ 
						for (var _c in map.baseLayers) {
							if (_c	== $('#basemap').val() ) {
								map.addLayer(map.baseLayers[_c]);
								map.active = _c;
								map.clearFeaturesLayers();
								sendMessage(map.webSocket, map.config.web_token, map.active);
							}
							else {
								// REMOVE OTHER BASE AND FEATURE LAYERS
								if (map.hasLayer(map.baseLayers[_c])) {
	            					map.removeLayer(map.baseLayers[_c]);
	            				}           				
							}
						}
					});
				// Apply Current Map
					map.addLayer(map.baseLayers[$('#basemap').val()]);
					map.active = $('#basemap').val();
					map.webSocket = openWS();
					sendMessage(map.webSocket, map.config.web_token, map.active);
				// TOGGLE TRACKING
					var legend = L.control({position: 'bottomright'});
					legend.onAdd = function (map) {
						var div = L.DomUtil.create('div', 'info legend');
						div.innerHTML += '<input type="checkbox" id="tracking"> Tracking <br>';
						div.innerHTML += '<input type="text" id="filter" style="width:100px;"> filter';
						return div;
					};
					legend.addTo(map);
					$("#tracking").prop('checked', true); // CHECK THE CHECKBOX
				});
			}
			catch(err) { console.log(err); }
		}
	// START MAP OBJ
		map.launch();
	// RETURN ENHANCED MAP OBJ
		return map;

	}

