// Package gameapp provides application initialization and runtime
// environment utilities for the MyHobieMMORPGGame server.
/*
 Copyright 2024 Akif-jpg

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package gameapp

import (
	gameconfig "github.com/Akif-jpg/MyHobieMMORPGGame/config"
	"github.com/gofiber/fiber/v3"
)

type RuntimeEnvEnum int

const (
	DEVELOPER RuntimeEnvEnum = iota
	PRODUCT
)

type GameApp struct {
	App        *fiber.App
	GameConfig *gameconfig.GameConfig
}

func New() *GameApp {
	ga := &GameApp{}
	ga.Init(DEVELOPER)
	return ga
}

func (g *GameApp) Init(ree RuntimeEnvEnum) {
	if g.GameConfig == nil {
		switch ree {
		case DEVELOPER:
			g.GameConfig = (&gameconfig.GameConfig{}).GetConfig(gameconfig.DEVELOPER)
		case PRODUCT:
			g.GameConfig = (&gameconfig.GameConfig{}).GetConfig(gameconfig.PRODUCT)
		default:
			g.GameConfig = (&gameconfig.GameConfig{}).GetConfig(gameconfig.DEVELOPER)
		}
	}
	g.App = fiber.New()
	g.App.Get("health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
	g.App.Listen(":" + g.GameConfig.Port)
}
