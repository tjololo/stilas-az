package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

func StartResumeOperation[T any](ctx context.Context, poller *runtime.Poller[T]) (done bool, result T, resumeToken string, err error) {
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
	}
	return
}
