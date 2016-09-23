window.addEventListener('load', function(e){

	var id = '184512_5f7f47e5b3c66207_x.jpg';
	
	var map = L.map('map', {
		center: [0, 0],
		crs: L.CRS.Simple,
		zoom: 1,
		minZoom: 1,
	});
	
	var info = location + 'tiles/' + id + '/info.json';

	var opts = {
		'quality': 'primitive:5,200,255',
		'tileFormat': 'gif',
	};

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
		       
		       	

