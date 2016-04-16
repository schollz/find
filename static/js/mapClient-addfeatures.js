
/* WEB SOCKET */

	var baseURL = document.URL.split('/addfeatures')[0].split('http://')[1];

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
			if (type == 'current') { console.log('current not displayed'); }
			else if (type == 'history') { console.log('history not displayed'); }
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

	// LAUNCH MAP OBJ
		map.launch = function() {
			try {
			// INITIATE MAP OBJ
				this.setView(
					[0,0], 1
				); // ZOOM
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
				});
			}
			catch(err) { console.log(err); }
		}
	// START MAP OBJ
		map.launch();
	// RETURN ENHANCED MAP OBJ
		return map;

	}

