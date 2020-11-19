# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [helm-api.proto](#helm-api.proto)
    - [Chart](#helmapi.Chart)
    - [ChartInfo](#helmapi.ChartInfo)
    - [DeleteHelmContextReq](#helmapi.DeleteHelmContextReq)
    - [Dependency](#helmapi.Dependency)
    - [Entry](#helmapi.Entry)
    - [File](#helmapi.File)
    - [GetChartReq](#helmapi.GetChartReq)
    - [HelmContextReq](#helmapi.HelmContextReq)
    - [HelmContextRes](#helmapi.HelmContextRes)
    - [Hook](#helmapi.Hook)
    - [HookExecution](#helmapi.HookExecution)
    - [Info](#helmapi.Info)
    - [InstallReq](#helmapi.InstallReq)
    - [KubeInfo](#helmapi.KubeInfo)
    - [ListChartReq](#helmapi.ListChartReq)
    - [ListReleaseReq](#helmapi.ListReleaseReq)
    - [ListReleaseRes](#helmapi.ListReleaseRes)
    - [Maintainer](#helmapi.Maintainer)
    - [Metadata](#helmapi.Metadata)
    - [Metadata.AnnotationsEntry](#helmapi.Metadata.AnnotationsEntry)
    - [Release](#helmapi.Release)
    - [ReleaseReq](#helmapi.ReleaseReq)
    - [ReleaseRollbackReq](#helmapi.ReleaseRollbackReq)
    - [RepoInfo](#helmapi.RepoInfo)
    - [SearchReq](#helmapi.SearchReq)
    - [SearchRes](#helmapi.SearchRes)
    - [UpdateRepoReq](#helmapi.UpdateRepoReq)
    - [UpgradeReq](#helmapi.UpgradeReq)
    - [Values](#helmapi.Values)
  
    - [HelmApiService](#helmapi.HelmApiService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="helm-api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## helm-api.proto



<a name="helmapi.Chart"></a>

### Chart
Chart is a helm package that contains metadata, a default config, zero or more
optionally parameterizable templates, and zero or more charts (dependencies).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metadata | [Metadata](#helmapi.Metadata) |  | Metadata is the contents of the Chartfile. |
| templates | [File](#helmapi.File) | repeated | Templates for this chart. |
| values | [string](#string) |  | Values are default config for this chart. |
| schema | [string](#string) |  | Schema is an optional JSON schema for imposing structure on Values |
| files | [File](#helmapi.File) | repeated | Files are miscellaneous files in a chart archive, e.g. README, LICENSE, etc. |






<a name="helmapi.ChartInfo"></a>

### ChartInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| repo_name | [string](#string) |  |  |
| chart_version | [string](#string) |  |  |
| app_version | [string](#string) |  |  |
| description | [string](#string) |  |  |
| icon | [string](#string) |  |  |






<a name="helmapi.DeleteHelmContextReq"></a>

### DeleteHelmContextReq
DeleteHelmContextReq


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name 唯一的context 标识。 |






<a name="helmapi.Dependency"></a>

### Dependency
Dependency describes a chart upon which another chart depends.

Dependencies can be used to express developer intent, or to capture the state
of a chart.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is the name of the dependency.

This must mach the name in the dependency&#39;s Chart.yaml. |
| version | [string](#string) |  | Version is the version (range) of this chart.

A lock file will always produce a single version, while a dependency may contain a semantic version range. |
| repository | [string](#string) |  | The URL to the repository.

Appending `index.yaml` to this string should result in a URL that can be used to fetch the repository index. |
| condition | [string](#string) |  | A yaml path that resolves to a boolean, used for enabling/disabling charts (e.g. subchart1.enabled ) |
| tags | [string](#string) | repeated | Tags can be used to group charts for enabling/disabling together |
| enabled | [bool](#bool) |  | Enabled bool determines if chart should be loaded |
| import_values | [string](#string) |  | ImportValues holds the mapping of source values to parent key to be imported. Each item can be a string or pair of child/parent sublist items. |
| alias | [string](#string) |  | Alias usable alias to be used for the chart |






<a name="helmapi.Entry"></a>

### Entry
Entry 提供 chart repository 信息


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| url | [string](#string) |  |  |
| username | [string](#string) |  |  |
| password | [string](#string) |  |  |
| certFile | [string](#string) |  |  |
| keyFile | [string](#string) |  |  |
| caFile | [string](#string) |  |  |
| insecure_skip_tls_verify | [bool](#bool) |  |  |






<a name="helmapi.File"></a>

### File



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is the path-like name of the template. |
| data | [string](#string) |  | Data is the template as byte data. |






<a name="helmapi.GetChartReq"></a>

### GetChartReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| repo | [string](#string) |  |  |
| chartName | [string](#string) |  |  |
| chartVersion | [string](#string) |  |  |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo 如果没有contextName，但是有RepoInfo会创建一个临时的context。 |






<a name="helmapi.HelmContextReq"></a>

### HelmContextReq
HelmContextReq


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name 唯一的context 标识。冲突则会覆盖。 |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 提供。 |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo 提供 chart repository 相关信息。 |
| incluster | [bool](#bool) |  | incluster 为true 时，会用sa 为helm-api 提供k8s的访问权限。 |
| expiry | [int64](#int64) |  | expiry 是超时的时间戳，超过此值，context会被回收。 |






<a name="helmapi.HelmContextRes"></a>

### HelmContextRes
HelmContextRes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="helmapi.Hook"></a>

### Hook
Hook defines a hook object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| kind | [string](#string) |  | Kind is the Kubernetes kind. |
| path | [string](#string) |  | Path is the chart-relative path to the template. |
| manifest | [string](#string) |  | Manifest is the manifest contents. |
| events | [string](#string) | repeated | Events are the events that this hook fires on. |
| last_run | [HookExecution](#helmapi.HookExecution) |  | LastRun indicates the date/time this was last run. |
| weight | [int64](#int64) |  | Weight indicates the sort order for execution among similar Hook type |
| delete_policies | [string](#string) | repeated | DeletePolicies are the policies that indicate when to delete the hook |






<a name="helmapi.HookExecution"></a>

### HookExecution
A HookExecution records the result for the last execution of a hook for a given release.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| started_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | StartedAt indicates the date/time this hook was started |
| completed_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | CompletedAt indicates the date/time this hook was completed. |
| phase | [string](#string) |  | Phase indicates whether the hook completed successfully |






<a name="helmapi.Info"></a>

### Info
Info describes release information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| first_deployed | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | FirstDeployed is when the release was first deployed. |
| last_deployed | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | LastDeployed is when the release was last deployed. |
| deleted | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Deleted tracks when this object was deleted. |
| description | [string](#string) |  | Description is human-friendly &#34;log entry&#34; about this release. |
| status | [string](#string) |  | Status is the current state of the release |
| notes | [string](#string) |  | Contains the rendered templates/NOTES.txt if available |






<a name="helmapi.InstallReq"></a>

### InstallReq
InstallReq 安装chart的请求


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| repoChartName | [string](#string) |  |  |
| chartVersion | [string](#string) |  | chartVersion 为空会使用最新版本（可选） |
| namespace | [string](#string) |  | namespace 为空会使用context 默认的空间 （可选） |
| values | [string](#string) |  | values 是yaml格式的value文件的内容。（可选） |
| releaseName | [string](#string) |  |  |
| dry_run | [bool](#bool) |  | dry_run 为true 不会真实的部署release。 |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |






<a name="helmapi.KubeInfo"></a>

### KubeInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kubeconfig | [string](#string) |  | kubeconfig kubeconfig 的内容，你可以在~/.kube.config下找到它。 |
| context | [string](#string) |  | context 指定kubeconfig中哪个context被使用。 |
| namespace | [string](#string) |  | namespace 指定使用的工作空间。 |






<a name="helmapi.ListChartReq"></a>

### ListChartReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName is name of CreateHelmContext return |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo or contextName |






<a name="helmapi.ListReleaseReq"></a>

### ListReleaseReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| namespace | [string](#string) |  | namespace 为空会使用context 默认的空间 （可选） |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 如果没有contextName，但是有KubeInfo会创建一个临时的context。 |






<a name="helmapi.ListReleaseRes"></a>

### ListReleaseRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| releases | [Release](#helmapi.Release) | repeated |  |






<a name="helmapi.Maintainer"></a>

### Maintainer
Maintainer describes a Chart maintainer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a user name or organization name |
| email | [string](#string) |  | Email is an optional email address to contact the named maintainer |
| url | [string](#string) |  | URL is an optional URL to an address for the named maintainer |






<a name="helmapi.Metadata"></a>

### Metadata
Metadata for a Chart file. This models the structure of a Chart.yaml file.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the chart |
| home | [string](#string) |  | The URL to a relevant project page, git repo, or contact person |
| sources | [string](#string) | repeated | Source is the URL to the source code of this chart |
| version | [string](#string) |  | A SemVer 2 conformant version string of the chart |
| description | [string](#string) |  | A one-sentence description of the chart |
| keywords | [string](#string) | repeated | A list of string keywords |
| maintainers | [Maintainer](#helmapi.Maintainer) | repeated | A list of name and URL/email address combinations for the maintainer(s) |
| icon | [string](#string) |  | The URL to an icon file. |
| api_version | [string](#string) |  | The API Version of this chart. |
| condition | [string](#string) |  | The condition to check to enable chart |
| tags | [string](#string) |  | The tags to check to enable chart |
| app_version | [string](#string) |  | The version of the application enclosed inside of this chart. |
| deprecated | [bool](#bool) |  | Whether or not this chart is deprecated |
| annotations | [Metadata.AnnotationsEntry](#helmapi.Metadata.AnnotationsEntry) | repeated | Annotations are additional mappings uninterpreted by Helm, made available for inspection by other applications. |
| kubeVersion | [string](#string) |  | KubeVersion is a SemVer constraint specifying the version of Kubernetes required. |
| dependencies | [Dependency](#helmapi.Dependency) | repeated | Dependencies are a list of dependencies for a chart. |
| type | [string](#string) |  | Specifies the chart type: application or library |






<a name="helmapi.Metadata.AnnotationsEntry"></a>

### Metadata.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="helmapi.Release"></a>

### Release
Release describes a deployment of a chart, together with the chart
and the variables used to deploy that chart.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is the name of the release |
| info | [Info](#helmapi.Info) |  | Info provides information about a release |
| chart | [Chart](#helmapi.Chart) |  | Chart is the chart that was released. |
| config | [string](#string) |  | Config is the set of extra Values added to the chart. These values override the default values inside of the chart. |
| manifest | [string](#string) |  | Manifest is the string representation of the rendered template. |
| hooks | [Hook](#helmapi.Hook) | repeated | Hooks are all of the hooks declared for this release. |
| version | [int64](#int64) |  | Version is an int which represents the version of the release. |
| namespace | [string](#string) |  | Namespace is the kubernetes namespace of the release. |






<a name="helmapi.ReleaseReq"></a>

### ReleaseReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| releaseName | [string](#string) |  |  |
| namespace | [string](#string) |  | namespace 为空会使用context 默认的空间 （可选） |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 如果没有contextName，但是有KubeInfo会创建一个临时的context。 |






<a name="helmapi.ReleaseRollbackReq"></a>

### ReleaseRollbackReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| releaseName | [string](#string) |  |  |
| reversion | [int32](#int32) |  |  |
| namespace | [string](#string) |  | namespace 为空会使用context 默认的空间 （可选） |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |






<a name="helmapi.RepoInfo"></a>

### RepoInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entrys | [Entry](#helmapi.Entry) | repeated | entrys 提供 一组 chart repository 信息 |






<a name="helmapi.SearchReq"></a>

### SearchReq
SearchReq


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keyword | [string](#string) |  |  |
| regexp | [bool](#bool) |  | use regular expressions for searching repositories you have added |
| contextName | [string](#string) |  | contextName is name of CreateHelmContext return |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo or contextName |






<a name="helmapi.SearchRes"></a>

### SearchRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| charts | [ChartInfo](#helmapi.ChartInfo) | repeated |  |






<a name="helmapi.UpdateRepoReq"></a>

### UpdateRepoReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。 |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  |  |






<a name="helmapi.UpgradeReq"></a>

### UpgradeReq
UpgradeReq 升级release的请求。


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextName | [string](#string) |  | contextName 唯一的context 标识。（可选） |
| repoChartName | [string](#string) |  |  |
| chartVersion | [string](#string) |  |  |
| namespace | [string](#string) |  | namespace 为空会使用context 默认的空间 （可选） |
| values | [string](#string) |  | values 是yaml格式的value文件的内容。（可选） |
| releaseName | [string](#string) |  |  |
| dry_run | [bool](#bool) |  | dry_run 为true 不会真实的部署release。 |
| history_max | [int32](#int32) |  |  |
| reset_values | [bool](#bool) |  |  |
| kubeinfo | [KubeInfo](#helmapi.KubeInfo) |  | KubeInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |
| repoinfo | [RepoInfo](#helmapi.RepoInfo) |  | RepoInfo 如果没有contextName，但是有KubeInfo和RepoInfo会创建一个临时的context。 |






<a name="helmapi.Values"></a>

### Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| yaml | [string](#string) |  |  |





 

 

 


<a name="helmapi.HelmApiService"></a>

### HelmApiService
HelmApiService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateContext | [HelmContextReq](#helmapi.HelmContextReq) | [HelmContextRes](#helmapi.HelmContextRes) | CreateContext 创建context |
| DeleteContext | [DeleteHelmContextReq](#helmapi.DeleteHelmContextReq) | [.google.protobuf.Empty](#google.protobuf.Empty) | DeleteContext 删除context |
| UpdateRepo | [UpdateRepoReq](#helmapi.UpdateRepoReq) | [.google.protobuf.Empty](#google.protobuf.Empty) | UpdateRepo 更新context 内repo 信息 |
| InstallRelease | [InstallReq](#helmapi.InstallReq) | [Release](#helmapi.Release) | InstallRelease 安装charts |
| UpgradeRelease | [UpgradeReq](#helmapi.UpgradeReq) | [Release](#helmapi.Release) | UpgradeRelease 更新release |
| UninstallRelease | [ReleaseReq](#helmapi.ReleaseReq) | [.google.protobuf.Empty](#google.protobuf.Empty) | UninstallRelease 删除release |
| GetChart | [GetChartReq](#helmapi.GetChartReq) | [Chart](#helmapi.Chart) | GetChart 获取 chart 信息 |
| GetRelease | [ReleaseReq](#helmapi.ReleaseReq) | [Release](#helmapi.Release) | GetRelease 获取某个release实例信息 |
| GetReleaseValues | [ReleaseReq](#helmapi.ReleaseReq) | [Values](#helmapi.Values) | GetReleaseValues 某个release实例values信息 |
| ListRelease | [ListReleaseReq](#helmapi.ListReleaseReq) | [ListReleaseRes](#helmapi.ListReleaseRes) | ListRelease 列出某个context下全部release. |
| RollbackRelease | [ReleaseRollbackReq](#helmapi.ReleaseRollbackReq) | [.google.protobuf.Empty](#google.protobuf.Empty) | RollbackRelease 回滚某个release |
| GetReleaseHistory | [ReleaseReq](#helmapi.ReleaseReq) | [ListReleaseRes](#helmapi.ListReleaseRes) | GetReleaseHistory 列出release 历史 |
| Search | [SearchReq](#helmapi.SearchReq) | [SearchRes](#helmapi.SearchRes) | Search 查找某个repo的chart |
| All | [ListChartReq](#helmapi.ListChartReq) | [SearchRes](#helmapi.SearchRes) | All 列出某个context所有的chart |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

