package process

import (
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/cache"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	"log"
)

func TransformURIWithInstructions(u iiifuri.URI, i IIIFInstructions, config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (iiifuri.URI, iiifimage.Image, error) {

	origin := u.Origin()
	target := u.Target()

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return nil, nil, err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, err
	}

	tr_uri, err := transformation.ToURI(target)

	if err != nil {
		return nil, nil, err
	}

	/*

	go run -mod vendor cmd/iiif-process/main.go -instructions-source file:///usr/local/aaronland/go-iiif/docker/config -config-source file:///usr/local/aaronland/go-iiif/docker/config 'idsecret:///spanking.jpg?id=99277455353'
	2019/10/31 14:31:13 MAKE NEW URI idsecret:///992/774/553/53/99277455353/full/!2048,1536/0/color.jpg
	2019/10/31 14:31:13 MAKE NEW URI idsecret:///992/774/553/53/99277455353/full/full/-1/color.jpg
	2019/10/31 14:31:13 failed to process idsecret:///spanking.jpg?id=99277455353&secret=YGuQcVhbgZApgECHpXSdnmeuzCEA5CzS&secret_o=86G6jZPwBtv9BSe7OtUJvbRuoYxJXAqD (b) : Missing id
	2019/10/31 14:31:13 failed to process idsecret:///spanking.jpg?id=99277455353&secret=YGuQcVhbgZApgECHpXSdnmeuzCEA5CzS&secret_o=86G6jZPwBtv9BSe7OtUJvbRuoYxJXAqD (o) : Missing id
	2019/10/31 14:31:13 MAKE NEW URI idsecret:///992/774/553/53/99277455353/-1,-1,320,320/full/0/dither.jpg
	2019/10/31 14:31:13 failed to process idsecret:///spanking.jpg?id=99277455353&secret=YGuQcVhbgZApgECHpXSdnmeuzCEA5CzS&secret_o=86G6jZPwBtv9BSe7OtUJvbRuoYxJXAqD (d) : Missing id
	{"/spanking.jpg":{"dimensions":{},"palette":[{"name":"#4e3c24","hex":"#4e3c24","reference":"vibrant"},{"name":"#9d8959","hex":"#9d8959","reference":"vibrant"},{"name":"#c7bca6","hex":"#c7bca6","reference":"vibrant"},{"name":"#5a4b36","hex":"#5a4b36","reference":"vibrant"}],"uris":{}}}

	*/

	str_uri := fmt.Sprintf("%s:///%s", u.Driver(), tr_uri)

	log.Println("MAKE NEW URI", str_uri)

	new_uri, err := iiifuri.NewURI(str_uri)

	if err != nil {
		return nil, nil, err
	}

	im, err := driver.NewImageFromConfigWithCache(config, source_cache, origin)

	if err != nil {
		return nil, nil, err
	}

	err = im.Transform(transformation)

	if err != nil {
		return nil, nil, err
	}

	new_target := new_uri.Target()
	err = dest_cache.Set(new_target, im.Body())

	if err != nil {
		return nil, nil, err
	}

	return new_uri, im, nil
}
