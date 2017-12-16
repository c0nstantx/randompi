package controllers

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"randompi/app/services"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

const mediaPath = "/media/videos"
const omxplayerPath = "/home/pi/omxplayer"

func (c App) Index() revel.Result {
	vl := services.VideoList(mediaPath)

	return c.Render(vl)
}

func (c App) Random() revel.Result {
	stopPlayer()
	vl := services.VideoList(mediaPath)
	fileList := ""
	for _, v := range vl {
		cmd := exec.Command(getOmxplayer(), "-p", "-o", "hdmi", v.Path)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(fileList)
	return c.Render()
}

func (c App) Play() revel.Result {
	stopPlayer()
	vl := services.VideoList(mediaPath)
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	cmd := exec.Command(getOmxplayer(), "-p", "-o", "hdmi", video.Path)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return c.Render(video)
}

func (c App) Stop() revel.Result {
	stopPlayer()
	return c.Redirect(App.Index)
}

func (c App) Pause() revel.Result {
	pausePlayer()
	vl := services.VideoList(mediaPath)
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	return c.Render(video)
}

func (c App) Resume() revel.Result {
	pausePlayer()
	vl := services.VideoList(mediaPath)
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	return c.Render(video)
}

func stopPlayer() {
	cmd := exec.Command(getDbusControl(), "stop")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func pausePlayer() {
	cmd := exec.Command(getDbusControl(), "pause")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getOmxplayer() string {
	return filepath.Join(omxplayerPath, "omxplayer")
}

func getDbusControl() string {
	return filepath.Join(omxplayerPath, "dbuscontrol.sh")
}
