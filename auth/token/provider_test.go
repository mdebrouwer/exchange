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
