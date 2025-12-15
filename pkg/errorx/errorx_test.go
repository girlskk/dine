package errorx

import (
	"net/http"
	"testing"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/e"
)

func TestError(t *testing.T) {
	// 业务错误，使用默认错误提示
	err := Fail(e.BadRequest, nil)
	t.Logf("%#v\n", err)

	err = FailWithStatus(e.ThirdPartyError, http.StatusInternalServerError, nil)
	t.Logf("%#v\n", err)

	err = Failf(e.ThirdPartyError, "third party error: %s", "test")
	t.Logf("%#v\n", err)

	err = FailWithStatusf(e.ThirdPartyError, http.StatusInternalServerError, "third party error: %s", "test")
	t.Logf("%#v\n", err)
}
