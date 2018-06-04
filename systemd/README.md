# systemd

```
cp iiif-server.service.example iiif-server.service

sudo mv iiif-server.service /lib/systemd/system/.

sudo systemctl enable iiif-service.service
sudo systemctl start iiif-service

```

## See also

* https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/