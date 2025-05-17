package media

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)


type Manager struct{
	Sessions map[string]*Session
	Mutex sync.RWMutex

}

func NewManager()*Manager{
	manager := &Manager{
		Sessions: map[string]*Session{},
	}

	return manager
}

func (m * Manager) CreateSession(config webrtc.Configuration)(string,*Session,error){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
		log.Printf("Session created")
	newSession ,err:= NewSession(config)
	if err!=nil{
		return "",nil,err
	}
	sessionId := uuid.New().String()
	m.Sessions[sessionId]= newSession
	return sessionId,newSession,nil

}

func (m * Manager) CreateSessionWithTrack(config webrtc.Configuration,videoTrack *Track)(string,*Session,error){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
		log.Printf("Session created")
	newSession ,err:= NewSession(config)
	if err != nil {
		return "", nil, err
	}

	newSession.PeerConnection.AddTrack(videoTrack.TrackLocal)
	newSession.PeerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		go func ()  {
			buf:= make([]byte,1500)
			for{
				n,_,err :=tr.Read(buf)
				if err != nil {
                    log.Println("Remote track read error:", err)
                    return
                }
				err= videoTrack.WriteToTrack(buf[:n])
				if err != nil {
                    log.Println("Failed to write RTP:", err)
                    return
                }
			}
		}()
	})
	newSession.LocalTracks = append(newSession.LocalTracks, videoTrack.TrackLocal)
	sessionId := uuid.New().String()
	m.Sessions[sessionId]= newSession
	return sessionId,newSession,nil

}



func (m * Manager) JoinSession(sessionID string)(*Session,error){
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	
	session,exist := m.Sessions[sessionID]
	if !exist{
		return nil,fmt.Errorf("to join session doesnt exist %v",sessionID)
	}
	return session,nil
	

}

func (m * Manager) RemoveSession(sessionID string)(error){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	session,exist := m.Sessions[sessionID]
	if !exist{
		return fmt.Errorf("for removing session doesnt exist")
	}
	if err:=session.Close();err!=nil{
		return err
	}
	return nil

}

func (m * Manager) BroadcastICE(sessionID string,candidate webrtc.ICECandidateInit)(error){
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	session,exist := m.Sessions[sessionID]
	if !exist{
		return fmt.Errorf("broadcast session doesnt exist")
	}
	return session.AddICECandidate(candidate)
}