package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/allscorpion/pokedexcli/internal/config"
)

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func BytesToLocationArea(data []byte) (LocationArea, error) {
    var la LocationArea
    if err := json.Unmarshal(data, &la); err != nil {
        return LocationArea{}, err
    }
    return la, nil
}

func GetLocationAreas(config *config.Config, url string) (LocationArea, error) {
	cacheData, exists := config.Cache.Get(url)
	fmt.Printf("Cache exists for %s: %v\n", url, exists)
	if exists {
		locationArea, err := BytesToLocationArea(cacheData);
		
		if err != nil {
			return LocationArea{}, err
		}

		fmt.Println("CACHE HIT!!")

		return locationArea, nil
	}


	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		return LocationArea{}, err
	}

	client := &http.Client{}
	res, err := client.Do(req);
	if err != nil {
		return LocationArea{}, err;
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationArea{}, fmt.Errorf("response failed with status code %v", res.StatusCode)
	}

	var la LocationArea;
	decoder := json.NewDecoder(res.Body);
	if err := decoder.Decode(&la); err != nil {
		return LocationArea{}, err;
	}

	data, err := json.Marshal(la)
	if err != nil {
		return LocationArea{}, err
	}

	fmt.Println("CACHE MISS!!")
	config.Cache.Add(url, data);

	return la, nil;
}