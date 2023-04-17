package kornet

import (
  "net/http"
  "net/url"
)

// fallback to use default values

func KRequestGetURL(req *http.Request) url.URL {

  URL := url.URL{
    User:     req.URL.User,
    Scheme:   "http",
    Host:     "127.0.0.1",
    Path:     "/",
    RawQuery: "",
  }

  if req.URL.Scheme != "" {

    URL.Scheme = req.URL.Scheme
  }

  if req.URL.Host != "" {

    URL.Host = req.URL.Host
  }

  if req.URL.Path != "" {

    URL.Path = req.URL.Path
  }

  if req.URL.RawQuery != "" {

    URL.RawQuery = req.URL.RawQuery
  }

  return URL
}
