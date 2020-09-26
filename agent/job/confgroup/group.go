package confgroup

import (
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/ilyam8/hashstructure"
)

type Group struct {
	Configs []Config
	Source  string
}

type Config map[string]interface{}

func (c Config) HashIncludeMap(_ string, k, _ interface{}) (bool, error) {
	s := k.(string)
	return !(strings.HasPrefix(s, "__") && strings.HasSuffix(s, "__")), nil
}

func (c Config) Name() string              { v, _ := c.get("name").(string); return v }
func (c Config) Module() string            { v, _ := c.get("module").(string); return v }
func (c Config) FullName() string          { return fullName(c.Name(), c.Module()) }
func (c Config) UpdateEvery() int          { v, _ := c.get("update_every").(int); return v }
func (c Config) AutoDetectionRetry() int   { v, _ := c.get("autodetection_retry").(int); return v }
func (c Config) Priority() int             { v, _ := c.get("priority").(int); return v }
func (c Config) Hash() uint64              { return calcHash(c) }
func (c Config) Source() string            { v, _ := c.get("__source__").(string); return v }
func (c Config) Provider() string          { v, _ := c.get("__provider__").(string); return v }
func (c Config) SetModule(source string)   { c.set("module", source) }
func (c Config) SetSource(source string)   { c.set("__source__", source) }
func (c Config) SetProvider(source string) { c.set("__provider__", source) }

func (c Config) set(key string, value interface{}) { c[key] = value }
func (c Config) get(key string) interface{}        { return c[key] }

func (c Config) Apply(def Default) {
	if c.UpdateEvery() <= 0 {
		v := firstPositive(def.UpdateEvery, module.UpdateEvery)
		c.set("update_every", v)
	}
	if c.AutoDetectionRetry() <= 0 {
		v := firstPositive(def.AutoDetectionRetry, module.AutoDetectionRetry)
		c.set("autodetection_retry", v)
	}
	if c.Priority() <= 0 {
		v := firstPositive(def.Priority, module.Priority)
		c.set("priority", v)
	}
	if c.UpdateEvery() < def.MinUpdateEvery && def.MinUpdateEvery > 0 {
		c.set("update_every", def.MinUpdateEvery)
	}
	if c.Name() == "" {
		c.set("name", c.Module())
	} else {
		c.set("name", cleanName(c.Name()))
	}
}

func cleanName(name string) string {
	return reSpace.ReplaceAllString(name, "_")
}

var reSpace = regexp.MustCompile(`\s+`)

func fullName(name, module string) string {
	if name == module {
		return name
	}
	return module + "_" + name
}

func calcHash(obj interface{}) uint64 {
	hash, _ := hashstructure.Hash(obj, nil)
	return hash
}

func firstPositive(value int, others ...int) int {
	if value > 0 || len(others) == 0 {
		return value
	}
	return firstPositive(others[0], others[1:]...)
}
