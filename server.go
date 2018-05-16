package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

type Group struct {
	Title string
	Dates []string
}

func readFiles() []Group {
	r, _ := regexp.Compile("data/links-(.*)-(.*?).json")
	files, ef := filepath.Glob("data/links-*.json")
	if ef != nil {
		log.Fatal(ef)
	}

	// reverse order
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	lq := make(map[string]Group)
	for _, file := range files {
		match := r.FindStringSubmatch(file)
		//fmt.Println(match[1])
		//fmt.Println(match[2])
		d := match[1]
		q := match[2]

		qo := lq[q]
		if qo.Title == "" {
			qo = Group{Title: q, Dates: []string{}}
		}
		qo.Dates = append(qo.Dates, d)
		lq[q] = qo
	}

	var groups []Group
	for _, q := range lq {
		groups = append(groups, q)
	}
	return groups
}

type Config struct {
	Port          string
	DataDir       string
	NodesFile     string
	Title         string
	Legend        template.HTML
	ConvertUTC    bool
	RadialMinView int
}

func main() {

	// Load config from file
	confPtr := flag.String("conf", "./confsample.ini", "INI config file")
	flag.Parse()

	// Config struct with default values
	config := Config{Port: ":8080",
		DataDir:       "./data",
		NodesFile:     "nodes.json",
		Title:         "Some data",
		ConvertUTC:    true,
		RadialMinView: 5}
	err := ini.MapTo(&config, *confPtr)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	title := config.Title
	legend := config.Legend
	nodes := config.NodesFile
	utc := config.ConvertUTC

	// Create router
	r := gin.Default()

	r.Use(static.Serve("/data", static.LocalFile(config.DataDir, false)))
	r.Use(static.Serve("/js", static.LocalFile("./views/js", false)))
	r.HTMLRender = gintemplate.Default()

	r.GET("/", func(c *gin.Context) {
		groups := readFiles()
		//fmt.Printf("==========\n%+v\n==========\n", groups)
		c.HTML(200, "index.html", gin.H{
			"title":  title,
			"Groups": groups,
			/*"Groups": []Group{
				{Title: "test", Dates: []string{"21h02", "23h03"}},
			},*/
		})
	})
	r.GET("/g/:search/:time", func(c *gin.Context) {
		c.HTML(200, "grap.html", gin.H{
			"title":  title,
			"nodes":  nodes,
			"legend": legend,
			"utc":    utc,
			"search": c.Param("search"),
			"time":   c.Param("time"),
		})
	})
	r.GET("/t/:search/:time", func(c *gin.Context) {
		c.HTML(200, "time.html", gin.H{
			"title":  title,
			"nodes":  nodes,
			"legend": legend,
			"utc":    utc,
			"search": c.Param("search"),
			"time":   c.Param("time"),
		})
	})
	r.GET("/r/:search/:time", func(c *gin.Context) {
		c.HTML(200, "radia.html", gin.H{
			"title":   title,
			"nodes":   nodes,
			"legend":  legend,
			"utc":     utc,
			"minview": config.RadialMinView,
			"search":  c.Param("search"),
			"time":    c.Param("time"),
		})
	})

	// Listen and Server
	r.Run(config.Port)
}
