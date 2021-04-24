package geoip

import (
	"errors"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	"github.com/oschwald/geoip2-golang"
	"net"
	"os"
)

var (
	Enabled     bool
	InitError   error
	ReadyToWork bool
)

const (
	DatabaseFileNotExist    = "MMDB Database file not exists"
	DatabaseFileNotReadable = "MMDB Database file is not readable"
	DatabaseFileCorrupted   = "MMDB Database file is corrupted"
	ModuleDisabledErr       = "GeoIP module is disabled"
)

func init() {
	Enabled = config.GetAppConfig().GeoIP.Enabled
	if Enabled {
		InitError = readDatabaseFile(config.GetAppConfig().GeoIP.DatabaseFile)
	} else {
		InitError = errors.New(ModuleDisabledErr)
	}
	ReadyToWork = InitError == nil
}

func LookupIP(ipString string) (*GeoInfo, error) {
	if !ReadyToWork {
		return nil, InitError
	}
	db, readDbErr := geoip2.Open(config.GetAppConfig().GeoIP.DatabaseFile)
	if readDbErr != nil {
		return nil, errors.New(DatabaseFileCorrupted)
	}
	defer db.Close()
	ip := net.ParseIP(ipString)
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}

	return &GeoInfo{
		Coordinates: Coordinates{
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
		},
		Country: Country{
			Code: record.Country.IsoCode,
			Name: record.Country.Names["en-US"],
		},
		Region: Region{
			Code: "",
		},
		City: City{
			Name: record.City.Names["en-US"],
		},
	}, nil
}

func readDatabaseFile(databaseFile string) error {
	databaseFileExists := util.IsFileExists(databaseFile)
	if databaseFileExists {
		_, openError := os.Open(databaseFile)
		if openError != nil {
			return errors.New(DatabaseFileNotReadable)
		} else {
			return nil
		}
	} else {
		return errors.New(DatabaseFileNotExist)
	}
}
