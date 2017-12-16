package controllers

import (
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"randompi/app/services"
	"syscall"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

const mediaPath = "/media/videos"
const omxplayerPath = "/home/pi/omxplayer"
const rootPassword = "pipass123"

func (c App) Index() revel.Result {
	vl := services.VideoList(mediaPath)

	return c.Render(vl)
}

func (c App) Reboot() revel.Result {
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	// cmd := exec.Command("echo", rootPassword, "|", "sudo", "-S", "reboot")
	// err := cmd.Run()
	// if err != nil {
	// 	log.Print(err)
	// }

	return c.Render()
}

func (c App) Random() revel.Result {
	stopPlayer()
	vl := services.VideoList(mediaPath)
	i := len(vl)
	rand := rand.Intn(i) + 1
	for _, v := range vl {
		if rand == i {
			cmd := exec.Command(getOmxplayer(), "-p", "-o", "hdmi", v.Path)
			err := cmd.Start()
			if err != nil {
				log.Print(err)
			}
			break
		}
		i--
	}

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
		log.Print(err)
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
	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

func pausePlayer() {
	cmd := exec.Command(getDbusControl(), "pause")
	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

func getOmxplayer() string {
	return filepath.Join(omxplayerPath, "omxplayer")
}

func getDbusControl() string {
	return filepath.Join(omxplayerPath, "dbuscontrol.sh")
}
