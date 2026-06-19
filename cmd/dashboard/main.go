package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Rtarun3606k/TakaTime/internal/Styles"
	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/debugger"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

func initModel(mongoURI string, theme types.ThemeConfig) Model {

	//new spinner
	s := spinner.New()

	s.Spinner = spinner.Dot // other options spinner.Line, spinner.MiniDot, or spinner.Points

	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1))
	return Model{
		Loading:   true,
		Err:       nil,
		MongoURI:  mongoURI, // Save this so Init() can use it
		TUITheme:  theme,
		AppStyles: Styles.InitStyles(theme),
		Spinner:   s,
		ActiveTab: "home", // default Home
	}
}

func (m Model) Init() tea.Cmd {
	// If no URI was provided, don't try to fetch
	if m.MongoURI == "" {
		return m.Spinner.Tick
	}
	// Start spinning the spinner and tell Bubble Tea to fetch data in the background!
	return tea.Batch(m.Spinner.Tick, fetchData(m.MongoURI, m.TUITheme))
}

func main() {
	//logging setup
	f, err := debugger.SetupDashboardLog()
	if err != nil {
		fmt.Printf("fatal error setting up the debugger: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var mongoDBString string
	var themeFlag string
	var versionFlag bool

	flag.BoolVar(&versionFlag, "version", false, "show version")

	flag.StringVar(&mongoDBString, "MongoDBString", "", "MongoDB String for dashboard to query data")

	flag.StringVar(&themeFlag, "theme", "dark", "Base theme: dark, light, dracula, nord, gruvbox, monokai, cyberpunk")
	flag.Parse()

	if versionFlag {
		fmt.Println(types.Version)
		return
	}
	//get the theme
	theme := utils.ThemeSwitcher(themeFlag)

	// Initialize the model with the URI. It will start with Loading: true
	appModel := initModel(mongoDBString, theme)
	log.Println("appModel initialized")

	// Start the program instantly. The Init() function handles the async loading.
	p := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error has occured while running the dashboard: %v", err)
		os.Exit(1)
	}
}
