package screwdriver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	apiVersion        = "v4"
	authEndpoint      = "auth/token"
	buildsEndpoint    = "builds"
	jobsEndpoint      = "jobs"
	pipelinesEndpoint = "pipelines"
)

type SD struct {
	JWT        string
	baseAPIURL string
	client     *http.Client
}

type Job struct {
	Name       string `json:"name"`
	PipelineId int    `json:"pipelineId"`
}

type Build struct {
	JobID     int    `json:"jobId"`
	Container string `json:"container"`
}

type Pipeline struct {
	Name    string `json:"name"`
	ScmRepo struct {
		Name    string `json:"name"`
		Branch  string `json:"branch"`
		URL     string `json:"url"`
		RootDir string `json:"rootDir"`
	} `json:"scmRepo"`
}

type tokenResponse struct {
	JWT string `json:"token"`
}

func New(token, baseAPIURL string) (*SD, error) {
	sd := new(SD)
	sd.baseAPIURL = baseAPIURL
	sd.client = new(http.Client)
	jwt, err := sd.jwt(token)
	if err != nil {
		return sd, fmt.Errorf("failed to create instance: %w", err)
	}

	sd.JWT = jwt

	return sd, nil
}

func (sd *SD) request(method, url string, body io.Reader, isRequestAuth bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to request to Screwdriver.cd: %w", err)
	}

	if isRequestAuth {
		prefix := "Bearer "
		req.Header.Add("Authorization", prefix+sd.JWT)
		req.Header.Add("Content-Type", "application/json")
	}

	return sd.client.Do(req)
}

func (sd *SD) jwt(token string) (string, error) {
	apiUrl, err := sd.makeURL(authEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get jwt from Screwdriver.cd: %w", err)
	}

	q := apiUrl.Query()
	q.Set("api_token", token)
	apiUrl.RawQuery = q.Encode()
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, false)
	if err != nil {
		return "", fmt.Errorf("failed to get jwt from Screwdriver.cd: %w", err)
	}

	defer res.Body.Close()

	tr := new(tokenResponse)
	err = json.NewDecoder(res.Body).Decode(tr)
	if err != nil {
		return "", fmt.Errorf("failed to get jwt from Screwdriver.cd: %w", err)
	}

	return tr.JWT, nil
}

func (sd *SD) Build(id string) (*Build, error) {
	build := new(Build)
	apiUrl, err := sd.makeURL(buildsEndpoint)
	if err != nil {
		return build, fmt.Errorf("failed to get builds from Screwdriver.cd: %w", err)
	}

	apiUrl.Path = path.Join(apiUrl.Path, id)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		return build, fmt.Errorf("failed to get builds from Screwdriver.cd: %w", err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(build)
	if err != nil {
		return build, fmt.Errorf("failed to decode json for build: %w", err)
	}

	return build, nil
}

func (sd *SD) makeURL(endpoint string) (*url.URL, error) {
	u, err := url.Parse(sd.baseAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to API URL parse: %w", err)
	}

	u.Path = path.Join(u.Path, apiVersion, endpoint)

	return u, nil
}

func (sd *SD) Job(jobId int) (*Job, error) {
	job := new(Job)
	apiUrl, err := sd.makeURL(jobsEndpoint)
	if err != nil {
		return job, fmt.Errorf("failed to get job from Screwdriver.cd: %w", err)
	}

	apiUrl.Path = path.Join(apiUrl.Path, strconv.Itoa(jobId))
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		return job, fmt.Errorf("failed to get job from Screwdriver.cd: %w", err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(job)
	if err != nil {
		return job, fmt.Errorf("failed to get job from Screwdriver.cd: %w", err)
	}

	return job, nil
}

func (sd *SD) Pipeline(pipelineId int) (*Pipeline, error) {
	pipeline := new(Pipeline)
	apiUrl, err := sd.makeURL(pipelinesEndpoint)
	if err != nil {
		return pipeline, fmt.Errorf("failed to get pipeline from Screwdriver.cd: %w", err)
	}

	apiUrl.Path = path.Join(apiUrl.Path, strconv.Itoa(pipelineId))
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		return pipeline, fmt.Errorf("failed to get pipeline from Screwdriver.cd: %w", err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(pipeline)
	if err != nil {
		return pipeline, fmt.Errorf("failed to get pipeline from Screwdriver.cd: %w", err)
	}

	return pipeline, nil
}
