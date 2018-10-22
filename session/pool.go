package session

import (
	"fmt"
	"sync"
)

func NewSessionPool() *SessionPool {
	return &SessionPool{
		sessions: map[SessionType]map[string]ISession{},
		mutex:    sync.Mutex{},
	}
}

type SessionPool struct {
	sessions map[SessionType]map[string]ISession
	mutex    sync.Mutex
}

func (o *SessionPool) Get(sessType SessionType, sessID string) (ISession, error) {
	typeToSessions, exist := o.sessions[sessType]
	if !exist {
		return nil, fmt.Errorf("try to get session from sessionPool. "+
			"SessionType not found. ID: %s, Type: %s", sessID, sessType)
	}

	sess, exist := typeToSessions[sessID]
	if !exist {
		return nil, fmt.Errorf("try to get sessID sessionPool. "+
			"SessID not found. ID: %v, Type: %v", sessID, sessType)
	}

	if sess.IsClose() {
		o.RemoveAndClose(sessType, sessID)

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

	_, exist := o.sessions[sessType]
	if !exist {
		o.sessions[sessType] = map[string]ISession{}
	}

	_, exist = o.sessions[sessType][sessID]
	if exist {
		return fmt.Errorf("session with same Type and ID already in sessionPool. "+
			"ID: %v, Type: %v", sessID, sessType)
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.sessions[sessType][sessID] = sess

	return nil
}

func (o *SessionPool) RemoveAndClose(sessType SessionType, sessID string) error {
	typeToSessions, exist := o.sessions[sessType]
	if !exist {
		return fmt.Errorf("try to get session from sessionPool. "+
			"SessionType not found. ID: %v, Type: %v", sessID, sessType)
	}

	sess, exist := typeToSessions[sessID]
	if !exist {
		return fmt.Errorf("try to get sessID from sessionPool. "+
			"SessID not found. ID: %v, Type: %v", sessID, sessType)
	}

	if !sess.IsClose() {
		sess.Close()
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	delete(o.sessions[sessType], sessID)

	return nil
}
