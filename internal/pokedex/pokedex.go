package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"

	c "github.com/allscorpion/pokedexcli/internal/config"
)



func ConvertBytesToJson[T any](data []byte) (T, error) {
    var obj T
    if err := json.Unmarshal(data, &obj); err != nil {
        return *new(T), err
    }
    return obj, nil
}

func GetLocationAreas(config *c.Config, url string) (c.LocationAreas, error) {
	cacheData, exists := config.Cache.Get(url)
	if exists {
		locationArea, err := ConvertBytesToJson[c.LocationAreas](cacheData)
		
		if err != nil {
			return c.LocationAreas{}, err
		}


		return locationArea, nil
	}


	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return c.LocationAreas{}, err
	}

	client := &http.Client{}
	res, err := client.Do(req);
	if err != nil {
		return c.LocationAreas{}, err;
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return c.LocationAreas{}, fmt.Errorf("response failed with status code %v", res.StatusCode)
	}

	var la c.LocationAreas;
	decoder := json.NewDecoder(res.Body);
	if err := decoder.Decode(&la); err != nil {
		return c.LocationAreas{}, err;
	}

	data, err := json.Marshal(la)
	if err != nil {
		return c.LocationAreas{}, err
	}

	config.Cache.Add(url, data);

	return la, nil;
}



func GetLocationArea(config *c.Config, url string) (c.LocationArea, error) {
	cacheData, exists := config.Cache.Get(url)
	if exists {
		locationArea, err := ConvertBytesToJson[c.LocationArea](cacheData)
		
		if err != nil {
			return c.LocationArea{}, err
		}


		return locationArea, nil
	}


	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return c.LocationArea{}, err
	}

	client := &http.Client{}
	res, err := client.Do(req);
	if err != nil {
		return c.LocationArea{}, err;
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return c.LocationArea{}, fmt.Errorf("response failed with status code %v", res.StatusCode)
	}

	var la c.LocationArea;
	decoder := json.NewDecoder(res.Body);
	if err := decoder.Decode(&la); err != nil {
		return c.LocationArea{}, err;
	}

	data, err := json.Marshal(la)
	if err != nil {
		return c.LocationArea{}, err
	}

	config.Cache.Add(url, data);

	return la, nil;
}



func GetPokemonByName(config *c.Config, name string) (c.Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name;

	cacheData, exists := config.Cache.Get(url)
	if exists {
		pokemon, err := ConvertBytesToJson[c.Pokemon](cacheData)
		
		if err != nil {
			return c.Pokemon{}, err
		}


		return pokemon, nil
	}


	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return c.Pokemon{}, err
	}

	client := &http.Client{}
	res, err := client.Do(req);
	if err != nil {
		return c.Pokemon{}, err;
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return c.Pokemon{}, fmt.Errorf("response failed with status code %v", res.StatusCode)
	}

	var la c.Pokemon;
	decoder := json.NewDecoder(res.Body);
	if err := decoder.Decode(&la); err != nil {
		return c.Pokemon{}, err;
	}

	data, err := json.Marshal(la)
	if err != nil {
		return c.Pokemon{}, err
	}

	config.Cache.Add(url, data);

	return la, nil;
}