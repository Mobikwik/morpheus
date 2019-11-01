package envConfig

import (
	"github.com/magiconair/properties"
	"log"
)

func LoadProperties(propertyFilePath string) *properties.Properties {

	log.Print("loading properties from path ", propertyFilePath)
	// init from a file
	p := properties.MustLoadFile(propertyFilePath, properties.UTF8)

	// or multiple files
	/*p = properties.MustLoadFiles([]string{
		"${HOME}/config.envConfig",
		"${HOME}/config-${USER}.envConfig",
	}, properties.UTF8, true)

	// or from a map
	p = properties.LoadMap(map[string]string{"key": "value", "abc": "def"})

	// or from a string
	p = properties.MustLoadString("key=value\nabc=def")

	// or from a URL
	p = properties.MustLoadURL("http://host/path")

	// or from multiple URLs
	p = properties.MustLoadURLs([]string{
		"http://host/config",
		"http://host/config-${USER}",
	}, true)

	// or from flags
	p.MustFlag(flag.CommandLine)

	// get values through getters
	_ = p.MustGetString("host")
	_ = p.GetInt("port", 8080)

	// or through Decode
	type Config struct {
		Host    string        `envConfig:"host"`
		Port    int           `envConfig:"port,default=9000"`
		Accept  []string      `envConfig:"accept,default=image/png;image;gif"`
		Timeout time.Duration `envConfig:"timeout,default=5s"`
	}
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	*/
	return p
}
