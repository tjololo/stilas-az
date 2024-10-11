package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func StartResumeOperation[T any](ctx context.Context, poller *runtime.Poller[T]) (done bool, result T, resumeToken string, err error) {
	logger := log.FromContext(ctx)
	done = false
	res, err := poller.Poll(ctx)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = runtime.NewResponseError(res)
		return
	}
	if poller.Done() {
		done = true
		result, err = poller.Result(ctx)
	} else {
		resumeToken, err = poller.ResumeToken()
		if err != nil {
			logger.Error(err, "Failed to get resume Token")
		}
	}
	return
}
