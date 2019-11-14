package brick

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"errors"

	graphql "github.com/graph-gophers/graphql-go"
)

// fork from github.com/graph-gophers/graphql-go/relay

func formatSchema(schema string) string {
	r := strings.Replace(schema, "\n", " ", -1)
	r = strings.Replace(r, "\t", " ", -1)
	r = strings.Trim(r, " ")
	r = regexp.MustCompile(`\s+`).ReplaceAllString(r, " ")
	return r
}

func relay(ctx context.Context, schema *graphql.Schema) {
	w := Response(ctx)
	r := Request(ctx)

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	access := Access(ctx)
	access["schema"] = formatSchema(params.Query)
	if os.Getenv("ENV") == "production" {
		access["schema_hash"] = fmt.Sprintf("%x", md5.Sum([]byte(access["schema"])))
	}
	if params.OperationName != "" {
		access["operation"] = params.OperationName
	}

	response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	is500 := false
	errorMsg := []string{}

	// https://github.com/graph-gophers/graphql-go/pull/207
	if response.Errors != nil {
		re := regexp.MustCompile(`{"code":\d+,"msg":".*?"}`)
		panicMsg := "graphql: panic occurred: "

		for _, rErr := range response.Errors {
			// extract business error
			if errMsg := re.FindString(rErr.Message); errMsg != "" {
				rErr.Message = errMsg
				continue
			}

			is500 = true

			// handle panic error
			if strings.Contains(rErr.Message, panicMsg) {
				// only log panic msg if it is not empty
				if rErr.Message != panicMsg {
					errMsg := rErr.Message
					errorMsg = append(errorMsg, errMsg)
				}

				// mask panic response
				rErr.Message = "masked panic"
				continue
			}

			errorMsg = append(errorMsg, rErr.Message)
		}
	}

	if len(errorMsg) > 0 {
		log.Error().Err(errors.New(strings.Join(errorMsg, ", "))).Send()
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	w.Header().Set("Content-Type", "application/json")
	if is500 {
		http.Error(w, string(responseJSON), http.StatusInternalServerError)
	} else {
		w.Write(responseJSON)
	}
}
