package controllers

import (
	"log"
	"os/exec"
	"randompi/app/services"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	vl := services.VideoList("/home/kostasx/Downloads/test_videos")

	return c.Render(vl)
}

func (c App) Random() revel.Result {
	closePlayers()
	vl := services.VideoList("/home/kostasx/Downloads/test_videos")
	// command := "vlc"
	for _, v := range vl {
		cmd := exec.Command("vlc", "--playlist-enqueue", v.Path)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
	return c.Render()
}

func (c App) Play() revel.Result {
	closePlayers()
	vl := services.VideoList("/home/kostasx/Downloads/test_videos")
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	cmd := exec.Command("vlc", "file://"+video.Path)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return c.Render(video)
}

func (c App) Stop() revel.Result {
	closePlayers()
	return c.Redirect(App.Index)
}

func (c App) Pause() revel.Result {
	cmd := exec.Command("vlc", "vlc://pause")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	vl := services.VideoList("/home/kostasx/Downloads/test_videos")
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	return c.Render(video)
}

func (c App) Resume() revel.Result {
	cmd := exec.Command("vlc", "vlc://pause")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	vl := services.VideoList("/home/kostasx/Downloads/test_videos")
	videoHash := c.Params.Query.Get("v")
	video := vl[videoHash]
	return c.Render(video)
}

func closePlayers() {
	// cmd := exec.Command("vlc", "vlc://stop")
	cmd := exec.Command("killall", "-9", "vlc")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
