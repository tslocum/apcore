// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package framework

import (
	"fmt"

	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/framework/config"
	"github.com/go-fed/apcore/util"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ini.v1"
)

const (
	postgresDB = "postgres"
)

func defaultConfig(dbkind string) (c *config.Config, err error) {
	var dbc config.DatabaseConfig
	dbc, err = defaultDatabaseConfig(dbkind)
	if err != nil {
		return
	}
	c = &config.Config{
		ServerConfig:      defaultServerConfig(),
		OAuthConfig:       defaultOAuth2Config(),
		DatabaseConfig:    dbc,
		ActivityPubConfig: defaultActivityPubConfig(),
	}
	return
}

func defaultServerConfig() config.ServerConfig {
	return config.ServerConfig{
		CookieMaxAge:   86400,
		SaltSize:       32,
		BCryptStrength: bcrypt.DefaultCost,
		RSAKeySize:     1024,
	}
}

func defaultOAuth2Config() config.OAuth2Config {
	return config.OAuth2Config{
		AccessTokenExpiry:  3600,
		RefreshTokenExpiry: 7200,
	}
}

func defaultDatabaseConfig(dbkind string) (d config.DatabaseConfig, err error) {
	d = config.DatabaseConfig{
		DatabaseKind: dbkind,
		// This default is implicit in Go but could change, so here we
		// make it explicit instead
		MaxIdleConns: 2,
		// This default is arbitrarily chosen
		DefaultCollectionPageSize: 10,
	}
	if dbkind != postgresDB {
		err = fmt.Errorf("unsupported database kind: %s", dbkind)
		return
	}
	d.PostgresConfig = defaultPostgresConfig()
	return
}

func defaultActivityPubConfig() config.ActivityPubConfig {
	return config.ActivityPubConfig{
		ClockTimezone:                    "UTC",
		OutboundRateLimitQPS:             10,
		OutboundRateLimitBurst:           50,
		HttpSignaturesConfig:             defaultHttpSignaturesConfig(),
		MaxInboxForwardingRecursionDepth: 50,
		MaxDeliveryRecursionDepth:        50,
	}
}

func defaultHttpSignaturesConfig() config.HttpSignaturesConfig {
	return config.HttpSignaturesConfig{
		Algorithms:      []string{"sha256", "sha512"},
		DigestAlgorithm: "SHA-256",
		GetHeaders:      []string{"(request-target)", "Date", "Digest"},
		PostHeaders:     []string{"(request-target)", "Date", "Digest"},
	}
}

func defaultPostgresConfig() config.PostgresConfig {
	return config.PostgresConfig{}
}

func loadConfigFile(filename string, a app.Application, debug bool) (c *config.Config, err error) {
	util.InfoLogger.Infof("Loading config file: %s", filename)
	var cfg *ini.File
	cfg, err = ini.Load(filename)
	if err != nil {
		return
	}
	c = &config.Config{}
	err = cfg.MapTo(c)
	if err != nil {
		return
	}
	appCfg := a.NewConfiguration()
	err = cfg.MapTo(appCfg)
	if err != nil {
		return
	}
	err = a.SetConfiguration(appCfg)
	if err != nil {
		return
	}
	if debug {
		c.ServerConfig.Host = "localhost"
	}
	return
}

func saveConfigFile(filename string, c *config.Config, others ...interface{}) error {
	util.InfoLogger.Infof("Saving config file: %s", filename)
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, c)
	if err != nil {
		return err
	}
	for _, o := range others {
		err = ini.ReflectFrom(cfg, o)
		if err != nil {
			return err
		}
	}
	return cfg.SaveTo(filename)
}