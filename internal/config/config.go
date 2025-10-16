package config

import "github.com/allscorpion/pokedexcli/internal/pokecache"

type Config struct {
	Next     string
	Previous string
	Cache   *pokecache.Cache
}