package middle

import (
	"github.com/weapons97/helm-api/pkg/helm"
	"github.com/weapons97/helm-api/pkg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func UnaryConnectionWaiter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	utils.ConnectionWaiter.Add(1)
	defer utils.ConnectionWaiter.Done()
	return handler(ctx, req)
}

func UnaryHelmClient(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	opt := helm.WithHelmContextNamespaceOpts(req)
	client, e := helm.NewHelmClient(req, opt)
	if e == helm.NoContextInfo {
		return handler(context.TODO(), req)
	}
	if e != nil {
		utils.Error().Err(e).Msg(`can't create helm context`)
		return nil, fmt.Errorf(`can't create helm contxt`)
	}
	withclient := context.WithValue(ctx, helm.HelmClientInContextName, client)
	return handler(withclient, req)
}
