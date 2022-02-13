package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"github.com/zpatrick/go-config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var defaults = Options{
	"url":          "http://127.0.0.1:8500",
	"pollinterval": "1",
	"httptimeout":  "10",
	"configmask":   "*.yml",
	"loglevel":     "info",
	"namepattern":  "postgresql-%s",
}

type Options map[string]string

type CLIFlag struct {
	Name       string
	DefaultVal string
	Usage      string
}

type CLIProvider struct {
	Flags          []CLIFlag
	ParentSettings Options
	Settings       Options
}

var (
	logError = logger.DefaultLog.Error
)

var settings Options

func initConfig() {
	var err error
	var s Options

	cli := NewCLIProvider([]CLIFlag{
		CLIFlag{
			Name:  "url",
			Usage: "Consul http endpoint\n",
		},
		CLIFlag{
			Name:  "cluster",
			Usage: "Stolon cluster name\n",
		},
		CLIFlag{
			Name:  "service",
			Usage: "Service name\n",
		},
		CLIFlag{
			Name:  "services",
			Usage: "List of service names (separated by commas)\n",
		},
		CLIFlag{
			Name:  "pollinterval",
			Usage: fmt.Sprintf("Consul polling interval (default %s)\n", defaults.get("pollinterval")),
		},
		CLIFlag{
			Name:  "httptimeout",
			Usage: fmt.Sprintf("Consul http client timeout (default %s)\n", defaults.get("httptimeout")),
		},
		CLIFlag{
			Name:  "config",
			Usage: "Config filename\n",
		},
		CLIFlag{
			Name: "configdir",
			Usage: fmt.Sprintf("Directory for configuration files ('%s') for clusters\n",
				defaults.get("configmask")),
		},
		CLIFlag{
			Name:  "configmask",
			Usage: fmt.Sprintf("File mask for config files (default '%s')\n", defaults.get("configmask")),
		},
		CLIFlag{
			Name:  "loglevel",
			Usage: fmt.Sprintf("Logging level (default '%s')\n", defaults.get("loglevel")),
		},
		CLIFlag{
			Name:  "logfile",
			Usage: "Log file\n",
		},
		CLIFlag{
			Name:  "logerrorfile",
			Usage: "Log error file\n",
		},
		CLIFlag{
			Name:  "logformat",
			Usage: "Logging prefix format\n",
		},
		CLIFlag{
			Name:  "deregister",
			Usage: "Deregister all services (any string to apply)\n",
		},
		CLIFlag{
			Name:  "namepattern",
			Usage: fmt.Sprintf("Pattern for consul discovery service name (%s)\n", defaults.get("namepattern")),
		},
	}, s)
	conf := config.NewConfig([]config.Provider{cli})
	configFile, _ := conf.String("config")

	defaultValues := config.NewStatic(defaults.getAll())
	conf.Providers = append(conf.Providers, defaultValues)

	s, err = conf.Settings()
	if err != nil {
		logError.Fatal(err)
	}

	if configFile != "" {
		yamlFile := config.NewYAMLFile(configFile)
		conf.Providers = append(conf.Providers, yamlFile)
	}

	s, err = conf.Settings()
	if err != nil {
		logError.Fatal(err)
	}
	settings = cli.override(s)

	var services []string
	servicesFlag, _ := conf.String("services")
	if servicesFlag != "" {
		services = strings.Split(servicesFlag, ",")
	}
	if service, err := conf.String("service"); err == nil {
		services = append(services, service)
	}
	settings["services"] = strings.Join(services, ",")

	if settings["logfile"] == "" {
		settings["logformat"] = fmt.Sprintf("%%s\t%s\t", settings["cluster"])
	} else {
		settings["logformat"] = "%s\t"
	}
}

func Get() Options {
	return settings
}

func (o Options) get(key string) string {
	return o[key]
}

func (o Options) getAll() Options {
	return o
}

func (o Options) IsSet(key string) bool {
	if _, exist := o[key]; exist && o[key] != "" {
		return true
	}
	return false
}

func NewCLIProvider(flags []CLIFlag, parent Options) *CLIProvider {
	return &CLIProvider{
		Flags:          flags,
		ParentSettings: parent,
		Settings:       Options{},
	}
}

func (o *CLIProvider) Load() (map[string]string, error) {
	ptrs := map[string]*string{}

	for _, f := range o.Flags {
		if _, exist := o.Settings[f.Name]; !exist {
			b := flag.String(f.Name, f.DefaultVal, f.Usage)
			ptrs[f.Name] = b
		}
	}
	flag.Parse()

	for _, f := range o.Flags {
		if _, exist := o.Settings[f.Name]; !exist {
			val := *ptrs[f.Name]
			if _, exist := o.ParentSettings[f.Name]; exist && val == "" {
				o.Settings[f.Name] = o.ParentSettings[f.Name]
			} else {
				o.Settings[f.Name] = val
			}
		}
	}

	return o.Settings, nil
}

func WalkDir(root, pattern string) ([]string, error) {
	var matches []string
	if err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
				return err
			} else if matched {
				matches = append(matches, path)
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return matches, nil
}

func Parse(configFile string) Options {
	var s Options

	if configFile != "" {
		var err error
		globals := config.NewStatic(settings.getAll())
		conf := config.NewConfig([]config.Provider{globals})
		yamlFile := config.NewYAMLFile(configFile)
		conf.Providers = append(conf.Providers, yamlFile)
		s, err = conf.Settings()
		if err != nil {
			logError.Fatal(err)
		}

		var services []string
		servicesFlag, _ := conf.String("services")
		if servicesFlag != "" {
			services = strings.Split(servicesFlag, ",")
		}
		if service, err := conf.String("service"); err == nil {
			services = append(services, service)
		}
		s["services"] = strings.Join(services, ",")

		if s["logfile"] == "" {
			s["logformat"] = fmt.Sprintf("%%s\t%s\t", s["cluster"])
		} else {
			s["logformat"] = "%s\t"
		}
		s["configdir"] = ""
	}

	return s
}

func (c CLIProvider) override(o Options) Options {
	var s = Options{}

	for k, v := range o {
		if _, exist := c.Settings[k]; exist && c.Settings[k] != "" {
			s[k] = c.Settings[k]
		} else {
			s[k] = v
		}
	}
	return s
}

func Validate(options Options) (err error) {
	if options["configdir"] == "" {
		if options["services"] == "" {
			return errors.New("one of the required parameters 'service' or 'services' is not set")
		}
	}

	if _, err = strconv.Atoi(options["pollinterval"]); err != nil {
		return errors.New("can't parse pollinterval")
	}

	if _, err = strconv.Atoi(options["httptimeout"]); err != nil {
		return errors.New("can't parse httptimeout")
	}

	return nil
}

func init() {
	initConfig()
}
