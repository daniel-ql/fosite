package openid

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ory/fosite"
	"github.com/ory/fosite/internal"
	"github.com/ory/fosite/token/jwt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var strat = &DefaultStrategy{
	RS256JWTStrategy: &jwt.RS256JWTStrategy{
		PrivateKey: internal.MustRSAKey(),
	},
}

var fooErr = errors.New("foo")

func TestGenerateIDToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	chgen := internal.NewMockOpenIDConnectTokenStrategy(ctrl)
	defer ctrl.Finish()

	ar := fosite.NewAccessRequest(nil)
	sess := &DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Subject: "peter",
		},
		Headers: &jwt.Headers{},
	}
	h := &IDTokenHandleHelper{IDTokenStrategy: chgen}

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail because generator failed",
			setup: func() {
				ar.Form.Set("nonce", "11111111111111111111111111111111111")
				ar.SetSession(sess)
				chgen.EXPECT().GenerateIDToken(nil, ar).Return("", fooErr)
			},
			expectErr: fooErr,
		},
		{
			description: "should pass",
			setup: func() {
				chgen.EXPECT().GenerateIDToken(nil, ar).AnyTimes().Return("asdf", nil)
			},
		},
	} {
		c.setup()
		token, err := h.generateIDToken(nil, ar)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		if err == nil {
			assert.NotEmpty(t, token, "(%d) %s", k, c.description)
		}
		t.Logf("Passed test case %d", k)
	}

}

func TestIssueExplicitToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	resp := internal.NewMockAccessResponder(ctrl)
	defer ctrl.Finish()

	ar := fosite.NewAuthorizeRequest()
	ar.Form = url.Values{"nonce": {"111111111111"}}
	ar.SetSession(&DefaultSession{Claims: &jwt.IDTokenClaims{
		Subject: "peter",
	}, Headers: &jwt.Headers{}})

	resp.EXPECT().SetExtra("id_token", gomock.Any())
	h := &IDTokenHandleHelper{IDTokenStrategy: strat}
	err := h.IssueExplicitIDToken(nil, ar, resp)
	assert.NoError(t, err)
}

func TestIssueImplicitToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	resp := internal.NewMockAuthorizeResponder(ctrl)
	defer ctrl.Finish()

	ar := fosite.NewAuthorizeRequest()
	ar.Form = url.Values{"nonce": {"111111111111"}}
	ar.SetSession(&DefaultSession{Claims: &jwt.IDTokenClaims{
		Subject: "peter",
	}, Headers: &jwt.Headers{}})

	resp.EXPECT().AddFragment("id_token", gomock.Any())
	h := &IDTokenHandleHelper{IDTokenStrategy: strat}
	err := h.IssueImplicitIDToken(nil, ar, resp)
	assert.NoError(t, err)
}
