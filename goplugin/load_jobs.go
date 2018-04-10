package goplugin

import (
	"io/ioutil"
	"path"

	"github.com/l2isbad/go.d.plugin/goplugin/job"
	"github.com/l2isbad/go.d.plugin/modules"
	"github.com/l2isbad/toml"
	"github.com/l2isbad/toml/ast"
)

const (
	keyGLOBAL   = "global"
	keyBASE     = "base"
	keySPECIFIC = "specific"
)

type jobStack []*job.Job

func (js *jobStack) pop() *job.Job {
	rv := (*js)[0]
	(*js)[0] = nil
	*js = (*js)[1:]
	return rv
}

func (js *jobStack) empty() bool {
	return len(*js) == 0
}

func (js *jobStack) push(j *job.Job) {
	*js = append(*js, j)
}

// loadJobs() reads and parses configuration files of enabled modules.
func (p *goPlugin) loadJobs() jobStack {
	loaded := jobStack{}
	for modName, creator := range p.modRun {
		conf, mod := job.NewConf(), creator()

		absPath := path.Join(p.modConf, modName+".conf")
		f, err := ioutil.ReadFile(absPath)
		_, noconf := mod.(modules.NoConfiger)

		if !noconf {
			p.Log.Debugf("'%s' can not be loaded without a configuration file", modName)
		}
		if !noconf && err != nil {
			p.Log.Errorf("'%s' skipped: %s", modName, err)
			continue
		}

		conf.SetModuleName(modName)
		setModDefaultConf(modName, conf)

		// PUSH (without configuration files)
		if err != nil {
			p.Log.Debugf("'%s': %s", modName, err)
			loaded.push(job.New(mod, conf))
			continue
		}
		p.Log.Debugf("'%s' configuration read success", modName)

		table, err := toml.Parse(f)
		if err != nil {
			p.Log.Errorf("'%s' skipped: %s", modName, err)
			continue
		}
		p.Log.Debugf("'%s' configuration TOML parse success", modName)

		for _, c := range parseConfAST(table) {
			if !noconf && !c.hasSpecific() {
				p.Log.Errorf("'%s' skipped: no \"specific\" section", modName)
				continue
			}
			confLocal, mod := *conf, creator()

			if c.hasGlobal() {
				if err := toml.UnmarshalTable(c.global, &confLocal); err != nil {
					p.Log.Errorf("'%s' skipped: %s", modName, err)
					break
				}
			}

			if c.hasBase() {
				if err := toml.UnmarshalTable(c.base, &confLocal); err != nil {
					p.Log.Errorf("'%s' %s skipped: %s", modName, c.name, err)
					continue
				}
			}

			if c.hasSpecific() {
				if err := toml.UnmarshalTable(c.specific, mod); err != nil {
					p.Log.Errorf("'%s' %s skipped: %s", modName, c.name, err)
					continue
				}
			}

			confLocal.SetJobName(c.name)
			loaded.push(job.New(mod, &confLocal))
		}
		p.Log.Debugf("'%s' number of loaded jobs: %d", modName, len(loaded))
	}
	return loaded
}

func setModDefaultConf(n string, c *job.Conf) {
	if v, ok := modules.GetDefault(n).Get(modules.UpdateEvery); ok {
		c.UpdEvery = v
	}
	if v, ok := modules.GetDefault(n).Get(modules.ChartCleanup); ok {
		c.ChartCleanup = v
	}
}

type jobRawConf struct {
	name     string
	global   *ast.Table
	base     *ast.Table
	specific *ast.Table
}

func (j *jobRawConf) hasGlobal() bool {
	return j.global != nil
}

func (j *jobRawConf) hasBase() bool {
	return j.base != nil
}

func (j *jobRawConf) hasSpecific() bool {
	return j.specific != nil
}

func parseConfAST(t *ast.Table) []jobRawConf {
	var rv []jobRawConf
	j := jobRawConf{}

	if _, ok := t.Fields[keyGLOBAL]; ok {
		if v, ok := t.Fields[keyGLOBAL].(*ast.Table); ok {
			j.global = v
			delete(t.Fields, keyGLOBAL)
		}
	}

	if len(t.Fields) == 0 && j.global != nil {
		rv = append(rv, j)
		return rv
	}

	for name, options := range t.Fields {
		v, ok := options.(*ast.Table)
		if !ok {
			continue
		}
		lj := jobRawConf{global: j.global}
		lj.name = name
		if v.Fields[keyBASE] != nil {
			if v, ok := v.Fields[keyBASE].(*ast.Table); ok {
				lj.base = v
			}
		}
		if v.Fields[keySPECIFIC] != nil {
			if v, ok := v.Fields[keySPECIFIC].(*ast.Table); ok {
				lj.specific = v
			}
		}
		if lj.base != nil || lj.specific != nil {
			rv = append(rv, lj)
		}
	}
	return rv
}
