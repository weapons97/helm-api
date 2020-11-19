/*
Copyright © 2020 weapons97@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package adapter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	protos "github.com/weapons97/helm-api/pkg/protos"
	"github.com/weapons97/helm-api/pkg/utils"

	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/helmpath"
	"net/http"
	"path/filepath"
	"unicode/utf8"

	"os"

	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var (
	NULL_POINT_ERROR = fmt.Errorf(`null pointer`)
)

func init() {
	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: http.ProxyFromEnvironment,
		},
	}
}

type Adapter interface {
	InstallCharts(releaseName, repoName, version string, userValues []byte, dryRun bool) (*Release, error)
	UpgradeRelease(releaseName, repoName, version string, userValues []byte,
		dryRun bool, maxhistory int, resetValue bool) (*Release, error)
	UninstallRelease(relName string) error
	GetChartDetail(chartPath, version string) (*Chart, error)

	GetRelease(release string) (*Release, error)
	GetValues(relName string) (string, map[string]interface{}, error)
	ListRelease() ([]*Release, error)

	RollbackRelease(releaseName string, ver int) error
	GetHistory(name string) ([]*Release, error)
	Search(keyword string, regexp bool) (charts []*ChartInfo, e error)
	All() (charts []*ChartInfo, e error)
}

type helmAdapter struct {
	cfg       *action.Configuration
	setting   *cli.EnvSettings
	namespace string
}

func NewAdapter(cfg *action.Configuration,
	setting *cli.EnvSettings,
	namespace string,
) Adapter {
	return &helmAdapter{
		cfg:       cfg,
		setting:   setting,
		namespace: namespace,
	}
}

type Entry repo.Entry

// ConvertEntry convert protos.Entry to adpter Entry
func ConvertEntry(en *protos.Entry) Entry {
	return Entry{
		Name:                  en.GetName(),
		URL:                   en.GetUrl(),
		Username:              en.GetUsername(),
		Password:              en.GetPassword(),
		CertFile:              en.GetCertFile(),
		KeyFile:               en.GetKeyFile(),
		CAFile:                en.GetCaFile(),
		InsecureSkipTLSverify: en.GetInsecureSkipTlsVerify(),
	}
}

type Metadata chart.Metadata

type Chart chart.Chart

// GetChartDetail
func (ha *helmAdapter) GetChartDetail(chartPath, version string) (*Chart, error) {
	client := action.NewInstall(nil)
	client.Version = version
	cp, err := client.ChartPathOptions.LocateChart(chartPath, ha.setting)
	if err != nil {
		return nil, err
	}
	utils.Debug().Str(`chart_path`, cp).Send()
	ct, e := LoadWithRawValues(cp)
	if e != nil {
		return nil, e
	}
	c := Chart(*ct)
	return &c, nil
}

// InstallCharts 安装charts
func (ha *helmAdapter) InstallCharts(releaseName, repoName, version string, userValues []byte,
	dryRun bool) (*Release, error) {
	chartPath := repoName
	client := action.NewInstall(ha.cfg)
	client.ReleaseName = releaseName
	client.Version = version
	client.DryRun = dryRun
	cp, err := client.ChartPathOptions.LocateChart(chartPath, ha.setting)
	if err != nil {
		return nil, err
	}
	fl := utils.NewFileLock(cp)
	err = fl.LockWithTimeout(time.Second * 30)
	if err != nil {
		return nil, err
	}
	defer fl.Unlock()

	utils.Debug().Str(`chart_path`, cp).Send()
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}
	p := getter.All(ha.setting)
	vals, err := formatYamlValues(userValues)
	if err != nil {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: ha.setting.RepositoryConfig,
					RepositoryCache:  ha.setting.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = ha.namespace
	r, e := client.Run(chartRequested, vals)
	if e != nil {
		return nil, e
	}
	res := Release(*r)
	return &res, nil
}

type Release release.Release

// ListRelease 列出某空间下的releases
func (ha *helmAdapter) ListRelease() ([]*Release, error) {
	client := action.NewList(ha.cfg)
	client.All = true
	client.SetStateMask()
	rs, e := client.Run()
	if e != nil {
		return nil, e
	}
	res := make([]*Release, len(rs))
	for i, r := range rs {
		t := Release(*r)
		res[i] = &t
	}
	return res, nil
}

// GetValues 返回指定release 的values, json 格式
func (ha *helmAdapter) GetValues(relName string) (string, map[string]interface{}, error) {
	client := action.NewGetValues(ha.cfg)
	client.AllValues = true
	vals, e := client.Run(relName)
	if e != nil {
		return "", vals, e
	}
	ym, e := yaml.Marshal(vals)
	return string(ym), vals, e
}

// UninstallRelease 卸载指定的release
func (ha *helmAdapter) UninstallRelease(relName string) error {
	client := action.NewUninstall(ha.cfg)
	res, err := client.Run(relName)
	if err != nil {
		return err
	}
	if res != nil && res.Info != "" {
		utils.Error().Msg(res.Info)
	}
	utils.Info().Msgf("release \"%s\" uninstalled\n", relName)
	return nil
}

// UpgradeRelease 更新指定release 的版本或values
func (ha *helmAdapter) UpgradeRelease(
	releaseName, repoName, version string, userValues []byte,
	dryRun bool, maxhistory int, resetValue bool) (*Release, error) {
	client := action.NewUpgrade(ha.cfg)
	client.Namespace = ha.namespace
	client.Version = version
	client.DryRun = dryRun
	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	vals, e := formatYamlValues(userValues)
	if e != nil {
		return nil, e
	}
	chartPath := repoName

	cp, e := client.ChartPathOptions.LocateChart(chartPath, ha.setting)
	if e != nil {
		return nil, e
	}
	fl := utils.NewFileLock(cp)
	e = fl.LockWithTimeout(time.Second * 30)
	if e != nil {
		return nil, e
	}
	defer fl.Unlock()

	ch, e := loader.Load(cp)
	if e != nil {
		return nil, e
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, e
		}
	}

	client.ResetValues = resetValue
	client.ReuseValues = true
	client.MaxHistory = maxhistory
	rel, err := client.Run(releaseName, ch, vals)
	if err != nil {
		return nil, errors.Wrap(err, "UPGRADE FAILED")
	}
	res := Release(*rel)

	return &res, nil
}

// RollbackRelease 回滚release
func (ha *helmAdapter) RollbackRelease(releaseName string, ver int) error {
	client := action.NewRollback(ha.cfg)
	client.Version = ver
	return client.Run(releaseName)
}

// MergeValues 合并yaml 值， 如果同一key上值冲突会用v2 覆盖v1 的值
func MergeValues(v1, v2 []byte) (v3 []byte, e3 error) {
	mv1, e3 := formatYamlValues(v1)
	if e3 != nil {
		return
	}
	mv2, e3 := formatYamlValues(v2)
	if e3 != nil {
		return
	}
	mv3 := mergeMaps(mv1, mv2)

	v3, e3 = yaml.Marshal(mv3)
	return
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func formatYamlValues(raw []byte) (vals map[string]interface{}, err error) {
	vals = make(map[string]interface{})
	err = yaml.Unmarshal(raw, &vals)
	if err != nil {
		return nil, err
	}
	yjs, err := yaml.Marshal(vals)
	if err != nil {
		return nil, err
	}
	vjs, err := utils.YamlToJson(yjs)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(vjs, &vals)
	if err != nil {
		return nil, err
	}

	return
}

func (ha *helmAdapter) GetRelease(name string) (*Release, error) {
	c := action.NewGet(ha.cfg)
	r, e := c.Run(name)
	if e != nil {
		return nil, e
	}
	res := Release(*r)
	return &res, nil
}

func (ha *helmAdapter) GetHistory(name string) ([]*Release, error) {
	c := action.NewHistory(ha.cfg)
	rs, e := c.Run(name)
	if e != nil {
		return nil, e
	}
	res := make([]*Release, len(rs))
	for i, r := range rs {
		t := Release(*r)
		res[i] = &t
	}
	return res, nil
}

type ChartInfo struct {
	Name         string `json:"name"`
	RepoName     string `json:"repo_name"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
	Icon         string `json:"icon"`
	Description  string `json:"description"`
}

func (ci *ChartInfo) Convert() *protos.ChartInfo {
	return &protos.ChartInfo{
		Name:         ci.Name,
		RepoName:     ci.RepoName,
		ChartVersion: ci.ChartVersion,
		AppVersion:   ci.AppVersion,
		Description:  ci.Description,
		Icon:         ci.Icon,
	}
}

func (ha *helmAdapter) Search(keyword string, regexp bool) (charts []*ChartInfo, e error) {

	repoFile := ha.setting.RepositoryConfig
	rf, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(rf.Repositories) == 0 {
		return nil, errors.New("no repositories configured")
	}

	i := search.NewIndex()
	for _, re := range rf.Repositories {
		n := re.Name
		f := filepath.Join(ha.setting.RepositoryCache, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			utils.Warn().Msgf("WARNING: Repo %q is corrupt or missing. Try 'helm repo update'.", n)
			continue
		}

		i.AddRepo(n, ind, true)
	}
	rs, e := i.Search(keyword, 25, regexp)
	if e != nil {
		return nil, e
	}

	for _, r := range rs {
		charts = append(charts, &ChartInfo{
			Name:         r.Chart.Name,
			RepoName:     r.Name,
			ChartVersion: r.Chart.Version,
			AppVersion:   r.Chart.AppVersion,
			Description:  r.Chart.Description,
			Icon:         r.Chart.Icon,
		})
	}
	return
}

func (ha *helmAdapter) All() (charts []*ChartInfo, e error) {

	repoFile := ha.setting.RepositoryConfig
	rf, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(rf.Repositories) == 0 {
		return nil, errors.New("no repositories configured")
	}

	i := search.NewIndex()
	for _, re := range rf.Repositories {
		n := re.Name
		f := filepath.Join(ha.setting.RepositoryCache, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			utils.Warn().Msgf("WARNING: Repo %q is corrupt or missing. Try 'helm repo update'.", n)
			continue
		}

		i.AddRepo(n, ind, true)
	}
	rs := i.All()

	for _, r := range rs {
		charts = append(charts, &ChartInfo{
			Name:         r.Chart.Name,
			RepoName:     r.Name,
			ChartVersion: r.Chart.Version,
			AppVersion:   r.Chart.AppVersion,
			Description:  r.Chart.Description,
			Icon:         r.Chart.Icon,
		})
	}
	return
}

func loadChart(chartPath, version string,
	settings *cli.EnvSettings) (*chart.Chart, error) {
	client := action.NewInstall(nil)
	client.Version = version
	cp, err := client.ChartPathOptions.LocateChart(chartPath, settings)
	if err != nil {
		return nil, err
	}

	return LoadWithRawValues(cp)
}

func LoadWithRawValues(cp string) (c *chart.Chart, e error) {
	c, e = loader.Load(cp)
	if e != nil {
		return nil, e
	}
	if fi, err := os.Stat(cp); err != nil {
		return nil, err
	} else if fi.IsDir() {
		return nil, errors.New("cannot load a directory")
	}

	raw, e := os.Open(cp)
	if e != nil {
		return nil, e
	}
	defer raw.Close()

	bfs, err := loader.LoadArchiveFiles(raw)
	if err != nil {
		if err == gzip.ErrHeader {
			return nil, fmt.Errorf("file '%s' does not appear to be a valid chart file (details: %s)", cp, err)
		}
	}
	for _, bf := range bfs {
		if bf.Name == `values.yaml` {
			c.Files = append(c.Files, &chart.File{
				Name: bf.Name,
				Data: bf.Data,
			})
			break
		}
	}
	return c, e
}

// ToStruct marshal interface to pstruct.Struct
func ToStruct(msg interface{}) (pbs *pstruct.Struct, e error) {

	byteArray, e := json.Marshal(msg)
	if e != nil {
		return
	}
	pbs = &pstruct.Struct{}
	e = jsonpb.Unmarshal(bytes.NewReader(byteArray), pbs)
	ms := jsonpb.Marshaler{}
	x, e := ms.MarshalToString(pbs)
	utils.Debug().Interface(`?`, x).Send()
	return
}

type Maintainer chart.Maintainer

func (mt *Maintainer) Convert() (res *protos.Maintainer, e error) {
	if mt == nil {
		e = NULL_POINT_ERROR
		return
	}
	res = &protos.Maintainer{
		Name:  mt.Name,
		Email: mt.Email,
		Url:   mt.URL,
	}
	return
}

type Dependency chart.Dependency

// Convert to protos for Dependency
func (d *Dependency) Convert() (res *protos.Dependency, e error) {
	if d == nil {
		e = NULL_POINT_ERROR
		return
	}
	yamlVal, e := yaml.Marshal(d.ImportValues)

	res = &protos.Dependency{
		Name:         d.Name,
		Version:      d.Version,
		Repository:   d.Repository,
		Condition:    d.Condition,
		Tags:         d.Tags,
		Enabled:      d.Enabled,
		ImportValues: string(yamlVal),
		Alias:        d.Alias,
	}
	return
}

// Convert to protos for Chart
func (md *Metadata) Convert() (res *protos.Metadata, e error) {
	if md == nil {
		e = NULL_POINT_ERROR
		return
	}
	maintainers := make([]*protos.Maintainer, len(md.Maintainers))
	for i, m := range md.Maintainers {
		t := Maintainer(*m)
		pm, e := t.Convert()
		if e != nil {
			return nil, e
		}
		maintainers[i] = pm
	}
	dependencies := make([]*protos.Dependency, len(md.Dependencies))
	for i, d := range md.Dependencies {
		t := Dependency(*d)
		pd, e := t.Convert()
		if e != nil {
			return nil, e
		}
		dependencies[i] = pd
	}
	res = &protos.Metadata{
		Name:         md.Name,
		Home:         md.Home,
		Sources:      md.Sources,
		Version:      md.Version,
		Description:  md.Description,
		Keywords:     md.Keywords,
		Maintainers:  maintainers,
		Icon:         md.Icon,
		ApiVersion:   md.APIVersion,
		Condition:    md.Condition,
		Tags:         md.Tags,
		AppVersion:   md.AppVersion,
		Deprecated:   md.Deprecated,
		Annotations:  md.Annotations,
		KubeVersion:  md.KubeVersion,
		Dependencies: dependencies,
		Type:         md.Type,
	}
	return
}

type File chart.File

// Convert to protos for File
func (f *File) Convert() (res *protos.File, e error) {
	if f == nil {
		e = NULL_POINT_ERROR
		return
	}
	fbody := `not vaild utf8 files`
	if utf8.Valid(f.Data) {
		fbody = string(f.Data)
	}
	res = &protos.File{
		Name: f.Name,
		Data: fbody,
	}
	return
}

type Hook release.Hook

// Convert to protos for Hook
func (h *Hook) Convert() (res *protos.Hook, e error) {
	if h == nil {
		e = NULL_POINT_ERROR
		return
	}
	events := make([]string, len(h.Events))
	for i, he := range h.Events {
		events[i] = string(he)
	}
	startedAt, e := ptypes.TimestampProto(h.LastRun.StartedAt.Time)
	if e != nil {
		return
	}
	completedAt, e := ptypes.TimestampProto(h.LastRun.CompletedAt.Time)
	if e != nil {
		return
	}
	deletePolicies := make([]string, len(h.DeletePolicies))
	for i, dp := range h.DeletePolicies {
		deletePolicies[i] = string(dp)
	}
	res = &protos.Hook{
		Name:     h.Name,
		Kind:     h.Kind,
		Path:     h.Path,
		Manifest: h.Manifest,
		Events:   events,
		LastRun: &protos.HookExecution{
			StartedAt:   startedAt,
			CompletedAt: completedAt,
			Phase:       string(h.LastRun.Phase),
		},
		Weight:         int64(h.Weight),
		DeletePolicies: deletePolicies,
	}
	return
}

// Convert to protos for Chart
func (c *Chart) Convert() (res *protos.Chart, e error) {
	if c == nil {
		e = NULL_POINT_ERROR
		return
	}
	metadata := Metadata(*c.Metadata)
	md, e := metadata.Convert()
	if e != nil {
		return
	}
	templates := make([]*protos.File, len(c.Templates))
	for i, tf := range c.Templates {
		t := File(*tf)
		pf, e := t.Convert()
		if e != nil {
			return nil, e
		}
		templates[i] = pf
	}

	yamlVal, e := yaml.Marshal(c.Values)
	files := make([]*protos.File, len(c.Files))
	for i, tf := range c.Files {
		t := File(*tf)
		pf, e := t.Convert()
		if e != nil {
			return nil, e
		}
		files[i] = pf
	}

	res = &protos.Chart{
		Metadata:  md,
		Templates: templates,
		Values:    string(yamlVal),
		Schema:    string(c.Schema),
		Files:     files,
	}
	return
}

// Convert to protos for Release
func (rel *Release) Convert() (res *protos.Release, e error) {
	if rel == nil {
		e = NULL_POINT_ERROR
		return
	}
	firstDeployed, e := ptypes.TimestampProto(rel.Info.FirstDeployed.Time)
	if e != nil {
		return
	}
	lastDeployed, e := ptypes.TimestampProto(rel.Info.LastDeployed.Time)
	if e != nil {
		return
	}
	deleted, e := ptypes.TimestampProto(rel.Info.Deleted.Time)
	if e != nil {
		return
	}
	ch := Chart(*rel.Chart)
	pc, e := ch.Convert()
	if e != nil {
		return
	}

	yamlConf, e := yaml.Marshal(rel.Config)
	hooks := make([]*protos.Hook, len(rel.Hooks))
	for i, h := range rel.Hooks {
		t := Hook(*h)
		ph, e := t.Convert()
		if e != nil {
			return nil, e
		}
		hooks[i] = ph
	}
	res = &protos.Release{
		Name: rel.Name,
		Info: &protos.Info{
			FirstDeployed: firstDeployed,
			LastDeployed:  lastDeployed,
			Deleted:       deleted,
			Description:   rel.Info.Description,
			Status:        rel.Info.Status.String(),
			Notes:         rel.Info.Notes,
		},
		Chart:     pc,
		Config:    string(yamlConf),
		Manifest:  rel.Manifest,
		Hooks:     hooks,
		Version:   int64(rel.Version),
		Namespace: rel.Namespace,
	}
	return
}
