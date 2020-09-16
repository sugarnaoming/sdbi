package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sdbi/config"
	"sdbi/display"
	"sdbi/kubernetes"
	"sdbi/screwdriver"
	"sync"

	"github.com/spf13/cobra"
)

type buildInfo struct {
	build    *screwdriver.Build
	job      *screwdriver.Job
	pipeline *screwdriver.Pipeline
	buildID  string
}

var Get = &cobra.Command{
	Use:   "get",
	Short: "get {Node Name1} {Node Name2} ...",
	Long:  "get {Node Name1} {Node Name2} ...",

	RunE: func(cmd *cobra.Command, args []string) error {
		k, err := kubernetes.New()
		if err != nil {
			// TODO error handling
			return err
		}

		buildPods, err := k.BuildPods()
		if err != nil {
			// TODO error handling
			return err
		}
		specificNode := k.NodeFilter(buildPods, args...)
		buildIds := k.ExtractId(specificNode)

		sdconf, err := config.New()
		if err != nil {
			// TODO error handling
			return err
		}
		err = sdconf.Load()
		if err != nil {
			// TODO error handling
			return err
		}

		usedConf, err := sdconf.CurrentConfig()
		if err != nil {
			// TODO error handling
			return err
		}

		sd, err := screwdriver.New(usedConf.UserToken, usedConf.APIURL)
		if err != nil {
			return err
		}

		c := make(chan buildInfo)
		wg := &sync.WaitGroup{}

		for _, id := range buildIds {
			wg.Add(1)
			go sdRequest(id, sd, c, wg)
		}

		go func(c chan buildInfo, wg *sync.WaitGroup) {
			wg.Wait()
			close(c)
		}(c, wg)

		d := display.New()
		d.CreateHeader("Org/Repo", "JobName", "Build URL")
		var contents []display.RowContents
		for info := range c {
			u, err := url.Parse(usedConf.UIURL)
			if err != nil {
				return err
			}

			u.Path = path.Join(u.Path, "builds", info.buildID)
			row := d.CreateRowContents(info.pipeline.ScmRepo.Name, info.job.Name, u.String())
			contents = append(contents, row)
		}

		d.View(contents)

		return nil
	},
}

func sdRequest(id string, sd *screwdriver.SD, c chan buildInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	pack := buildInfo{}
	build, err := sd.Build(id)
	if err != nil {
		// TODO error handling
		fmt.Println(err)
		os.Exit(1)
	}
	pack.build = build

	job, err := sd.Job(pack.build.JobID)
	if err != nil {
		// TODO error handling
		fmt.Println(err)
		os.Exit(1)
	}
	pack.job = job

	pipeline, err := sd.Pipeline(pack.job.PipelineId)
	if err != nil {
		// TODO error handling
		fmt.Println(err)
		os.Exit(1)
	}
	pack.pipeline = pipeline

	pack.buildID = id
	c <- pack
}
