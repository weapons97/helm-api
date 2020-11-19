# helm-api

通过 grpc 和 http 接口部署 helm charts (仅支持3.0.0+)

### 为什么要使用helm-api

- 有时候您希望通过程序在集群上部署和管理helm release。
- 有时候您希望通过http接口在集群上部署和管理helm release。

## 实现接口

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

### Key 关键词
1. Context 
   * context 存储了k8s权限和helm repository。所以用户可以通过context 复用这些信息。
2. ContextName
   * 如果用户请求中带有ContextName 请求会使用对应context 存储的k8s权限和helm repository。如果请求不带有contextName，请求会根据RepoInfo 和 KubeInfo创建一个临时的context。

### Installation
* git clone git@github.com:weapons97/helm-api.git
* cd helm-api
* kubectl install -f deploy

### DOC
请看 doc/doc.md

### HTTP
请看 swagger/helm-api.swagger.json

### GRPC
请看 protos/helm-api.proto

### Examples

cd ./examples 打开此目录，运行一下脚本以验证接口。

* 确保你安装了 nodejs 环境。

0. 指定ip和端口
   ```bash
   export host=<你的 helm-api host>
   export port=<你的 helm-api port>
   ```
1. 创建context。

   ``` bash
   # node createContext.js <your kubeconfig> <your context> <your namespace>
   node createContext.js ~/kube.config "" "test"
   ```

   创建 incluster context。 incluster context 使用service account 为操作授权。

   Create an incluster context. Incluster context uses service account to authorize operations.

   ```bash
   # node inclusterContext.js <your namespace>
   node inclusterContext.js default
   ```

2. 更新 repo

   ```bash
   node repoupdate.js 
   ```

3. 获得 chart 信息

   ```bash
   #  node getchart.js test <repo name> <chart name>
   node getchart.js test bitnami mariadb
   ```

4. 创建一个mysql release

   ```bash
   # node install.js <release name>
   node install.js mysql
   ```

5. 获取mysql release values

   ```bash
   # node getreleasevalues.js <context name> <release name>
   node getreleasevalues.js test mysql
   ```

6. 更新mysql release 

   ```bash
   # node upgrade <release name>
   node upgrade mysql
   ```

7. 获取mysql release all

   ```bash
   # node getreleasevalues.js <context name> <release name>
   node getreleasevalues.js test mysql
   ```

8. 获取 release history

   ```bash
   # node history.js <context name> <release name>
   node history.js test mysql 
   ```

9. 回滚 release 

   ```bash
   # node rollback.js <context name> <release name> <version>
   node rollback.js test mysql 1
   ```

10. 删除 release

   ```bash
   # node deleterelease.js <context name> <release name>
   node deleterelease.js test mysql
   ```

11. 删除context

    ```bash
    # node deleteContext.js <context name>
    node deleteContext.js test
    ```

    

