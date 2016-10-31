package session

import (
	"net/http"

	"github.com/cookies"
	uuid "github.com/satori/go.uuid"
)

const (
	DefaultProjectId = "waterandboards-auth" // currently set like this because of I have to change the name in app engine
)

type Session interface {
	Set(w http.ResponseWriter, val map[string]string) error
	Get(req *http.Request) (map[string]string, error)
	Del(w http.ResponseWriter)
}

// The UserSession struct holds the pointer to the underlying database and a few other settings
type SessionType struct {
	cookie cookies.CookieMng // abstract a cookie secure cookie manager
	store  Store
}

// NewSession creates a default session with good default prod settings
func NewSession(projectId) (*SessionType, error) {
	if projectId == "" {
		projectId = DefaultProjectId
	}

	mng, err := NewUserManager(projectId)
	if err != nil {
		return nil, err
	}

	conf := &cookie.Conf{true, true, 0} // max age 0 means typical session when
	// the broswer is close the session ends
	return &Session{cookies.New("_session", conf), NewFileSystemStore()}, nil
}

// NewTestSession creates a new TestSession struct that can be used for managing users.
// with good default test settings
func NewTestSession(projectId) (*SessionType, error) {
	if projectId == "" {
		projectId = DefaultProjectId
	}

	mng, err := NewUserManager(projectId)
	if err != nil {
		return nil, err
	}

	conf := &cookie.Conf{false, false, 0} // max age 0 means typical session when
	// Cookies lasts for 24 hours by default. Specified in seconds.
	return &Session{cookies.New("_session", conf), NewMapStore()}, nil
}

// NewSessionWithConf creates a new Session with that configuration.
func NewSessionWithConf(projectId, ssname string, conf *cookieConf) (*SessionType, error) {
	if projectId == "" {
		projectId = DefaultProjectId
	}

	mng, err := NewUserManager(projectId)
	if err != nil {
		return nil, err
	}

	return &Session{cookies.New(ssname), NewFileSystemStore()}, nil
}

func (ss *Session) Set(w http.ResponseWriter, value map[string]string) error {
	uid := uuid.NewV4()
	value["uuid"] = uid

	ss.cookie.SetCookieVal(w, nil, value)

	return nil
}

func (ss *Session) Get(req http.Request) (map[string]string, error) {
	return ss.cookie.GetCookieVal(nil, req), nil

}

func (ss *Session) Del(w http.ResponseWriter) {
	ss.cookie.Del(w, nil)
}
