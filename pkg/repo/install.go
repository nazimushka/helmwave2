package repo

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func (rep *Config) Install(settings *helm.EnvSettings, f *repo.File) error {
	if !rep.Force && f.Has(rep.Name) {
		existing := f.Get(rep.Name)
		if rep.Entry != *existing {
			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return fmt.Errorf("repository name (%s) already exists, please specify a different name", rep.Name)
		}

		// The add is idempotent so do nothing
		log.Infof("%q already exists with the same configuration, skipping", rep.Name)
		return nil
	}

	chartRepo, err := repo.NewChartRepository(&rep.Entry, getter.All(settings))
	if err != nil {
		return err
	}

	chartRepo.CachePath = settings.RepositoryCache

	_, err = chartRepo.DownloadIndexFile()
	if err != nil {
		log.Warnf("⚠️ looks like %v is not a valid chart repository or cannot be reached", rep.URL)
	}

	f.Update(&rep.Entry)

	return nil
}
