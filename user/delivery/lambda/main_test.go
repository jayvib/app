// +build unit

package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

var update = flag.Bool("update", false, "Update the .golden file")

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestHandler(t *testing.T) {
	//log.SetOutput(ioutil.Discard)
	usecaseMock := new(mocks.Usecase)
	handler := newHandler(usecaseMock)
	tests := []struct {
		name       string
		handler    HandlerFunc
		req        events.APIGatewayProxyRequest
		assertResp func(t *testing.T, resp events.APIGatewayProxyResponse)
		mock       func()
		wantErr    bool
	}{
		{
			name:    "known user",
			handler: handler,
			req: events.APIGatewayProxyRequest{
				Path: "/user/dummyid",
			},
			mock: func() {
				want := &model.User{
					ID:        "uniqueid99",
					Email:     "sanji.vinsmoke@onepiece.com",
					Firstname: "Sanji",
					Lastname:  "Vinsmoke",
					Username:  "sanji.vinsmoke",
				}
				usecaseMock.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
					Return(want, nil).Once()
			},
			assertResp: func(t *testing.T, resp events.APIGatewayProxyResponse) {
				want := &model.User{
					ID:        "uniqueid99",
					Email:     "sanji.vinsmoke@onepiece.com",
					Firstname: "Sanji",
					Lastname:  "Vinsmoke",
					Username:  "sanji.vinsmoke",
				}

				var got model.User
				err := json.Unmarshal([]byte(resp.Body), &got)
				require.NoError(t, err)
				assert.Equal(t, &got, want)
			},
		},
		{
			name:    "unknown user",
			handler: handler,
			req: events.APIGatewayProxyRequest{
				Path: "/user/unknownid",
			},
			mock: func() {
				usecaseMock.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, apperr.New(apperr.NoItemFound, "not found item", nil)).
					Once()
			},
			assertResp: func(t *testing.T, resp events.APIGatewayProxyResponse) {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := tt.handler(context.Background(), tt.req)
			require.NoError(t, err)
			tt.assertResp(t, resp)
		})
	}
}
