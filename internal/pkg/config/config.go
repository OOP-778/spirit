/*
 * Copyright 2020-2021 Luke Whrit, Jack Dorland

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * The configuration system we use in Curiosity is not very advanced.
 * It's powered entirely by koanf.

 * First we load some default values from a confmap (L26-L37).
 * Then, on top of that, we load from the `config.toml` file (L47-L49).
 * And then, finally, on top of that file we load from environment variables.

 * We decided on this order, notably having environment variables on top, because of
 * Curiosity possibly being used in dockerized environments, where environment variables
 * are the preferred way of configuration.
 */

package config

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

// Config is the loaded config object
var Config struct {
	Server struct {
		Host              string         `koanf:"host"`
		Port              int            `koanf:"port"`
		CompresssionLevel compress.Level `koanf:"compression_level"`
		Prefork           bool           `koanf:"prefork"`

		Ratelimits struct {
			Requests int           `koanf:"requests"`
			Duration time.Duration `koanf:"duration"`
		} `koanf:"ratelimits"`
	}

	Documents struct {
		IDLength          int   `koanf:"id_length"`
		MaxDocumentLength int   `koanf:"max_document_length"`
		MaxAge            int64 `koanf:"max_age"`
	} `koanf:"documents"`

	Database struct {
		Dialect       string `koanf:"dialect"`
		ConnectionURI string `koanf:"connection_uri"`
	} `koanf:"database"`

	Security struct {
		UseCors bool `koanf:"use_cors"`
	}
}

// Load configuration from file
func Load() error {
	// Set some default values
	k.Load(confmap.Provider(map[string]interface{}{
		"server.host":                   "0.0.0.0",
		"server.port":                   9000,
		"server.compression_level":      -1,
		"server.prefork":                false,
		"server.ratelimits.requests":    200,
		"server.ratelimits.duration":    300_000,
		"documents.id_length":           8,
		"documents.max_document_length": 400_000,
		"documents.max_age":             2592000,
		"security.use_cors":             true,
	}, "."), nil)

	// Load configuration from TOML on top of default values
	if err := k.Load(file.Provider("./config.toml"), toml.Parser()); err != nil {
		log.Fatalf("Error when loading config from file (toml): %v", err)
	}

	// Load environment variables on top of TOML and default values
	err := k.Load(env.Provider("SPACEBIN_", ".", func(s string) string {
		// Strip the `SPACEBIN_` prefix and replace any `_` with `.` so hierarchy is correctly represented.
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "SPACEBIN_")), "_", ".", -1)
	}), nil)

	if err != nil {
		log.Fatalf("Error when loading config from environment: %v", err)
	}

	err = k.Unmarshal("", &Config)

	if err != nil {
		log.Fatalf("Error when un-marshaling config to struct: %v", err)
	}

	return nil
}
