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
		log.Printf("OnTrack triggered: SSRC=%d, ID=%s, PayloadType=%d", tr.SSRC(), tr.ID(), tr.PayloadType())
		localTrack, err := NewTrack(tr.ID(), tr.Kind(), tr.Codec().RTPCodecCapability)
		if err != nil {
			log.Println("Failed to create local track:", err)
			return
		}
	
		// Step 2: Add local track to the PeerConnection
		_, err = newSession.PeerConnection.AddTrack(localTrack.TrackLocal)
		if err != nil {
			log.Println("Failed to add local track to PeerConnection:", err)
			return
		}
	
		// Step 3: Start relaying media from remote to local track
		go RelayTrack(tr, localTrack)
	
		// Step 4: Optionally store the local track for later
		newSession.LocalTracks = append(newSession.LocalTracks, localTrack.TrackLocal)
	})
	
	newSession.LocalTracks = append(newSession.LocalTracks, videoTrack.TrackLocal)
	sessionId := "123X"
	m.Sessions[sessionId]= newSession
	return sessionId,newSession,nil

}



func (m * Manager) JoinSession(sessionID string)(*Session,error){
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	log.Printf("Joining session")

	session,exist := m.Sessions["123X"]
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