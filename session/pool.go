package session

import (
	"fmt"
	"sync"
)

func NewSessionPool() *SessionPool {
	return &SessionPool{
		sessions: map[string]ISession{},
		mutex:    sync.Mutex{},
	}
}

type SessionPool struct {
	sessions map[string]ISession
	mutex    sync.Mutex
}

func (o *SessionPool) Get(sessID string) (ISession, error) {
	sess, exist := o.sessions[sessID]
	if !exist {
		return nil, fmt.Errorf("try to get sessID sessionPool. "+
			"SessID not found. ID: %v", sessID)
	}

	if sess.IsClose() {
		o.RemoveAndClose(sessID)

		return nil, fmt.Errorf("get closed session from sessionPool. Remove session. "+
			"ID: %v, Type: %v", sess.GetId(), sess.GetType())
	}

	return sess, nil
}

func (o *SessionPool) Put(sess ISession) error {
	sessType := sess.GetType()
	sessID := sess.GetId()

	if sess.IsClose() {
		return fmt.Errorf("try to put in sessionPool closed session. "+
			"ID: %v, Type: %v", sessID, sessType)
	}

	_, exist := o.sessions[sessID]
	if exist {
		return fmt.Errorf("session with same ID already in sessionPool. "+
			"ID: %v, Type: %v", sessID, sessType)
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.sessions[sessID] = sess

	return nil
}

func (o *SessionPool) RemoveAndClose(sessID string) error {
	sess, exist := o.sessions[sessID]
	if !exist {
		return fmt.Errorf("try to get sessID from sessionPool. "+
			"SessID not found. ID: %v", sessID)
	}

	if !sess.IsClose() {
		sess.Close()
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	delete(o.sessions, sessID)

	return nil
}
