package checker

import (
	"errors"
	"fmt"
	"github.com/bartosian/sui_helpers/suimon/pkg/env"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/oschwald/geoip2-golang"
	"path/filepath"

	"github.com/bartosian/sui_helpers/suimon/cmd/checker/config"
)

const (
	httpClientTimeout        = 3 * time.Second
	suimonGeoDBPath          = "%s/.suimon/geodb.mmdb"
	invalidGeoDBPathProvided = `provide valid geodb path by setting SUIMON_GEODB_PATH env variable
or set it in suimon.yaml`
)

func (checker *Checker) ParseData() error {
	suimonConfig, nodeConfig := checker.suimonConfig, checker.nodeConfig

	if len(nodeConfig.P2PConfig.SeedPeers) == 0 {
		return errors.New("no peers found in config file")
	}

	if suimonConfig.HostLookupConfig.EnableLookup {
		geoDBClient, err := initGeoDBClient(suimonConfig)
		if err != nil {
			return err
		}

		defer geoDBClient.Close()

		checker.geoDbClient = geoDBClient
	}

	var (
		wg         sync.WaitGroup
		errChan    = make(chan error)
		errCounter int
		err        error
	)

	monitorsConfig := suimonConfig.MonitorsConfig

	// parse data for the RPC table
	if monitorsConfig.RPCTable.Display {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := checker.parseRPCHosts(); err != nil {
				errChan <- err
			}
		}()
	}

	// parse data for the NODE table
	if monitorsConfig.NodeTable.Display {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := checker.parseNode(); err != nil {
				errChan <- err
			}
		}()
	}

	// parse data for the PEERS table
	if monitorsConfig.PeersTable.Display {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := checker.parsePeers(); err != nil {
				errChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for parseErr := range errChan {
		err = multierror.Append(err, parseErr)

		errCounter++
	}

	if errCounter == 3 {
		return err
	}

	return nil
}

func initGeoDBClient(suimonConfig config.SuimonConfig) (*geoip2.Reader, error) {
	dbPath := suimonConfig.HostLookupConfig.GeoDbPath

	filePath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, errors.New(invalidGeoDBPathProvided)
	}

	db, err := geoip2.Open(filePath)
	if err != nil {
		home := os.Getenv("HOME")
		dbPath = env.GetEnvWithDefault("SUIMON_GEODB_PATH", fmt.Sprintf(suimonGeoDBPath, home))

		filePath, err := filepath.Abs(dbPath)
		if err != nil {
			return nil, errors.New(invalidGeoDBPathProvided)
		}

		return geoip2.Open(filePath)
	}

	return db, nil
}