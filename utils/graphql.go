package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	req "github.com/imroc/req"
	be "github.com/pickjunk/brick/error"
)

// Graphql struct
type Graphql struct {
	URL       string
	Query     string
	Variables map[string]interface{}
	Operation string
	Headers   req.Header
}

// Fetch execute a graphql api
func (g *Graphql) Fetch(ctx context.Context, result interface{}) error {
	if os.Getenv("DEBUG") == "true" {
		req.Debug = true
		defer func() {
			req.Debug = false
		}()
	}

	r := req.New()

	params := map[string]interface{}{
		"query": g.Query,
	}
	if g.Variables != nil {
		params["variables"] = g.Variables
	}
	if g.Variables != nil {
		params["operation"] = g.Operation
	}
	res, err := r.Post(g.URL, g.Headers, req.BodyJSON(params), ctx)
	if err != nil {
		return err
	}

	code := res.Response().StatusCode
	if !(code >= 200 && code < 300) {
		return fmt.Errorf("http status error: %d", code)
	}

	var e struct {
		Errors []struct {
			Message string
		}
	}
	err = res.ToJSON(&e)
	if err != nil {
		return err
	}
	if len(e.Errors) > 0 {
		var bErr be.BusinessError
		err = json.Unmarshal([]byte(e.Errors[0].Message), &bErr)
		if err != nil {
			return errors.New(e.Errors[0].Message)
		}
		return &bErr
	}

	err = res.ToJSON(result)
	if err != nil {
		return err
	}

	return nil
}
