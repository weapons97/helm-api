# helm-api

通过grpc 和 http 接口部署 helm charts (仅支持3.0.0+)
deploy your helm charts for Kubernetes via GRPC AND HTTP endpoints （only for 3.0.0+）

### When to use this solution

- Automating deployments in the cluster.
- Programmatically managing the cluster from the code.

## Coverage

HelmApiService

| Method Name       | Request Type                                          | Response Type                                    | Description                                                  |
| ----------------- | ----------------------------------------------------- | ------------------------------------------------ | ------------------------------------------------------------ |
| CreateContext     | [HelmContextReq](#helmapi.HelmContextReq)             | [HelmContextRes](#helmapi.HelmContextRes)        | CreateContext 创建context context 持有了k8s集群资源操作权限和harbor登录信息。所以这些信息可以通过context复用。 The context holds k8s cluster resource operation authority and harbor login information. So this information can be reuse through context. |
| DeleteContext     | [DeleteHelmContextReq](#helmapi.DeleteHelmContextReq) | [.google.protobuf.Empty](#google.protobuf.Empty) | DeleteContext 删除context                                    |
| UpdateRepo        | [UpdateRepoReq](#helmapi.UpdateRepoReq)               | [.google.protobuf.Empty](#google.protobuf.Empty) | UpdateRepo 更新context 内repo 信息                           |
| InstallRelease    | [InstallReq](#helmapi.InstallReq)                     | [Release](#helmapi.Release)                      | InstallRelease 安装charts                                    |
| UpgradeRelease    | [UpgradeReq](#helmapi.UpgradeReq)                     | [Release](#helmapi.Release)                      | UpgradeRelease 更新release                                   |
| UninstallRelease  | [ReleaseReq](#helmapi.ReleaseReq)                     | [.google.protobuf.Empty](#google.protobuf.Empty) | UninstallRelease 删除release                                 |
| GetChart          | [GetChartReq](#helmapi.GetChartReq)                   | [Chart](#helmapi.Chart)                          | GetChart 获取 chart 信息                                     |
| GetRelease        | [ReleaseReq](#helmapi.ReleaseReq)                     | [Release](#helmapi.Release)                      | GetRelease 获取某个release实例信息                           |
| GetReleaseValues  | [ReleaseReq](#helmapi.ReleaseReq)                     | [Values](#helmapi.Values)                        | GetReleaseValues 某个release实例values信息                   |
| ListRelease       | [ListReleaseReq](#helmapi.ListReleaseReq)             | [ListReleaseRes](#helmapi.ListReleaseRes)        | ListRelease 列出某个context下全部release.                    |
| RollbackRelease   | [ReleaseRollbackReq](#helmapi.ReleaseRollbackReq)     | [.google.protobuf.Empty](#google.protobuf.Empty) | RollbackRelease 回滚某个release                              |
| GetReleaseHistory | [ReleaseReq](#helmapi.ReleaseReq)                     | [ListReleaseRes](#helmapi.ListReleaseRes)        | GetReleaseHistory 列出release 历史                           |



### Key
1. Context 
   * context 存储了k8s权限和helm repository。所以用户可以通过context 复用这些信息。
   * The context stores k8s permissions and helm repository. So users can reuse this information through context.
2. contextName
   * 如果用户请求中带有contextName 请求会使用对应context 存储的k8s权限和helm repository。如果请求不带有contextName，请求会根据RepoInfo 和 KubeInfo创建一个临时的context。
   * If the user request contains a contextName request, the k8s permission and helm repository stored in the corresponding context will be used. If the request does not have a contextName, the request will create a temporary context based on RepoInfo and KubeInfo.

### Installation
* kubectl install -f deploy

### DOC
just look at doc/doc.md

### HTTP
just look at swagger/helm-api.swagger.json

### GRPC
just look at protos/helm-api.proto

### Examples

cd ./examples 打开此目录，运行一下脚本以验证接口。

* 确保你安装了 nodejs 环境。

install nodejs

```bash
# ubuntu
sudo apt-get install nodejs
# mac
brew install node
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

    

