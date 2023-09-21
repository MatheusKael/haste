package cache

import (
	"encoding/json"
	"os"
)

const TempDir = "./tmp/body.json"

func ReadCacheData() interface{} {

	var jsonTmpData interface{}

	cacheData, err := os.OpenFile(TempDir, os.O_RDWR|os.O_CREATE, 0666)

	defer cacheData.Close()

	decoder := json.NewDecoder(cacheData)

	err = decoder.Decode(&jsonTmpData)

	if err != nil {
		panic(err)
	}

	return jsonTmpData
}


func WriteToCache(text string) { 

		cacheData, err := os.OpenFile(TempDir, os.O_RDWR|os.O_CREATE, 0666)

		cacheData.Truncate(0)
		cacheData.Seek(0, 0)

		defer cacheData.Close()

		_, err = cacheData.Write([]byte(text))

		if err != nil {
			panic(err)
		}
}
