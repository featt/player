package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	m "github.com/featt/player/pkg/models"
	"github.com/gin-gonic/gin"
)



func AddSong(c *gin.Context)  {
    file, err := c.FormFile("file")
    if err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    }

    filename := filepath.Join("tracks", file.Filename)
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.Status(http.StatusCreated)
}

var pause chan bool 

func Play(c *gin.Context) {
    var playlist *m.Track 
    track, err := getCurrentTrack(c, playlist)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    filename := filepath.Join("tracks", track.Name)
    c.Header("Current-Track", track.Name)

    f, err := os.Open(filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer f.Close()

    duration := time.Duration(track.Duration) * time.Second
    timer := time.NewTimer(duration)

    done := make(chan bool) 

    go func() {
        <-timer.C 
        log.Println("Finished playing track:", track.Name)
        done <- true
    }()

    for {
        select {
        case <-done:
            next := track.Next
            if next == nil {
                log.Println("Reached end of playlist")
                return
            }

            log.Println("Switching to next track:", next.Name)
            sendTrack(c, next)

            track = next
            filename = filepath.Join("tracks", track.Name)
            c.Header("Current-Track", track.Name)

            f, err = os.Open(filename)
            if err != nil {
                log.Println("Error opening track file:", err)
                return
            }
            defer f.Close()

            duration = time.Duration(track.Duration) * time.Second
            timer = time.NewTimer(duration)

            done = make(chan bool)

            go func() {
                <-timer.C 
                log.Println("Finished playing track:", track.Name)
                done <- true
            }()

        case <-c.Done():
            log.Println("Playback interrupted:", c.Err())
            return

        case <-pause:           
            log.Println("Pausing playback")
            timer.Stop()
            <-done 
        }

        if len(pause) > 0 {           
            log.Println("Waiting for playback to resume")
            <-pause
            log.Println("Resuming playback")
            timer = time.NewTimer(duration)
            done = make(chan bool)

            go func() {
                <-timer.C 
                log.Println("Finished playing track:", track.Name)
                done <- true
            }()
        }
    }
}

func Pause(c *gin.Context) {
    pause <- true
    c.Status(http.StatusOK)
}


func Next(c *gin.Context) {
    var playlist *m.Track 
    track, err := getCurrentTrack(c, playlist)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    track = track.Next
    if track == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no next track"})
        return
    }

    sendTrack(c, track)
}

func Prev(c *gin.Context) {
    var playlist *m.Track 
    track, err := getCurrentTrack(c, playlist)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    track = track.Prev
    if track == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no previous track"})
        return
    }

    sendTrack(c, track)
}



func getCurrentTrack(c *gin.Context, playlist *m.Track) (*m.Track, error) {
    currentTrack := c.GetHeader("Current-Track")
    if currentTrack == "" {
        return nil, fmt.Errorf("Current-Track header not set")
    }

    for track := playlist; track != nil; track = track.Next {
        if track.Name == currentTrack {
            return track, nil
        }
    }

    return nil, fmt.Errorf("current track not found in playlist")
}

func sendTrack(c *gin.Context, track *m.Track) {
    filename := filepath.Join("tracks", track.Name)
    c.Header("Current-Track", track.Name)

    go func() {
        c.File(filename)
    }()
}
