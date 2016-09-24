window.addEventListener('load', function(e){

	var qs = window.location.search;
	qs = qs.substring(1);

	var qs = window.location.search;
	qs = qs.substring(1);
	
	var params = {};
	var queries = qs.split("&amp;");
	var count = queries.length;
	
	for ( var i = 0; i < count; i++ ) {
		temp = queries[i].split('=');
		params[temp[0]] = temp[1];
	}   

	var mode = 5;
	
	if (params['mode'] == 'triangles'){
		mode = 1;
	}

	if (params['mode'] == 'circles'){
		mode = 4;
	}
	
	var id = '184512_5f7f47e5b3c66207_x.jpg';
	
	var map = L.map('map', {
		center: [0, 0],
		crs: L.CRS.Simple,
		zoom: 1,
		minZoom: 1,
	});

	var root = location.href.replace(location.search, '');
	var info = root + 'tiles/' + id + '/info.json';

	var opts = {
		'quality': 'primitive:' + mode + ',200,255',
		'tileFormat': 'gif',
	};

	console.log(info);
	console.log(opts);
	
	var layer = L.tileLayer.iiif(info, opts);

	layer.on('loading', function(){
		var el = document.getElementById("status");
		el.innerText = "loading tiles";
	});

	layer.on('load', function(){
		var el = document.getElementById("status");
		el.innerText = "";
	});

	map.addLayer(layer);    

});
		       
		       	

