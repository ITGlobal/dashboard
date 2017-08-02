package teamcity

import (
	"fmt"
	"time"

	"github.com/itglobal/dashboard/tile"
	"github.com/kapitanov/go-teamcity"
	log "github.com/kpango/glg"
	uuid "github.com/satori/go.uuid"
)

type teamcityDataItem struct {
	ID              tile.ID
	BuildTypeID     string
	Title           string
	State           tile.State
	DescriptionText string
	HasProgress     bool
	Progress        int
	IsVisible       bool
}

type teamcityData map[string]*teamcityDataItem

type teamcityProvider struct {
	uid                string
	teamcityClient     teamcity.Client
	manager            tile.Manager
	fetchInterval      time.Duration
	itemsByID          map[tile.ID]*teamcityDataItem
	itemsByBuildTypeID map[string]*teamcityDataItem
}

// Gets provider unique ID
func (p *teamcityProvider) ID() string {
	return p.uid
}

// Gets provider type key
func (p *teamcityProvider) Type() string {
	return providerType
}

// Initializes a provider
func (p *teamcityProvider) Init() error {
	// TODO drop old build types if they are not available va API
	go func() {
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
	}()
	return nil
}

type fetchBuildsResult int

const (
	frOK fetchBuildsResult = iota
	frError
	frGotNewProjects
)

const fetchBuildsCount = 25

func (p *teamcityProvider) fetchBuilds() (fetchBuildsResult, error) {
	// Fetch N last builds and put them into data container
	builds, err := p.teamcityClient.GetBuilds(fetchBuildsCount)
	if err != nil {
		log.Errorf("[teamcity] fetchBuilds -> frError %s", err)
		return frError, err
	}

	// Scanning list of build in reverse order
	for i := len(builds) - 1; i >= 0; i-- {
		build := builds[i]
		item, exists := p.itemsByBuildTypeID[build.BuildTypeID]
		if !exists {
			// Got a project that was not fetched yet, need to refetch list of the projects
			log.Debugf("[teamcity] fetchBuilds -> frGotNewProjects (buildTypeId = '%s')", build.BuildTypeID)
			return frGotNewProjects, nil
		}

		item.IsVisible = true
		setupDashItem(&build, item)
	}

	return frOK, nil
}

func (p *teamcityProvider) fetchProjects() {
	// Fetch list of projects
	projects, err := p.teamcityClient.GetProjects()
	if err != nil {
		log.Errorf("[teamcity] fetchProjects -> error! %s", err)
		return
	}

	// For each project - fetch a list of build types
	for _, project := range projects {
		p.fetchBuildTypes(&project)
	}
}

func (p *teamcityProvider) fetchBuildTypes(project *teamcity.Project) {
	// Fetch a list of build types for a projects
	buildTypes, err := p.teamcityClient.GetBuildTypesForProject(project.ID)
	if err != nil {
		log.Errorf("[teamcity] fetchBuildTypes(%s) -> error! %s", project.ID, err)
		return
	}

	// For each build type - add a dashboard item if not exists
	for _, buildType := range buildTypes {
		item, exists := p.itemsByBuildTypeID[buildType.ID]
		if exists {
			continue
		}

		item = &teamcityDataItem{
			ID:          tile.ID(uuid.NewV4().String()),
			BuildTypeID: buildType.ID,
			Title:       buildType.Name,
			State:       tile.StateDefault,
		}

		p.itemsByID[item.ID] = item
		p.itemsByBuildTypeID[buildType.ID] = item

		log.Debugf("[teamcity] fetchBuildTypes(%s): new item %s", project.ID, buildType.ID)
	}
}

func (p *teamcityProvider) fetchUnknownItems() {
	for _, item := range p.itemsByBuildTypeID {
		if item.State == tile.StateDefault {
			p.fetchItemStatus(item)
		}
	}
}

func (p *teamcityProvider) fetchItemStatus(item *teamcityDataItem) {
	builds, err := p.teamcityClient.GetBuildsForBuildType(item.BuildTypeID, 1)
	if err != nil {
		log.Errorf("[teamcity] fetchItemStatus(%s) -> error! %s", item.BuildTypeID, err)
		return
	}

	l := len(builds)
	if l == 0 {
		log.Warnf("[teamcity] fetchItemStatus(%s) -> no builds are found", item.BuildTypeID)
		item.State = tile.StateDefault
		item.DescriptionText = "No build are found"
		item.IsVisible = false
		return
	}

	item.IsVisible = true
	setupDashItem(&builds[0], item)
}

func (p *teamcityProvider) syncItems() {
	u := p.manager.BeginUpdate(p)
	defer u.EndUpdate()

	for _, t := range u.GetTiles() {
		if _, exists := p.itemsByID[t.ID()]; !exists {
			u.RemoveTile(t.ID())
		}
	}

	for _, item := range p.itemsByID {
		if !item.IsVisible {
			u.RemoveTile(item.ID)
			continue
		}

		t := u.AddOrUpdateTile(item.ID)

		t.SetType(tile.TypeTextStatusProgress)
		t.SetSize(tile.Size2x)
		t.SetState(item.State)
		t.SetTitleText(item.Title)
		t.SetDescriptionText(item.DescriptionText)
		if item.HasProgress {
			t.SetStatusValue(item.Progress)
		} else {
			t.SetNoStatusValue()
		}
	}
}

func setupDashItem(build *teamcity.Build, item *teamcityDataItem) {
	if build.StatusText != "" {
		item.DescriptionText = fmt.Sprintf("%s: %s", build.Number, build.StatusText)
	} else {
		item.DescriptionText = build.Number
	}

	item.IsVisible = true

	switch build.Status {
	case teamcity.StatusSuccess:
		item.State = tile.StateSuccess
		item.HasProgress = false
		break
	case teamcity.StatusRunning:
		item.State = tile.StateIndeterminate
		item.Progress = build.Progress
		item.HasProgress = true
		break
	case teamcity.StatusFailure:
		item.State = tile.StateError
		item.HasProgress = false
		break
	}
}
