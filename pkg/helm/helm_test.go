package helm

import (
	"github.com/weapons97/helm-api/pkg/helm/adapter"
	"github.com/davecgh/go-spew/spew"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	_ "google.golang.org/grpc"
	"testing"
)

func init() {
	viper.Set(`TMP`, `testdata`)
	viper.Set(`DEBUG`, `true`)
}

func TestHelmContext(t *testing.T) {
	ctxName := `test`
	entry := adapter.Entry{
		Name:     "testRepo",
		URL:      "https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts",
		Username: "",
		Password: "",
		CertFile: "",
		KeyFile:  "",
		CAFile:   "",
	}
	resHelmContext := &HelmContext{
		ContextDir:     "testdata/test",
		Name:           "test",
		KubeConfigPath: "testdata/test/kubeconfig.yaml",
		KubeContext:    "",
		KubeNamespace:  "default",
	}
	hcn, e := NewHelmContext(ctxName, false,
		[]byte(`testconfig context`), "", "default",
		time.Now(),
		entry,
	)
	require.NoError(t, e)
	require.Equal(t, ctxName, hcn.Name)

	hc, e := LoadHelmContext(ctxName)
	require.NoError(t, e)
	require.Equal(t, resHelmContext, hc)

	e = DelHelmContext(ctxName)
	require.NoError(t, e)
}

func TestListContext(t *testing.T) {
	hcs, e := listContext()
	require.NoError(t, e)
	spew.Dump(hcs)
}
