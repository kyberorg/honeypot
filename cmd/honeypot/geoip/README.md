# GeoIP

GeoIP enriches access log with geoip information (city, region, country) based on IP. 
Needs [Maxmind GeoIP2 City-database](https://www.maxmind.com/en/geoip2-city) to work correctly.

### GeoIP information added
* Country (code, name)
* City name (name)
* Region (code, name)
* Location (latitude, longitude)

### GeoIP2 database
Maxmind GeoIP2 database can be downloaded for free from its [official site](https://www.maxmind.com/en/geoip2-city).

Database should be in MMDB format, once downloaded please refer to downloaded file in [mmdb-file param](#DatabaseFile)

For this you need to register and get `accountID` and `LicenseKey`.

#### Arch/Manjaro way
For Arch/Manjaro Linux you can install `geoipupdate` package
```shell
pacman -S geoipupdate
```
* Add your `accountID` and `LicenseKey` to `/etc/GeoIP.conf`
* Run `geoipupdate`
```shell
geoipupdate
```
and find mmdb file at `/var/lib/GeoIP/GeoLite2-City.mmdb`

## Params

#### DatabaseFile
* Path to MMDB database
```shell
--geoip-mmdb-file=/path/to/GeoLite2-City.mmdb
```
