package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ab22/abcd/router/httputils"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

type MiddlewareFunc func(httputils.APIHandler) httputils.APIHandler

// Go's http.FileServer by default, lists the directories and files
// of the specified folder to serve and cannot be disabled.
// To prevent directory listing, noDirListing checks if the
// path requests ends in '/'. If it does, then the client is requesting
// to explore a folder and we return a 404 (Not found), else, we just
// call the http.Handler passed as parameter.
func NoDirListing(h httputils.APIHandler) httputils.APIHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		urlPath := r.URL.Path

		if strings.HasSuffix(urlPath, "/") {
			http.NotFound(w, r)
			return nil
		}

		return h(ctx, w, r)
	}
}

// Validates that the user cookie is set up before calling the handler
// passed as parameter.
func ValidateAuth(h httputils.APIHandler) httputils.APIHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var (
			sessionData *SessionData
			err         error
			ok          bool
			cookieStore *sessions.CookieStore
			session     *sessions.Session
		)

		cookieStore, ok = ctx.Value("cookieStore").(*sessions.CookieStore)

		if !ok {
			return fmt.Errorf("validate auth: could not cast value as cookie store:", ctx.Value("cookieStore"))
		}

		session, err = cookieStore.Get(r, SessionCookieName)

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		}

		if sessionData, ok = session.Values["data"].(*SessionData); !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		}

		authenticatedContext := context.WithValue(ctx, "sessionData", sessionData)
		return h(authenticatedContext, w, r)
	}
}

// JsonHandler is a middleware function for the ApiHandler type. The
// JsonHandler sets the content type of the response as application/json. It
// also checks if the ApiHandler returned an error. When the ApiHandler returns
// an error, then a http.Error is filled with the error data return from the
// ApiHandler. If there's no error, the payload that returned from the
// ApiHandler is parsed to json (if any) and written to the
// http.ResponseWriter.

/*func JsonHandler(h ApiHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload, apiError := h(w, r)
		if apiError != nil {
			if apiError.Error != nil {
				log.Println(apiError.Error)
			}

			// If no message is set, get the default error message from the
			// http module.
			if apiError.Message == "" {
				apiError.Message = http.StatusText(apiError.HttpCode)
			}

			http.Error(w, apiError.Message, apiError.HttpCode)
			return
		}

		// If nothing is to be converted to json then return
		if payload == nil {
			return
		}

		// Convert payload to json and return HTTP OK (200)
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})
}*/

// gzipContent is a middleware function for handlers to encode content to gzip.
func GzipContent(h httputils.APIHandler) httputils.APIHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Add("Vary", "Accept-Encoding")

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return h(ctx, w, r)
		}

		w.Header().Set("Content-Encoding", "gzip")

		gzipResponseWriter := httputils.NewGzipResponseWriter(w)
		defer gzipResponseWriter.Close()

		return h(ctx, gzipResponseWriter, r)
	}
}

// Authorize validates privileges for the current user. Each route must have
// an array of privileges that point which users can make a call to it.
//
// Note:
//
// It is assumed that ValidateAuth was called before this function, or at
// least some other session check was done before this.

//func Authorize(requiredRoles []string, h http.Handler) http.HandlerFunc {
func Authorize(h httputils.APIHandler) httputils.APIHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var (
			requiredRoles []string
			sessionData   *SessionData
			route         Route
			ok            bool
		)

		sessionData, ok = ctx.Value("sessionData").(*SessionData)

		if !ok {
			return fmt.Errorf("authorize: could not cast value as session data:", ctx.Value("sessionData"))
		}

		route, ok = ctx.Value("route").(Route)

		if !ok {
			return fmt.Errorf("authorize: could not cast value as route:", ctx.Value("route"))
		}

		requiredRoles = route.RequiredRoles()

		if len(requiredRoles) == 0 {
			return h(ctx, w, r)
		}

		for _, role := range requiredRoles {
			if role == "ADMIN" && sessionData.IsAdmin {
				return h(ctx, w, r)
			} else if role == "TEACHER" && sessionData.IsTeacher {
				return h(ctx, w, r)
			}
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}
}