/*
 * Copyright 2020 Luke Whrit, Jack Dorland; The Spacebin Authors

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 *  Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/config"
)

// Start initializes the server
func Start() *fiber.App {
	app := fiber.New(&fiber.Settings{
		Prefork: config.Config.Server.Prefork,
	})

	registerMiddlewares(app)
	registerRoutes(app)

	return app
}