package auth

import (
	"errors"
	"net/http"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

// LambdaAuthenticator
type LambdaAuthenticator struct{}

func (l *LambdaAuthenticator) Authenticate(r *http.Request) (Claims, error) {
	ctx := r.Context()

	req, ok := core.GetAPIGatewayContextFromContext(ctx)
	if !ok {
		return Claims{}, errors.New("cound not get API Gateway context from request")
	}

	rawClaims, ok := req.Authorizer["claims"].(map[string]any)
	if !ok {
		return Claims{}, errors.New("could not retrieve authorizer claims")
	}

	sub, ok := rawClaims["sub"].(string)
	if !ok {
		return Claims{}, errors.New("could not parse sub field")
	}

	email, ok := rawClaims["email"].(string)
	if !ok {
		return Claims{}, errors.New("could not parse email field")
	}

	c := Claims{
		Sub:   sub,
		Email: email,
	}

	return c, nil
}
