package media

import (
	"fmt"
	"io"
	"sync"

	"github.com/pion/webrtc/v3"
)


type Track struct{
	ID  string
	Kind webrtc.RTPCodecType
	TrackLocal *webrtc.TrackLocalStaticRTP
	Mutex sync.Mutex
	IsClosed bool
}


func NewTrack(ID string,kind webrtc.RTPCodecType,codec webrtc.RTPCodecCapability)(*Track,error){
	var track *webrtc.TrackLocalStaticRTP
	var err error
	if kind ==webrtc.RTPCodecTypeAudio{
		track,err= webrtc.NewTrackLocalStaticRTP(codec,ID,"audio")
	}else if  kind == webrtc.RTPCodecTypeVideo{
		track,err= webrtc.NewTrackLocalStaticRTP(codec,ID,"video")
	} else {
		return nil, fmt.Errorf("invalid media type")
	}
	if err != nil {
        return nil, fmt.Errorf("failed to create track: %w", err)
    }
	return &Track{
		ID: ID,
		Kind: kind,
		TrackLocal: track,
		IsClosed: false,
	},nil
}

 func (t *Track) WriteToTrack(packet []byte)(error){
	t.Mutex.Lock()
    defer t.Mutex.Unlock()
	if t.IsClosed{
		return fmt.Errorf("track is closed")
	}
	if _,err:=t.TrackLocal.Write(packet);err!=nil{
        return fmt.Errorf("failed to write RTP packet: %w", err)
	}
	return nil
}

func (t *Track) Close()(error){
	t.Mutex.Lock()
    defer t.Mutex.Unlock()
	if t.IsClosed{
		return fmt.Errorf("track is closed")
	}
	t.IsClosed=true
	return nil
}


func RelayTrack(remoteTrack *webrtc.TrackRemote, localTrack *Track) {
	go func() {
		buf := make([]byte, 1500)
		for {
			n, _, err := remoteTrack.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Remote track closed")
				} else {
					fmt.Println("Error reading from remote track:", err)
				}
				localTrack.Close()
				break
			}

			if writeErr := localTrack.WriteToTrack(buf[:n]); writeErr != nil {
				fmt.Println("Error writing to local track:", writeErr)
				break
			}
		}
	}()
}
