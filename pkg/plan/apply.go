package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

func (p *Plan) Apply() (err error) {
	if len(p.body.Releases) == 0 {
		return release.ErrEmpty
	}

	log.Info("🗄 Sync repositories...")
	err = p.syncRepositories()
	if err != nil {
		return err
	}

	log.Info("🛥 Sync releases...")
	err = p.syncReleases()
	if err != nil {
		return err
	}

	return nil
}

func (p *Plan) syncRepositories() error {
	settings := helm.New()
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return err
	}

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Repositories))

	for i := range p.body.Repositories {
		go func(wg *parallel.WaitGroup, i int) {
			defer wg.Done()
			err := p.body.Repositories[i].Install(settings, f)
			if err != nil {
				log.Fatal(err)
			}
		}(wg, i)
	}

	return f.WriteFile(settings.RepositoryConfig, 0644)
}

func (p *Plan) syncReleases() (err error) {

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, i int) {
			defer wg.Done()
			_, err := p.body.Releases[i].Sync()
			if err != nil {
				log.Fatal(err)
			}
		}(wg, i)
	}

	return nil
}
