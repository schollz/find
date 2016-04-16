/*!
	Wheelzoom 3.0.0
	license: MIT
	http://www.jacklmoore.com/wheelzoom
*/
window.wheelzoom = (function(){
	var defaults = {
		zoom: 0.10
	};
	var canvas = document.createElement('canvas');
	var lastPositionX;
	var lastPositionY;
	var dragging = false;

	function setSrcToBackground(img) {
		img.style.backgroundImage = "url('"+img.src+"')";
		img.style.backgroundRepeat = 'no-repeat';
		canvas.width = img.naturalWidth;
		canvas.height = img.naturalHeight;
		img.src = canvas.toDataURL();
	}

	main = function(img, options){
		if (!img || !img.nodeName || img.nodeName !== 'IMG') { return; }

		var settings = {};
		var width;
		var height;
		var bgWidth;
		var bgHeight;
		var bgPosX;
		var bgPosY;
		var previousEvent;
		var currentZoom;

		function updateBgStyle() {
			if (bgPosX > 0) {
				bgPosX = 0;
			} else if (bgPosX < width - bgWidth) {
				bgPosX = width - bgWidth;
			}

			if (bgPosY > 0) {
				bgPosY = 0;
			} else if (bgPosY < height - bgHeight) {
				bgPosY = height - bgHeight;
			}
			
			
			if ($('#dot').length > 0)
				$('#dot').remove();
			var color = '#FF0000';
			var size = '10px';
			$("#map-image-wrapper").append(
			$("<div></div>")
				.attr('id','dot')
				.css('position', 'absolute')
				.css('top', lastPositionY*(bgHeight/height)+bgPosY  + 'px')
				.css('left', lastPositionX*(bgWidth/width)+bgPosX + 'px')
				.css('width', size)
				.css('height', size)
				.css('background-color', color)
			);
			img.style.backgroundSize = bgWidth+'px '+bgHeight+'px';
			img.style.backgroundPosition = bgPosX+'px '+bgPosY+'px';
		}

		function reset() {
			bgWidth = width;
			bgHeight = height;
			bgPosX = bgPosY = 0;
			updateBgStyle();
		}

		function onwheel(e) {
			var deltaY = 0;

			e.preventDefault();

			if (e.deltaY) { // FireFox 17+ (IE9+, Chrome 31+?)
				deltaY = e.deltaY;
			} else if (e.wheelDelta) {
				deltaY = -e.wheelDelta;
			}

			// As far as I know, there is no good cross-browser way to get the cursor position relative to the event target.
			// We have to calculate the target element's position relative to the document, and subtrack that from the
			// cursor's position relative to the document.
			var rect = img.getBoundingClientRect();
			var offsetX = e.pageX - rect.left - document.body.scrollLeft;
			var offsetY = e.pageY - rect.top - document.body.scrollTop;

			// Record the offset between the bg edge and cursor:
			var bgCursorX = offsetX - bgPosX;
			var bgCursorY = offsetY - bgPosY;
			
			// Use the previous offset to get the percent offset between the bg edge and cursor:
			var bgRatioX = bgCursorX/bgWidth;
			var bgRatioY = bgCursorY/bgHeight;

			// Update the bg size:
			if (deltaY < 0) {
				bgWidth += bgWidth*settings.zoom;
				bgHeight += bgHeight*settings.zoom;
				currentZoom += currentZoom*settings.zoom;
			} else {
				bgWidth -= bgWidth*settings.zoom;
				bgHeight -= bgHeight*settings.zoom;
				currentZoom -= currentZoom*settings.zoom;
			}

			// Take the percent offset and apply it to the new size:
			bgPosX = offsetX - (bgWidth * bgRatioX);
			bgPosY = offsetY - (bgHeight * bgRatioY);
			// Prevent zooming out beyond the starting size
			if (bgWidth <= width || bgHeight <= height) {
				reset();
			} else {
				updateBgStyle();
			}
						
			if ($('#dot').length > 0)
				$('#dot').remove();
			var color = '#FF0000';
			var size = '10px';
			$("#map-image-wrapper").append(
			$("<div></div>")
				.attr('id','dot')
				.css('position', 'absolute')
				.css('top', lastPositionY*(bgHeight/height)+bgPosY  + 'px')
				.css('left', lastPositionX*(bgWidth/width)+bgPosX + 'px')
				.css('width', size)
				.css('height', size)
				.css('background-color', color)
			);

		}
	


		function drag(e) {
			e.preventDefault();
			bgPosX += (e.pageX - previousEvent.pageX);
			bgPosY += (e.pageY - previousEvent.pageY);
			
			previousEvent = e;
			dragging=true;
			updateBgStyle();
		}

		function removeDrag() {
			document.removeEventListener('mouseup', removeDrag);
			document.removeEventListener('mousemove', drag);
		}

		// Make the background draggable
		function draggable(e) {
			e.preventDefault();
			previousEvent = e;
			document.addEventListener('mousemove', drag);
			document.addEventListener('mouseup', removeDrag);
		}

	/*
	$('#map-image').click(function(e) {
	if (dragging) {
		dragging = false;
	} else {
		if ($('#dot').length > 0)
			$('#dot').remove();
		var offset = $(this).offset();
		var X = (e.pageX - offset.left);
		var Y = (e.pageY - offset.top);
		
		var img = $('#map-image')[0]
		actualX = X * img.naturalWidth / img.width;
		actualY = Y * img.naturalHeight / img.height;
		
		if (bgPosX > 0) {
			bgPosX = 0;
		} else if (bgPosX < width - bgWidth) {
			bgPosX = width - bgWidth;
		}

		if (bgPosY > 0) {
			bgPosY = 0;
		} else if (bgPosY < height - bgHeight) {
			bgPosY = height - bgHeight;
		}
		
		actualX = X;
		actualY = Y;
		var rect = img.getBoundingClientRect();
		var offsetX = e.pageX - rect.left - document.body.scrollLeft;
		var offsetY = e.pageY - rect.top - document.body.scrollTop;
		lastPositionX=(offsetX-bgPosX)/(bgWidth/width);
		lastPositionY=(offsetY-bgPosY)/(bgHeight/height);
		$('#x-coord').val(lastPositionX);
		$('#y-coord').val(lastPositionY);
        var color = '#FF0000';
        var size = '10px';
        $("#map-image-wrapper").append(
            $("<div></div>")
				.attr('id','dot')
                .css('position', 'absolute')
                .css('top', offsetY + 'px')
                .css('left', offsetX + 'px')
                .css('width', size)
                .css('height', size)
                .css('background-color', color)
        );
	}
	});	*/
		
	function click(e) {
	if (dragging) {
		dragging = false;
	} else{
		if ($('#dot').length > 0)
			$('#dot').remove();
		
		if (bgPosX > 0) {
			bgPosX = 0;
		} else if (bgPosX < width - bgWidth) {
			bgPosX = width - bgWidth;
		}

		if (bgPosY > 0) {
			bgPosY = 0;
		} else if (bgPosY < height - bgHeight) {
			bgPosY = height - bgHeight;
		}

		var rect = img.getBoundingClientRect();
		var offsetX = e.pageX - rect.left - document.body.scrollLeft;
		var offsetY = e.pageY - rect.top - document.body.scrollTop;
		lastPositionX=(offsetX-bgPosX)/(bgWidth/width);
		lastPositionY=(offsetY-bgPosY)/(bgHeight/height);
		$('#x-coord').val(lastPositionX/rect.width);
		$('#y-coord').val(lastPositionY/rect.height);
        var color = '#FF0000';
        var size = '10px';
        $("#map-image-wrapper").append(
            $("<div></div>")
				.attr('id','dot')
                .css('position', 'absolute')
                .css('top', offsetY + 'px')
                .css('left', offsetX + 'px')
                .css('width', size)
                .css('height', size)
                .css('background-color', color)
        );
	}
	}
	
		function loaded() {
			var computedStyle = window.getComputedStyle(img, null);

			width = parseInt(computedStyle.width, 10);
			height = parseInt(computedStyle.height, 10);
			bgWidth = width;
			bgHeight = height;
			bgPosX = 0;
			bgPosY = 0;

			setSrcToBackground(img);

			img.style.backgroundSize =  width+'px '+height+'px';
			img.style.backgroundPosition = '0 0';
			img.addEventListener('wheelzoom.reset', reset);

			img.addEventListener('wheel', onwheel);
			img.addEventListener('mousedown', draggable);
			img.addEventListener('click', click);
		}

		img.addEventListener('wheelzoom.destroy', function (originalProperties) {
			console.log(originalProperties);
			img.removeEventListener('wheelzoom.destroy');
			img.removeEventListener('wheelzoom.reset', reset);
			img.removeEventListener('load', onload);
			img.removeEventListener('mouseup', removeDrag);
			img.removeEventListener('mousemove', drag);
			img.removeEventListener('mousedown', draggable);
			img.removeEventListener('wheel', onwheel);
			img.removeEventListener('click', click);

			img.style.backgroundImage = originalProperties.backgroundImage;
			img.style.backgroundRepeat = originalProperties.backgroundRepeat;
			img.src = originalProperties.src;
		}.bind(null, {
			backgroundImage: img.style.backgroundImage,
			backgroundRepeat: img.style.backgroundRepeat,
			src: img.src
		}));

		options = options || {};

		Object.keys(defaults).forEach(function(key){
			settings[key] = options[key] !== undefined ? options[key] : defaults[key];
		});

		if (img.complete) {
			loaded();
		} else {
			function onload() {
				img.removeEventListener('load', onload);
				loaded();
			}
			img.addEventListener('load', onload);
		}
	};

	// Do nothing in IE8
	if (typeof window.getComputedStyle !== 'function') {
		return function(elements) {
			return elements;
		}
	} else {
		return function(elements, options) {
			if (elements && elements.length) {
				Array.prototype.forEach.call(elements, main, options);
			} else if (elements && elements.nodeName) {
				main(elements, options);
			}
			return elements;
		}
	}
}());
