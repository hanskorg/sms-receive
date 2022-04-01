package conf

import (
	"flag"
	"fmt"
	"github.com/hanskorg/logkit"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

var (
	C        *Conf
	Debug    bool
	confPath string
	bind     string
	logger   io.Writer
)

func init() {
	flag.StringVar(&bind, "bind", ":9909", "http server listen, default: 0.0.0.0:9909.")
	flag.StringVar(&confPath, "conf", "conf/config.yaml", "config file path, default: conf/config.yaml")
}

type Zone struct {
	ID     int    `yaml:"id" json:"id"`
	Short  string `yaml:"short" json:"short"`
	County string `yaml:"county" json:"county"`
	CN     string `yaml:"simple" json:"simple"`
	EN     string `yaml:"en" json:"en"`
	Active bool   `yaml:"active" json:"active"`
}

type Twilio struct {
	Token string `yaml:"authToken" json:"authToken"`
	SID   string `yaml:"accountSID" json:"accountSID"`
	Hook  string `yaml:"hook" json:"hook"`
}

type Database struct {
	DSN           string `yaml:"dsn" json:"dsn"`
	MaxConnection int    `yaml:"maxConnection" json:"maxConnection"`
	MaxIdle       int    `yaml:"maxIdle" json:"maxIdle"`
}

type Conf struct {
	Zone        []*Zone   `yaml:"zone" json:"zone"`
	Debug       bool      `yaml:"debug" json:"debug"`
	Bind        string    `yaml:"bind" json:"bind"`
	TemplateDir string    `yaml:"templateDir" json:"templateDir"`
	Twilio      *Twilio   `yaml:"twilio" json:"twilio"`
	Database    *Database `yaml:"database" json:"database"`
}

func Init() *Conf {
	C = &Conf{
		Zone:     make([]*Zone, 0),
		Bind:     "",
		Twilio:   &Twilio{},
		Database: &Database{},
	}
	bs, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("config file load fail: %s", err.Error()))
	}
	yaml.Unmarshal(bs, C)
	if bind != "" {
		C.Bind = bind
	}
	logkit.SetAlsoStdout(true)
	logkit.SetWithCaller("file")
	logkit.SetName("smsfree")
	logger, _ = logkit.Init()
	return C
}
