package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
)

func init() {
	os.Setenv("MYAPP_DEBUG", "false")
	os.Setenv("MYAPP_PORT", "8080")
	os.Setenv("MYAPP_USER", "Gopher")
	os.Setenv("MYAPP_RATE", "0.5")
	os.Setenv("MYAPP_TIMEOUT", "3m")
	os.Setenv("MYAPP_ROLES", "get,put,delete")
	os.Setenv("MYAPP_CODES", "key1:10,key2:100")

	os.Setenv("MYAPP_MANUAL_OVERRIDE_1", "manual override value")
	os.Setenv("MYAPP_MANUAL_OVERRIDE1", "never got")

	os.Setenv("MYAPP_AUTO_SPLIT_VAR", "auto split value")
	os.Setenv("MYAPP_AUTOSPLITVAR", "never got")

	os.Setenv("MYAPP_REQUIREDVAR", "required value")
	os.Setenv("MYAPP_IGNOREDVAR", "ignored value")

	os.Setenv("BACKEND_CREDENTIAL", "wheninthegodoasthegophersdo")

	os.Setenv("MYAPP_URL", "https://golang.org")
}

type Cfg struct {
	Debug   bool
	Port    int
	User    string
	Roles   []string
	Rate    float64
	Timeout time.Duration
	Codes   map[string]int

	ManualOverride1   string `envconfig:"manual_override_1"`
	AutoSplitVar      string `split_words:"true"`
	DefaultVar        string `default:"default_value"`
	RequiredVar       string `required:"true"`
	IgnoredVar        string `ignored:"true"`
	BackendCredential string `envconfig:"BACKEND_CREDENTIAL"`

	URL URLDecorder
}

type URLDecorder url.URL

func (d *URLDecorder) Decode(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return err
	}
	*d = URLDecorder(*u)
	return nil
}

func main() {
	var cfg Cfg
	err := envconfig.Process("myapp", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(cfg)
}
