package mock

//go:generate mockgen -package=mock -destination=job.go    github.com/l2isbad/go.d.plugin/internal/godplugin/job Job
//go:generate mockgen -package=mock -destination=module.go github.com/l2isbad/go.d.plugin/internal/modules       Module
