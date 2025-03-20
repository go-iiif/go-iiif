# example

This is a modified version of https://github.com/mejackreed/Leaflet-IIIF

Modifications include:

* No remote Javascript
* Updating jQuery to 3.x
* Hardcoding the `quality` parameter in the `getTileUrl` method in leaflet-iiif.js because otherwise it throws a temper tantrum - no idea...
* In `index.html` you will need to change the image identifier to something that actually exists the `--root` directory of your iiif server
