package bouncer

import "strings"

type config struct {
	DefaultRedirect string          `toml:"default_redirect"`
	Packages        []packageConfig `toml:"packages"`
}

type packageConfig struct {
	Prefix   string `toml:"prefix"`
	Import   string `toml:"import"`
	Redirect string `toml:"redirect"`
}

func (c *config) FindPackage(path string) packageConfig {
	for _, pkgConf := range c.Packages {
		prefix := strings.TrimSuffix(pkgConf.Prefix, "/")
		if !strings.HasPrefix(path, prefix) {
			continue
		}

		rest := path[len(prefix):]
		if len(rest) != 0 && !strings.HasPrefix(rest, "/") {
			continue
		}

		return pkgConf
	}

	return packageConfig{}
}
