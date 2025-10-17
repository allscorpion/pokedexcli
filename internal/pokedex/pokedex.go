package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"

	c "github.com/allscorpion/pokedexcli/internal/config"
)



func convertBytesToJson[T any](data []byte) (T, error) {
    var obj T
    if err := json.Unmarshal(data, &obj); err != nil {
        return *new(T), err
    }
    return obj, nil
}

func getApiRequest[T any](config *c.Config, url string) (T, error) {
	var zero T

	cacheData, exists := config.Cache.Get(url)
	if exists {
		parsedCacheData, err := convertBytesToJson[T](cacheData)
		
		if err != nil {
			return zero, err
		}

		return parsedCacheData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return zero, fmt.Errorf("response failed with status code %v", res.StatusCode)
	}

	var obj T
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&obj); err != nil {
		return zero, err
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return zero, err
	}

	config.Cache.Add(url, data);

	return obj, nil
}

func GetLocationAreas(config *c.Config, url string) (c.LocationAreas, error) {
	return getApiRequest[c.LocationAreas](config, url)
}

func GetLocationArea(config *c.Config, url string) (c.LocationArea, error) {
	return getApiRequest[c.LocationArea](config, url)
}

func GetPokemonByName(config *c.Config, name string) (c.Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name;
	return getApiRequest[c.Pokemon](config, url)
}