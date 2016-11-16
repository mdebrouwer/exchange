package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/mdebrouwer/exchange/auth"
	"github.com/mdebrouwer/exchange/auth/token"
)

type memStore struct {
	tokens map[string]token.User
}

func (s *memStore) Register(token string, user token.User) (bool, error) {
	_, ok := s.tokens[token]
	if !ok {
		s.tokens[token] = user
	}
	return !ok, nil
}

func cookies(w *httptest.ResponseRecorder) []*http.Cookie {
	r := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
  return r.Cookies()
}

var _ = Describe("token", func() {
	var provider auth.Provider

	BeforeEach(func() {
		provider = token.NewProvider(
			[]string{"some_token", "some_other_token"},
			&memStore{tokens: make(map[string]token.User)},
		)
	})

	Describe("when creating a session", func() {
		var reqToken string
		var response *httptest.ResponseRecorder

		JustBeforeEach(func() {
			request := httptest.NewRequest("POST", "/sessions", strings.NewReader(reqToken))
			response = httptest.NewRecorder()
			provider.SessionHandler()(response, request)
		})

		Describe("when using a valid token", func() {
			BeforeEach(func() {
				reqToken = "some_token"
			})

			It("should be successful", func() {
				Expect(response.Code).To(Equal(http.StatusCreated))
			})

			XDescribe("when performing a subsequent request", func() {
				var subsequentReqToken string
				var subsequentResponse *httptest.ResponseRecorder
				JustBeforeEach(func() {
					request := httptest.NewRequest("POST", "/sessions", strings.NewReader(subsequentReqToken))
					for _, c := range cookies(response) {
						request.AddCookie(c)
					}
					subsequentResponse = httptest.NewRecorder()
					provider.SessionHandler()(response, request)
				})

				BeforeEach(func() {
					subsequentReqToken = reqToken
				})

				It("should be successful", func() {
					Expect(response.Code).To(Equal(http.StatusOK))
				})

				Describe("when the subsequent token is invalid", func() {
					BeforeEach(func() {
						subsequentReqToken = "some_invalid_token"
					})

					It("should not be successful", func() {
						Expect(subsequentResponse.Code).To(Equal(http.StatusUnauthorized))
					})
				})
			})
		})

		Describe("when using an invalid token", func() {
			BeforeEach(func() {
				reqToken = "some_invalid_token"
			})

			It("should not be successful", func() {
				Expect(response.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
