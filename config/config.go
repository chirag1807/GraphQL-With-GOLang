package config

import (
	"articlewithgraphql/api/model/dto"
	"articlewithgraphql/constants"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var DatabaseConfig dto.Database
var JWtSecretConfig dto.JWTSecret

func LoadEnv() {

	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("D:/Article Project With GraphQL/.config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	if err := viper.Unmarshal(&DatabaseConfig); err != nil {
		fmt.Println("Error While Decoding .env File")
	}

	fileContent, err := os.ReadFile(constants.SECRET_JSON_FILE_PATH)

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(fileContent, &JWtSecretConfig)

	if err != nil {
		fmt.Println(err)
	}

}
