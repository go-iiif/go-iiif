window.addEventListener('load', function(e){

    var qs = window.location.search;
    qs = qs.substring(1);
    
    var params = {};
    var queries = qs.split("&amp;");
    var count = queries.length;
    
    for ( var i = 0; i < count; i++ ) {
	temp = queries[i].split('=');
	params[temp[0]] = temp[1];
    }
    
    var id = '184512_5f7f47e5b3c66207_x.jpg';	// disk
    
    if (params['id']){
	id = params['id'];
    }			 

    var map = L.map('map', {
	center: [0, 0],
	crs: L.CRS.Simple,
	zoom: 1,
	minZoom: 1,
    });
    
    var i = document.getElementById("image");
    i.onclick = function(){
	
	leafletImage(map, function(err, canvas) {
	    
    	    if (err){
    		console.log(err);
    		alert("Argh! There was a problem capturing your image");
    		return false;
    	    }
	    
            var dt = new Date();
            var iso = dt.toISOString();
            var iso = iso.split('T');
            var ymd = iso[0];
            ymd = ymd.replace("-", "", "g");
	    
            var bounds = map.getPixelBounds();
	    var zoom = map.getZoom();

	    var pos = [
		bounds.min.x,
		bounds.min.y,
		bounds.max.x,
		bounds.max.y,
		zoom
	    ];

	    pos = pos.join("-");

            var name = id + "-" + ymd + "-" + pos + ".png";
	    
    	    canvas.toBlob(function(blob) {
    		saveAs(blob, name);
            });
	    
    	    // window.open(body);
	});
    };
    
    var info = 'http://' + location.host + '/' + id + '/info.json';

    var opts = {
	'quality': 'default',
	'tileFormat': 'jpg',
    };

    map.addLayer(L.tileLayer.iiif(info, opts));
    
});
