package teamcity

import (
	"fmt"
	"log"
	"sort"
	"time"

	dash "github.com/itglobal/dashboard/api"

	"github.com/kapitanov/go-teamcity"
)

const key = "teamcity"

type teamcityData map[string]*dash.Item

type teamcityProvider struct {
	callback       dash.Callback
	teamcityClient teamcity.Client
	fetchInterval  time.Duration
	data           teamcityData
}

func (p *teamcityProvider) Key() string {
	return key
}

func factory(config dash.Config, callback dash.Callback) (dash.Provider, error) {
	url, err := config.GetString("url")
	if err != nil {
		return nil, err
	}

	username := config.GetStringOrDefault("username", "")
	password := config.GetStringOrDefault("password", "")

	fetchInterval, err := time.ParseDuration(config.GetStringOrDefault("timer", "20s"))

	p := new(teamcityProvider)
	p.callback = callback

	var auth teamcity.Authorizer
	if username != "" && password != "" {
		auth = teamcity.BasicAuth(username, password)
	} else {
		auth = teamcity.GuestAuth()
	}

	p.teamcityClient = teamcity.NewClient(url, auth)
	p.fetchInterval = fetchInterval
	p.data = make(teamcityData)

	go p.run()

	return p, nil
}

func init() {
	dash.RegisterFactory(key, factory)
}

func (p *teamcityProvider) run() {
	p.fetchProjects()
	p.fetchBuilds()
	p.fetchUnknownItems()
	p.syncItems()

	for {
		time.Sleep(p.fetchInterval)

		res, err := p.fetchBuilds()
		if err != nil {
			continue
		}

		if res == frGotNewProjects {
			p.fetchProjects()
			p.fetchBuilds()
			p.fetchUnknownItems()
		}

		p.syncItems()
	}
}

type fetchBuildsResult int

const (
	frOK fetchBuildsResult = iota
	frError
	frGotNewProjects
)

const fetchBuildsCount = 25

func (p *teamcityProvider) fetchBuilds() (fetchBuildsResult, error) {
	log.Printf("[teamcity] fetchBuilds")

	// Fetch N last builds and put them into data container
	builds, err := p.teamcityClient.GetBuilds(fetchBuildsCount)
	if err != nil {
		log.Printf("[teamcity] fetchBuilds -> frError %s", err)
		return frError, err
	}

	// Scanning list of build in reverse order
	for i := len(builds) - 1; i >= 0; i-- {
		build := builds[i]
		//for _, build := range builds {
		item, exists := p.data[build.BuildTypeID]
		if !exists {
			// Got a project that was not fetched yet, need to refetch list of the projects
			log.Printf("[teamcity] fetchBuilds -> frGotNewProjects (buildTypeId = '%s')", build.BuildTypeID)
			return frGotNewProjects, nil
		}

		setupDashItem(&build, item)
	}

	log.Printf("[teamcity] fetchBuilds -> frOK")
	return frOK, nil
}

func (p *teamcityProvider) fetchProjects() {
	log.Printf("[teamcity] fetchProjects")

	// Fetch list of projects
	projects, err := p.teamcityClient.GetProjects()
	if err != nil {
		log.Printf("[teamcity] fetchProjects -> error! %s", err)
		return
	}

	// For each project - fetch a list of build types
	for _, project := range projects {
		p.fetchBuildTypes(&project)
	}
}

func (p *teamcityProvider) fetchBuildTypes(project *teamcity.Project) {
	log.Printf("[teamcity] fetchBuildTypes(%s)", project.ID)

	// Fetch a list of build types for a projects
	buildTypes, err := p.teamcityClient.GetBuildTypesForProject(project.ID)
	if err != nil {
		log.Printf("[teamcity] fetchBuildTypes(%s) -> error! %s", project.ID, err)
		return
	}

	// For each build type - add a dashboard item if not exists
	for _, buildType := range buildTypes {
		item, exists := p.data[buildType.ID]
		if exists {
			continue
		}

		item = &dash.Item{}
		item.Key = buildType.ID
		item.Name = buildType.Name
		item.ProviderKey = key
		item.Status = dash.StatusUnknown
		item.StatusText = ""
		item.Progress = 0

		p.data[buildType.ID] = item

		log.Printf("[teamcity] fetchBuildTypes(%s): new dash item %s", project.ID, buildType.ID)
	}
}

func (p *teamcityProvider) fetchUnknownItems() {
	for id, item := range p.data {
		if item.Status == dash.StatusUnknown {
			p.fetchItemStatus(id, item)
		}
	}
}

func (p *teamcityProvider) fetchItemStatus(id string, item *dash.Item) {
	builds, err := p.teamcityClient.GetBuildsForBuildType(id, 1)
	if err != nil {
		log.Printf("[teamcity] fetchItemStatus(%s) -> error! %s", id, err)
		return
	}

	l := len(builds)
	if l == 0 {
		log.Printf("[teamcity] fetchItemStatus(%s) -> no builds are found", id)
		item.Status = dash.StatusBad
		item.StatusText = "No build are found"
		return
	}

	setupDashItem(&builds[0], item)
}

func (p *teamcityProvider) syncItems() {
	items := make([]*dash.Item, len(p.data))
	keys := make([]string, len(p.data))

	i := 0
	for key := range p.data {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	i = 0
	for _, key := range keys {
		items[i] = p.data[key]
		i++
	}

	p.callback(p, items)
}

func setupDashItem(build *teamcity.Build, item *dash.Item) {
	if build.StatusText != "" {
		item.StatusText = fmt.Sprintf("%s: %s", build.Number, build.StatusText)
	} else {
		item.StatusText = build.Number
	}

	switch build.Status {
	case teamcity.StatusSuccess:
		item.Status = dash.StatusGood
		item.Progress = dash.NoProgress
		break
	case teamcity.StatusRunning:
		item.Status = dash.StatusPending
		item.Progress = build.Progress
		break
	case teamcity.StatusFailure:
		item.Status = dash.StatusBad
		item.Progress = dash.NoProgress
		break
	}
}
