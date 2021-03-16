package clients

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func DotEnv() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dotenvfiles := []string{}

	switch os.Getenv("APP_ENV") {
	case "production":
		dotenvfiles = append(dotenvfiles, ".env.production")
	}

	dotenvfiles = append(dotenvfiles, ".env")

	files := []string{}
	dir := wd
	olddir := dir + "a"
	for dir != olddir {
		for _, dotenvfile := range dotenvfiles {

			filename := filepath.Join(dir, dotenvfile)
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				continue
			}

			files = append(files, filename)
		}

		olddir = dir
		dir = filepath.Dir(dir)
	}

	godotenv.Load(files...)
	return nil
}
