package configs

import (
	"github.com/kamontat/fthelper/shared/fs"
	"github.com/kamontat/fthelper/shared/maps"
	"github.com/kamontat/fthelper/shared/xtemplates"
)

func LoadConfigFromFileSystem(files []fs.FileSystem, data maps.Mapper, strategy maps.Mapper) (maps.Mapper, error) {
	var result = maps.New()
	for _, file := range files {
		if file.IsDir() {
			var files, err = file.ReadDir()
			if err != nil {
				return result, err
			}

			output, err := LoadConfigFromFileSystem(files, data, strategy)
			if err != nil {
				return result, err
			}

			result = maps.Merger(result).Add(output).SetConfig(strategy).Merge()
		} else {
			// read content
			var content, err = file.Read()
			if err != nil {
				return result, err
			}

			// compile template data only if data is not empty
			// If data is empty, then no point to parse templates
			if !data.IsEmpty() {
				str, err := xtemplates.Text(string(content), data)
				if err != nil {
					return result, err
				}
				content = []byte(str)
			}

			// convert content to mapper
			output, err := maps.FromJson(content)
			if err != nil {
				return result, err
			}

			// merge result together
			result = maps.Merger(result).Add(output).SetConfig(strategy).Merge()
		}
	}

	return result, nil
}
