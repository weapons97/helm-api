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
package services

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/weapons97/helm-api/pkg/helm"
	"github.com/weapons97/helm-api/pkg/helm/adapter"
	protos "github.com/weapons97/helm-api/pkg/protos"
	"github.com/weapons97/helm-api/pkg/utils"
	"time"
)

var NotHelmFromContext = fmt.Errorf(`can't find helm client from context'`)

func HelmFromContext(ctx context.Context) (adapter.Adapter, bool) {
	u, ok := ctx.Value(helm.HelmClientInContextName).(adapter.Adapter)
	return u, ok
}

type HelmApiService struct {
}

// CreateHelmContext
func (has *HelmApiService) CreateContext(ctx context.Context, creq *protos.HelmContextReq) (hr *protos.HelmContextRes, e error) {

	entrys := creq.GetRepoinfo().GetEntrys()
	rentrys := make([]adapter.Entry, len(entrys))
	for i, pe := range entrys {
		entry := adapter.ConvertEntry(pe)
		rentrys[i] = entry
	}
	var expiryDate time.Time

	if creq.GetExpiry() == 0 {
		expiryDate = time.Date(2222, time.February, 2, 2, 2, 22, 222, time.Local)
	} else {
		expiryDate = time.Now().Add(time.Second * time.Duration(creq.GetExpiry()))
	}
	hc, e := helm.NewHelmContext(creq.GetName(), creq.GetIncluster(),
		[]byte(creq.GetKubeinfo().GetKubeconfig()), creq.GetKubeinfo().GetContext(), creq.GetKubeinfo().GetNamespace(),
		expiryDate,
		rentrys...,
	)
	if e != nil {
		utils.Error().Err(e).Send()
		helm.DelHelmContext(creq.Name)
		return
	}
	hr = &protos.HelmContextRes{
		Name: hc.Name,
	}
	return
}

// DeleteHelmContext
func (has *HelmApiService) DeleteContext(ctx context.Context, dreq *protos.DeleteHelmContextReq) (*empty.Empty, error) {
	return &empty.Empty{}, helm.DelHelmContext(dreq.GetName())
}

// UpdateRepo
func (has *HelmApiService) UpdateRepo(ctx context.Context, req *protos.UpdateRepoReq) (emp *empty.Empty, e error) {
	emp = &empty.Empty{}
	utils.Debug().Interface(`req`, req).Send()
	hc, e := helm.LoadHelmContext(req.GetContextName())
	if e != nil {
		return
	}
	utils.Debug().Interface(`hc`, hc).Send()
	entrys := req.GetRepoinfo().GetEntrys()
	rentrys := make([]adapter.Entry, len(entrys))
	for i, pe := range entrys {
		entry := adapter.ConvertEntry(pe)
		rentrys[i] = entry
	}
	e = hc.UpdateRepo(rentrys...)
	return
}

// GetChart 如果chart 里面含有非文本数据会返回错误
func (has *HelmApiService) GetChart(ctx context.Context, gcr *protos.GetChartReq) (c *protos.Chart, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	rc, e := client.GetChartDetail(
		fmt.Sprintf(`%s/%s`, gcr.GetRepo(), gcr.GetChartName()),
		gcr.GetChartVersion())
	if e != nil {
		return
	}
	c, e = rc.Convert()
	return
}

// GetRelease 如果release 对应chart 里面含有非文本数据会返回错误
func (has *HelmApiService) GetRelease(ctx context.Context, grr *protos.ReleaseReq) (r *protos.Release, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	pr, e := client.GetRelease(grr.GetReleaseName())
	if e != nil {
		return
	}
	r, e = pr.Convert()
	return
}

// GetRelease 如果release 对应chart 里面含有非文本数据会返回错误
func (has *HelmApiService) GetReleaseValues(ctx context.Context, grr *protos.ReleaseReq) (v *protos.Values, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	pr, _, e := client.GetValues(grr.GetReleaseName())
	v = &protos.Values{
		Yaml: pr,
	}
	return
}

// GetRelease 如果release 对应chart 里面含有非文本数据会返回错误
func (has *HelmApiService) ListRelease(ctx context.Context, lr *protos.ListReleaseReq) (lres *protos.ListReleaseRes, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	rels, e := client.ListRelease()
	if e != nil {
		return
	}
	lres = &protos.ListReleaseRes{
		Releases: make([]*protos.Release, len(rels)),
	}
	for i, r := range rels {
		lres.Releases[i], e = r.Convert()
		if e != nil {
			return nil, e
		}

	}
	return
}

// InstallCharts
func (has *HelmApiService) InstallRelease(ctx context.Context, ir *protos.InstallReq) (r *protos.Release, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}

	rel, e := client.InstallCharts(
		ir.GetReleaseName(),
		ir.GetRepoChartName(),
		ir.GetChartVersion(),
		[]byte(ir.GetValues()),
		ir.GetDryRun(),
	)
	if e != nil {
		return
	}
	r, e = rel.Convert()
	return
}

// Upgrade
func (has *HelmApiService) UpgradeRelease(ctx context.Context, ur *protos.UpgradeReq) (r *protos.Release, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}

	rel, e := client.UpgradeRelease(
		ur.GetReleaseName(),
		ur.GetRepoChartName(),
		ur.GetChartVersion(),
		[]byte(ur.GetValues()),
		ur.GetDryRun(),
		int(ur.GetHistoryMax()),
		ur.GetResetValues(),
	)
	if e != nil {
		return
	}
	r, e = rel.Convert()
	return
}

// RollbackRelease
func (has *HelmApiService) RollbackRelease(ctx context.Context, rr *protos.ReleaseRollbackReq) (emp *empty.Empty, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	emp = &empty.Empty{}
	e = client.RollbackRelease(rr.GetReleaseName(), int(rr.GetReversion()))
	return
}

// UninstallRelease
func (has *HelmApiService) UninstallRelease(ctx context.Context, rr *protos.ReleaseReq) (emp *empty.Empty, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	emp = &empty.Empty{}
	e = client.UninstallRelease(rr.GetReleaseName())
	return
}

// GetReleaseHistory
func (has *HelmApiService) GetReleaseHistory(ctx context.Context, rr *protos.ReleaseReq) (history *protos.ListReleaseRes, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	res, e := client.GetHistory(rr.GetReleaseName())
	history = &protos.ListReleaseRes{
		Releases: make([]*protos.Release, len(res)),
	}
	for i, r := range res {
		history.Releases[i], e = r.Convert()
		if e != nil {
			return
		}
	}
	return
}

// Search
func (has *HelmApiService) Search(ctx context.Context, sr *protos.SearchReq) (res *protos.SearchRes, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	cx, e := client.Search(sr.GetKeyword(), sr.GetRegexp())
	if e != nil {
		return
	}
	charts := make([]*protos.ChartInfo, len(cx))
	for i, c := range cx {
		charts[i] = c.Convert()
	}
	res = &protos.SearchRes{
		Charts: charts,
	}
	return
}

// All
func (has *HelmApiService) All(ctx context.Context, lcr *protos.ListChartReq) (res *protos.SearchRes, e error) {
	client, ok := HelmFromContext(ctx)
	if !ok {
		return nil, NotHelmFromContext
	}
	cx, e := client.All()
	if e != nil {
		return
	}
	charts := make([]*protos.ChartInfo, len(cx))
	for i, c := range cx {
		charts[i] = c.Convert()
	}
	res = &protos.SearchRes{
		Charts: charts,
	}
	return
}
