package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"

	"github.com/mdebrouwer/exchange/auth"
	"github.com/mdebrouwer/exchange/auth/token"
)

type memStore struct {
	tokens map[string]token.User
}

func (s *memStore) Register(token token.Token, user token.User) (bool, error) {
	str := string(token)
	_, ok := s.tokens[str]
	if !ok {
		s.tokens[str] = user
	}
	return !ok, nil
}

func (s *memStore) Get(token token.Token) (user token.User, ok bool, err error) {
	user, ok = s.tokens[string(token)]
	return
}

func cookies(w *httptest.ResponseRecorder) []*http.Cookie {
	r := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	return r.Cookies()
}

var _ = Describe("token", func() {
	var provider auth.Provider
	var store *memStore

	BeforeEach(func() {
		store = &memStore{tokens: make(map[string]token.User)}
		provider = token.NewProvider(
			securecookie.GenerateRandomKey(32),
			securecookie.GenerateRandomKey(32),
			[]string{"some_token", "some_other_token"},
			store,
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

	Describe("when registering a user", func() {
		var userRequest *http.Request
		var userResponse *httptest.ResponseRecorder

		BeforeEach(func() {
			userRequest = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name":"Bruce Wayne","email":"batman@example.com"}`))
		})

		JustBeforeEach(func() {
			userResponse = httptest.NewRecorder()
			provider.UserHandler()(userResponse, userRequest)
		})

		Describe("when the user has a session", func() {
			BeforeEach(func() {
				request := httptest.NewRequest("POST", "/sessions", strings.NewReader("some_token"))
				response := httptest.NewRecorder()
				provider.SessionHandler()(response, request)

				for _, c := range cookies(response) {
					userRequest.AddCookie(c)
				}
			})

			It("should be successful", func() {
				Expect(userResponse.Code).To(Equal(http.StatusCreated))
			})

			It("should save the user", func() {
				user := store.tokens["some_token"]
				Expect(user.Name).To(Equal("Bruce Wayne"))
				Expect(user.Email).To(Equal("batman@example.com"))
			})

			Describe("when the token has already been registered", func() {
				BeforeEach(func() {
					store.tokens["some_token"] = token.User{Name: "Alfred Pennyworth", Email: "best_butler@example.com"}
				})

				It("should not be successful", func() {
					Expect(userResponse.Code).To(Equal(http.StatusConflict))
				})
			})
		})

		Describe("when there is no session", func() {
			It("should not be successful", func() {
				Expect(userResponse.Code).To(Equal(http.StatusUnauthorized))
			})

			It("should save not the user", func() {
				Expect(store.tokens).To(BeEmpty())
			})
		})
	})

	Describe("when getting the current user details", func() {
		var userRequest *http.Request
		var userResponse *httptest.ResponseRecorder

		BeforeEach(func() {
			userRequest = httptest.NewRequest("GET", "/user", strings.NewReader(""))
		})

		JustBeforeEach(func() {
			userResponse = httptest.NewRecorder()
			provider.UserHandler()(userResponse, userRequest)
		})

		Describe("when the user has a session", func() {
			BeforeEach(func() {
				request := httptest.NewRequest("POST", "/sessions", strings.NewReader("some_token"))
				response := httptest.NewRecorder()
				provider.SessionHandler()(response, request)

				for _, c := range cookies(response) {
					userRequest.AddCookie(c)
				}
			})

			Describe("when the user has been registered", func() {
				BeforeEach(func() {
					store.tokens["some_token"] = token.User{Name: "Alfred Pennyworth", Email: "best_butler@example.com"}
				})

				It("should be successful", func() {
					Expect(userResponse.Code).To(Equal(http.StatusOK))
				})

				It("should return the user", func() {
					Expect(userResponse.Body.String()).To(MatchJSON(`{"name":"Alfred Pennyworth","email":"best_butler@example.com"}`))
				})
			})

			Describe("when the user has not been registered", func() {
				It("should not be successful", func() {
					Expect(userResponse.Code).To(Equal(http.StatusNotFound))
				})
			})
		})

		Describe("when there is no session", func() {
			It("should not be successful", func() {
				Expect(userResponse.Code).To(Equal(http.StatusUnauthorized))
			})

			It("should save not the user", func() {
				Expect(store.tokens).To(BeEmpty())
			})
		})
	})

	Describe("when getting accessing a filtered resource", func() {
		var resourceRequest *http.Request
		var resourceResponse *httptest.ResponseRecorder
		var resourceAccessed bool
		var accessingUser token.User

		BeforeEach(func() {
			accessingUser = token.User{}
			resourceAccessed = false
			resourceRequest = httptest.NewRequest("GET", "/filtered_resource", strings.NewReader(""))
		})

		JustBeforeEach(func() {
			resourceResponse = httptest.NewRecorder()
			provider.Filter(func(w http.ResponseWriter, r *http.Request) {
				accessingUser = context.Get(r, "user").(token.User)
				resourceAccessed = true
			})(resourceResponse, resourceRequest)
		})

		Describe("when the user has a session", func() {
			BeforeEach(func() {
				request := httptest.NewRequest("POST", "/sessions", strings.NewReader("some_token"))
				response := httptest.NewRecorder()
				provider.SessionHandler()(response, request)

				for _, c := range cookies(response) {
					resourceRequest.AddCookie(c)
				}
			})

			Describe("when the user has been registered", func() {
				BeforeEach(func() {
					store.tokens["some_token"] = token.User{Name: "Alfred Pennyworth", Email: "best_butler@example.com"}
				})

				It("should call the downstream resource", func() {
					Expect(resourceAccessed).To(BeTrue())
				})

				It("should pass down the user in the context", func() {
					Expect(accessingUser.Name).To(Equal("Alfred Pennyworth"))
				})
			})

			Describe("when the user has not been registered", func() {
				It("should not be successful", func() {
					Expect(resourceResponse.Code).To(Equal(http.StatusForbidden))
				})

				It("should not call the downstream resource", func() {
					Expect(resourceAccessed).To(BeFalse())
				})
			})
		})

		Describe("when there is no session", func() {
			It("should not be successful", func() {
				Expect(resourceResponse.Code).To(Equal(http.StatusUnauthorized))
			})

			It("should not call the downstream resource", func() {
				Expect(resourceAccessed).To(BeFalse())
			})
		})
	})
})
