package infra

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func NewViper(filename string) *viper.Viper {
	v := viper.New()

	// 1. Set Nama File Config & Path
	v.SetConfigName(filename)
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")

	// 2. Baca File Config
	if err := v.ReadInConfig(); err != nil {
		// cek: Apakah errornya karena "File Tidak Ditemukan"?
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// 3. Enable Environment Variables (Override)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return v
}
