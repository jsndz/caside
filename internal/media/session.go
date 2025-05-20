package media

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
)

type Session struct{
	PeerConnection *webrtc.PeerConnection
	LocalTracks []*webrtc.TrackLocalStaticRTP
	RemoteTracks []*webrtc.TrackRemote
	ICECandidate []webrtc.ICECandidateInit
	Mutex sync.Mutex
}


func  NewSession(config webrtc.Configuration) (*Session,error){
	peerConc, err := webrtc.NewPeerConnection(config)
	peerConc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Println("ICE Connection State has changed to:", state.String())
	})
	
	peerConc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Println("Peer Connection State has changed to:", state.String())
	})
	
	
	
	if err!=nil{
		return nil, err
	}
	session := &Session{
		PeerConnection: peerConc,
		LocalTracks: []*webrtc.TrackLocalStaticRTP{},
		RemoteTracks :[]*webrtc.TrackRemote{},
		ICECandidate :[]webrtc.ICECandidateInit{},
	}
	return session,nil
}


func (s *Session) AddICECandidate(candidate webrtc.ICECandidateInit)(error){
	s.Mutex.Lock()
    defer s.Mutex.Unlock()

	if err:= s.PeerConnection.AddICECandidate(candidate);err!= nil {
		return err
	}
	s.ICECandidate = append(s.ICECandidate, candidate)
	return nil
}

func (s *Session) HandleSDP(sdp webrtc.SessionDescription)(*webrtc.SessionDescription,error){
	s.Mutex.Lock()
    defer s.Mutex.Unlock()

	if sdp.Type == webrtc.SDPTypeOffer{
		if err:= s.PeerConnection.SetRemoteDescription(sdp);err!=nil{
			return nil,err
		}
		answer,err := s.PeerConnection.CreateAnswer(nil)
		if err!=nil{
			return nil,err
		}
		if err:= s.PeerConnection.SetLocalDescription(answer);err!=nil{
			return nil,err
		} 
		return &answer,nil
	} else if sdp.Type==webrtc.SDPTypeAnswer{
		if err:= s.PeerConnection.SetRemoteDescription(sdp);err!=nil{
			return nil,err
		}
	}
	return nil,nil
}

func (s *Session) Close() error {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()

	s.LocalTracks = nil

    if err := s.PeerConnection.Close(); err != nil {
        return err
    }

    return nil
}
