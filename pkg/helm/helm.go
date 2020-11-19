package helm

import (
	"github.com/weapons97/helm-api/pkg/helm/adapter"
	"github.com/weapons97/helm-api/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	CONTEXT_CONFIG = `config.yaml`
	KUBECONFIG     = `kubeconfig.yaml`
	LOCKFILE       = `/helm-api.lock`
)

// HelmContext
type HelmContext struct {
	ContextDir                                 string
	Name                                       string
	KubeConfigPath, KubeContext, KubeNamespace string
	Incluster                                  bool
	Expiry                                     int64
}

// NewHelmConifg create helm config by kubeconfig
func (hc *HelmContext) NewHelmConifg() (actionConfig *action.Configuration, e error) {
	actionConfig = new(action.Configuration)

	kubeConfig := &genericclioptions.ConfigFlags{}
	if hc.Incluster {
		config, e := rest.InClusterConfig()
		if e != nil {
			return nil, e
		}
		kubeConfig := genericclioptions.NewConfigFlags(false)
		kubeConfig.APIServer = &config.Host
		kubeConfig.BearerToken = &config.BearerToken
		kubeConfig.CAFile = &config.CAFile
		kubeConfig.Namespace = &hc.KubeNamespace

	} else {
		kubeConfig = kube.GetConfig(hc.KubeConfigPath, hc.KubeContext, hc.KubeNamespace)
	}

	kc := kube.New(kubeConfig)
	kc.Log = func(s string, i ...interface{}) {
		utils.Info().Msgf(s, i...)
	}

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		utils.Error().Err(err).Msg(`can't create helm cfg'`)
		return nil, fmt.Errorf(`can't create helm cfg' bad kubeconfig`)
	}

	var store *storage.Storage
	// 固定用secrets 存储release
	d := driver.NewSecrets(clientset.CoreV1().Secrets(hc.KubeNamespace))
	d.Log = func(s string, i ...interface{}) {
		utils.Info().Msgf(s, i...)
	}
	store = storage.Init(d)

	actionConfig.RESTClientGetter = kubeConfig
	actionConfig.KubeClient = kc
	actionConfig.Releases = store
	actionConfig.Log = func(s string, i ...interface{}) {
		utils.Info().Msgf(s, i...)
	}
	return
}

// UpdateRepo same as command `helm repo update`
func (hc *HelmContext) UpdateRepo(entrys ...adapter.Entry) error {
	settings := hc.NewHelmSettings()
	repoFile := settings.RepositoryConfig
	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		utils.Error().Err(err).Send()
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		utils.Error().Err(err).Send()
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		utils.Error().Err(err).Send()
		return err
	}

	rentrys := make([]*repo.Entry, len(entrys))
	for i, entry := range entrys {
		_entry := repo.Entry(entry)
		r, err := repo.NewChartRepository(&_entry, getter.All(settings))
		if err != nil {
			utils.Error().Err(err).Send()
			return err
		}
		r.CachePath = settings.RepositoryCache

		utils.Info().Interface(`config`, r.Config).Msg(`DownloadIndexFile`)
		if _, err := r.DownloadIndexFile(); err != nil {
			utils.Error().Err(err).Send()
			return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", entry.URL)
		}
		rentrys[i] = &_entry
	}

	f.Update(rentrys...)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		utils.Error().Err(err).Send()
		return err
	}
	return nil
}

// NewHelmSettings create setting we need to create helm action
func (hc *HelmContext) NewHelmSettings() (settings *cli.EnvSettings) {
	settings = cli.New()
	settings.RepositoryConfig = hc.ContextDir + `/repositories.yaml`
	settings.RepositoryCache = hc.ContextDir + `/repository`
	return
}

// NewHelmClient
func (hc *HelmContext) NewHelmClient() (adapter.Adapter, error) {
	cfg, e := hc.NewHelmConifg()
	if e != nil {
		return nil, e
	}
	settings := hc.NewHelmSettings()
	return adapter.NewAdapter(cfg, settings, hc.KubeNamespace), e
}

func writefile(file string, content []byte) (e error) {

	f, e := os.Create(file)
	if e != nil {
		return
	}
	defer f.Close()
	n, e := f.Write(content)
	if n != len(content) {
		return fmt.Errorf(`can't write file %s`, file)
	}
	return
}

// NewHelmContext
func NewHelmContext(name string, incluster bool,
	kubeconfig []byte, context, namespace string,
	ExpiryTime time.Time,
	entrys ...adapter.Entry,
) (hc *HelmContext, e error) {
	hc = new(HelmContext)
	if name == "" {
		hc.Name = utils.RandomString(16)
	} else {
		hc.Name = name
	}
	tmp := viper.GetString(`TMP`)
	helmContextDir := tmp + `/` + hc.Name

	hc.ContextDir = helmContextDir
	hc.KubeContext = context
	hc.KubeNamespace = namespace
	hc.Incluster = incluster
	hc.Expiry = ExpiryTime.Unix()

	e = os.MkdirAll(hc.ContextDir, 0777)
	if e != nil {
		return
	}
	fl := utils.NewFileLock(hc.ContextDir + LOCKFILE)
	e = fl.LockWithTimeout(time.Second * 30)
	if e != nil {
		utils.Error().Err(e).Send()
		return
	}
	defer fl.Unlock()

	e = hc.writeconfig(kubeconfig)
	if e != nil {
		return
	}
	for _, entry := range entrys {
		entry, e = hc.writeCert([]byte(entry.CertFile),
			[]byte(entry.KeyFile),
			[]byte(entry.CAFile),
			entry)
		if e != nil {
			return
		}
		if entry.URL == "" {
			e = fmt.Errorf(`%s didn't have url'`, entry.Name)
			return
		}

	}
	e = hc.UpdateRepo(entrys...)
	return
}

// LoadHelmContext
func LoadHelmContext(helmContextName string) (hc *HelmContext, e error) {

	tmp := viper.GetString(`TMP`)
	helmContextDir := tmp + `/` + helmContextName
	fl := utils.NewFileLock(helmContextDir + LOCKFILE)
	e = fl.LockWithTimeout(time.Second * 30)
	if e != nil {
		if strings.Contains(e.Error(), `no such file or directory`) {
			e = fmt.Errorf(`no such context %v`, helmContextName)
		}
		return hc, e
	}
	defer fl.Unlock()

	config := viper.New()
	configPath := helmContextDir + `/` + CONTEXT_CONFIG
	f, e := os.Open(configPath)
	if e != nil {
		return hc, fmt.Errorf(`no such context %s`, helmContextName)
	}
	config.SetConfigType("yaml")
	e = config.ReadConfig(f)
	if e != nil {
		return hc, e
	}

	hc = &HelmContext{
		ContextDir:     config.GetString(`helmcontext.contextdir`),
		Name:           config.GetString(`helmcontext.name`),
		KubeConfigPath: config.GetString(`helmcontext.kubeconfigpath`),
		KubeContext:    config.GetString(`helmcontext.kubecontext`),
		KubeNamespace:  config.GetString(`helmcontext.kubenamespace`),
		Expiry:         config.GetInt64(`helmcontext.expiry`),
		Incluster:      config.GetBool(`helmcontext.incluster`),
	}
	return
}

// DelHelmContext
func DelHelmContext(helmContextName string) error {
	tmp := viper.GetString(`TMP`)
	helmContextDir := tmp + `/` + helmContextName
	fl := utils.NewFileLock(helmContextDir + LOCKFILE)
	e := fl.LockWithTimeout(time.Second * 30)
	if e != nil {
		if strings.Contains(e.Error(), `no such file or directory`) {
			e = fmt.Errorf(`no such context %v`, helmContextName)
		}
		return e
	}
	defer fl.Unlock()
	return os.RemoveAll(helmContextDir)
}

func (hc *HelmContext) writeCert(certFile, keyFile, caFile []byte, entry adapter.Entry) (res adapter.Entry, e error) {

	if len(certFile) != 0 {
		certPath := hc.ContextDir + "/" + entry.Name + "/repo.crt"
		e = writefile(certPath, certFile)
		if e != nil {
			utils.Error().Err(e).Send()
			return
		}
		entry.CertFile = certPath
	}
	if len(keyFile) != 0 {
		keyPath := hc.ContextDir + "/" + entry.Name + "/repo.key"
		e = writefile(keyPath, keyFile)
		if e != nil {
			utils.Error().Err(e).Send()
			return
		}
		entry.KeyFile = keyPath
	}
	if len(caFile) != 0 {
		caPath := hc.ContextDir + "/" + entry.Name + "/repo.ca"
		e = writefile(caPath, caFile)
		if e != nil {
			utils.Error().Err(e).Send()
			return
		}
		entry.CAFile = caPath
	}

	return entry, nil
}

func (hc *HelmContext) writeconfig(kubeconfig []byte) error {
	kubeConfigPath := hc.ContextDir + `/` + KUBECONFIG

	kubeconfigf, e := os.Create(kubeConfigPath)
	if e != nil {
		utils.Error().Err(e).Send()
		return e
	}
	defer kubeconfigf.Close()

	_, e = kubeconfigf.Write(kubeconfig)
	if e != nil {
		utils.Error().Err(e).Send()
		return e
	}

	config := viper.New()

	hc.KubeConfigPath = kubeConfigPath
	config.Set(`helmcontext`, hc)

	configdir := hc.ContextDir + `/` + CONTEXT_CONFIG
	return config.WriteConfigAs(configdir)
}

func listContext() (hcs []*HelmContext, e error) {
	tmp := viper.GetString(`TMP`)
	tmpdir, e := os.Open(tmp)
	if e != nil {
		return
	}
	fi, e := tmpdir.Readdir(0)
	if e != nil {
		return
	}
	for _, contextdir := range fi {
		if !contextdir.IsDir() {
			continue
		}
		hc, e := LoadHelmContext(contextdir.Name())
		if e != nil {
			continue
		}
		hcs = append(hcs, hc)
	}
	return
}

func ContextGC() {
	for {
		utils.ConnectionWaiter.Wait()
		hcs, e := listContext()
		if e != nil {
			utils.Warn().Err(e).Send()
			continue
		}
		for _, hc := range hcs {
			if hc.Expiry < time.Now().Unix() {
				e = DelHelmContext(hc.Name)
				if e != nil {
					utils.Warn().Err(e).Send()
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}
