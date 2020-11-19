package adapter

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/repo"
	"log"
	"strings"
	"sync"
	"testing"
)

var (
	tentry = repo.Entry{
		Name:     "testRepo",
		URL:      "https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts",
		Username: "",
		Password: "",
		CertFile: "",
		KeyFile:  "",
		CAFile:   "",
	}
)

func testConcurrency(n int, tf func()) {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Printf("%T, %v \n", r, r)
				}
			}()
			tf()
		}()
	}
	wg.Wait()
}

func TestMergeValues(t *testing.T) {
	v1 := `
aaa: vaaa1
aaa2:
  bbb:
    ccc: cvvv4
aaa3:
  bbb: bvvv
aaa4:
  ccc: 1
`
	v2 := `
aaa:
  bbb: bvvv1
  bbbn: bvvv4
aaa2:
  bbb3:
    ccc: cvvv4
aaa3:
  bbb4: bvvv3
`
	res := `
aaa:
  bbb: bvvv1
  bbbn: bvvv4
aaa2:
  bbb:
    ccc: cvvv4
  bbb3:
    ccc: cvvv4
aaa3:
  bbb: bvvv
  bbb4: bvvv3
aaa4:
  ccc: 1`
	v3, e := MergeValues([]byte(v1), []byte(v2))
	require.NoError(t, e)
	require.Equal(t, bytes.TrimSpace(v3), strings.TrimSpace(res))
}
