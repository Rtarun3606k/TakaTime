package main

import (
	"github.com/Rtarun3606k/TakaTime/internal/Styles"
	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type dataLoadedMsg struct {
	updatedModel Model
	err          error
	FromCache    bool
}

// fetchData loads dashboard data asynchronously, resolving the active theme from
// SQLite config when available and falling back to fallbackTheme (the --theme flag value).
func fetchData(uri string, fallbackTheme types.ThemeConfig) tea.Cmd {
	return func() tea.Msg {
		// Initialize SQLite connection
		sqliteDB, err := db.InitSQLite()
		if err != nil {
			// Fallback if DB fails to open
			log.Println("sqlite DB connection initialization failed, falling back: ", err)
			return dataLoadedMsg{updatedModel: Model{}, err: err}
		}
		defer sqliteDB.Close()

		// Resolve the active theme: prefer the saved SQLite config, fall back to
		// the theme provided via the --theme flag.
		var activeTheme types.ThemeConfig
		userConfig, configErr := db.LoadConfig(sqliteDB)
		if configErr == nil && userConfig.Theme != "" {
			activeTheme = utils.ThemeSwitcher(userConfig.Theme)
		} else {
			activeTheme = fallbackTheme
		}
		loadedStyles := Styles.BuildStyles(activeTheme)

		// Check Cache First
		cachedData, err := db.GetDashboardCache(sqliteDB)
		if err == nil && cachedData != nil {
			// Cache HIT
			recolorList(cachedData.Languages, activeTheme)
			recolorList(cachedData.Projects, activeTheme)
			recolorList(cachedData.OS, activeTheme)
			recolorList(cachedData.Editors, activeTheme)

			tempModel := Model{
				AppStyles:         loadedStyles, // 👉 NEW 2: Attach styles to the Cache HIT model
				LanguageListStats: cachedData.Languages,
				ProjectListStats:  cachedData.Projects,
				OsListStats:       cachedData.OS,
				editorListStats:   cachedData.Editors,
				TimeStats:         cachedData.TimeStats,
				ActivityData:      cachedData.Activity,
				Streak:            cachedData.Streak,
				TodayHours:        cachedData.TodayHours,
				AverageHours:      cachedData.AverageHours,
				DailyHistory:      cachedData.DailyHistory,
			} 
			
			
			return dataLoadedMsg{updatedModel: tempModel, err: nil, FromCache: true}
		} else {
			log.Println("Failed to get dashboard Cache: ", err)
		}

		// Cache MISS! Fetch fresh from MongoDB
		tempModel := Model{TUITheme: activeTheme}
		filledModel, _, err := tempModel.GetData(uri)

		// 👉 NEW 3: Attach styles to the Cache MISS model before returning it
		filledModel.AppStyles = loadedStyles

		if err != nil {
			log.Println("Failed to fetch data from MongoDB:", err)
			return dataLoadedMsg{updatedModel: filledModel, err: err}
		}

		// Save to SQLite Cache for next time
		if sqliteDB != nil {
			db.SaveDashboardCache(sqliteDB, types.CacheData{
				Languages:    filledModel.LanguageListStats,
				Projects:     filledModel.ProjectListStats,
				OS:           filledModel.OsListStats,
				Editors:      filledModel.editorListStats,
				TimeStats:    filledModel.TimeStats,
				Activity:     filledModel.ActivityData,
				Streak:       filledModel.Streak,
				TodayHours:   filledModel.TodayHours,
				AverageHours: filledModel.AverageHours,
				DailyHistory: filledModel.DailyHistory,
			})
		}

		return dataLoadedMsg{
			updatedModel: filledModel,
			err:          nil,
			FromCache:    false,
		}
	}
}

// updateViewport swaps the content based on the active tab
func (m *Model) updateViewport() {
	if m.ActiveTab == "about" {
		m.Viewport.SetContent(m.generateAboutContent())
	} else {
		m.Viewport.SetContent(m.generateScrollableContent())
	}
}

// ----------------------------
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	//spinner

	case spinner.TickMsg:
		var spinCmd tea.Cmd
		m.Spinner, spinCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spinCmd)

		// If we are still loading, return immediately so the UI repaints the new frame
		if m.Loading {
			return m, tea.Batch(cmds...)
		}

		return m, cmd

	case tea.WindowSizeMsg:
		// Save the raw dimensions
		m.Width = msg.Width
		m.Height = msg.Height

		// We know our Header is roughly 3 lines, and Footer is 3 lines.
		// So the vertical margin we need to subtract from the viewport is 6.
		headerHeight := 6
		footerHeight := 3
		verticalMargin := headerHeight + footerHeight

		if !m.Ready {
			// 1. Initialize the viewport on the first render
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.Viewport.YPosition = headerHeight                  // Start it below the header
			m.Viewport.SetContent(m.generateScrollableContent()) // Helper func we will write
			m.Ready = true
		} else {
			// 2. If already ready, just resize it
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMargin
			// m.Viewport.SetContent(m.updateViewport())
		}
		m.updateViewport()
		return m, nil

	case dataLoadedMsg:
		m.Loading = false
		m.CacheData = msg.FromCache
		m.AppStyles = msg.updatedModel.AppStyles

		//  ASSIGN THE DATA HERE!
		if msg.err == nil {
			m.LanguageListStats = msg.updatedModel.LanguageListStats
			m.ProjectListStats = msg.updatedModel.ProjectListStats
			m.OsListStats = msg.updatedModel.OsListStats
			m.editorListStats = msg.updatedModel.editorListStats
			m.TimeStats = msg.updatedModel.TimeStats
			m.ActivityData = msg.updatedModel.ActivityData
			m.Streak = msg.updatedModel.Streak
			m.TodayHours = msg.updatedModel.TodayHours
			m.DailyHistory = msg.updatedModel.DailyHistory
		}

		// 3. Update the viewport content now that we have data!
		if m.Ready {
			// m.Viewport.SetContent(m.generateScrollableContent())
			m.updateViewport()
		}
		return m, nil

	case tea.KeyMsg:

		//settings model
		if m.ShowSettings {
			switch msg.String() {
			case "esc", "q", "s", "S":
				m.ShowSettings = false // Close modal without saving
				return m, nil

			// Navigation (Supports wrapping around columns!)
			case "up", "k":
				if m.SettingsCursor > 1 {
					m.SettingsCursor -= 2
				} // Jump up a row
				return m, nil
			case "down", "j":
				if m.SettingsCursor < len(types.AvailableThemes)-2 {
					m.SettingsCursor += 2
				} // Jump down a row
				return m, nil
			case "left", "h":
				if m.SettingsCursor > 0 {
					m.SettingsCursor--
				} // Move left
				return m, nil
			case "right", "l":
				if m.SettingsCursor < len(types.AvailableThemes)-1 {
					m.SettingsCursor++
				} // Move right
				return m, nil

			case "enter":
				// 1. Grab the name of the selected theme
				selectedThemeName := types.AvailableThemes[m.SettingsCursor]

				// 2. Fetch the new color palette using your function
				newThemeConfig := utils.ThemeSwitcher(selectedThemeName)

				// 3. Rebuild the AppStyles using the new colors!
				m.AppStyles = Styles.BuildStyles(newThemeConfig)
				m.TUITheme = newThemeConfig
				m.ShowSettings = false

				// Recolor lists using the new theme colors
				recolorList(m.LanguageListStats, newThemeConfig)
				recolorList(m.ProjectListStats, newThemeConfig)
				recolorList(m.OsListStats, newThemeConfig)
				recolorList(m.editorListStats, newThemeConfig)

				// Persist config and updated cache to SQLite
				sqliteDB, _ := db.InitSQLite()
				if sqliteDB != nil {
					db.SaveConfig(sqliteDB, types.CacheData{Theme: selectedThemeName})
					db.ClearDashboardCache(sqliteDB)
					// db.SaveDashboardCache(sqliteDB, types.CacheData{
					// 	Languages:    m.LanguageListStats,
					// 	Projects:     m.ProjectListStats,
					// 	OS:           m.OsListStats,
					// 	Editors:      m.editorListStats,
					// 	TimeStats:    m.TimeStats,
					// 	Activity:     m.ActivityData,
					// 	Streak:       m.Streak,
					// 	TodayHours:   m.TodayHours,
					// 	AverageHours: m.AverageHours,
					// 	DailyHistory: m.DailyHistory,
					// })
					sqliteDB.Close()
				}

				// 5. Force the viewport to redraw with the new colors!
				if m.Ready {
					m.updateViewport()
				}
				return m, nil
				return m, nil
			}
		}

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.Loading = true
			return m, fetchData(m.MongoURI, m.TUITheme)

		case "m":
			// Toggle the boolean (True becomes False, False becomes True)
			m.ViewMore = !m.ViewMore

			// Recalculate the viewport content with the new size!
			if m.Ready {
				// m.Viewport.SetContent(m.generateScrollableContent())
				m.updateViewport()
			}
			return m, nil

		case "h", "H":
			if m.ActiveTab != "home" {
				m.ActiveTab = "home"
				m.Viewport.GotoTop() // Reset scroll position when switching tabs
				m.updateViewport()
			}
			return m, nil

		case "a", "A":
			if m.ActiveTab != "about" {
				m.ActiveTab = "about"
				m.Viewport.GotoTop() // Reset scroll position when switching tabs
				m.updateViewport()
			}
			return m, nil

		case "s", "S": //  ADD THE HOTKEY TO OPEN SETTINGS!
			m.ShowSettings = true
			m.SettingsCursor = 0 // Reset cursor to top
			return m, nil

		}
	}

	// 4. Pass ANY unhandled messages (like scrolling keys) to the viewport!
	if m.Ready {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func recolorList(stats []types.ListStats, theme types.ThemeConfig) {
	colors := []string{theme.Color1, theme.Color2, theme.Color3, theme.Color4, theme.TextColor}
	for i := range stats {
		stats[i].Color = colors[i%len(colors)]
	}
}
