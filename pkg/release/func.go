package release

import (
	"crypto/sha256"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/chart/loader"
	"os"
	"strings"
)

func (rel *Config) In(a []Config) bool {
	for _, r := range a {
		if rel == &r {
			return true
		}
	}
	return false
}

func (rel *Config) PlanValues() {

	for i := len(rel.Values) - 1; i >= 0; i-- {
		if _, err := os.Stat(rel.Values[i]); err != nil {
			if os.IsNotExist(err) {
				rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
			}
		}
	}

}

func (rel *Config) RenderValues(dir string) {
	rel.PlanValues()

	h := sha256.New()

	for i, v := range rel.Values {

		h.Write([]byte(rel.Name + "->" + rel.Options.Namespace))
		s := v + "." + hex.EncodeToString(h.Sum(nil))[:16] + ".plan"

		p := dir + s
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel})
		if err != nil {
			log.Warn(err)
		}

		rel.Values[i] = s
	}

}

func (rel *Config) ReposDeps() (repos []string, err error) {
	chart, err := loader.Load(rel.Chart)
	if err != nil {
		return nil, err
	}

	deps := chart.Metadata.Dependencies

	for _, d := range deps {
		if strings.HasPrefix(d.Repository, "@") {
			d.Repository = helper.TrimFirstRune(d.Repository)
		}
		repos = append(repos, d.Repository)
	}

	return repos, nil
}
