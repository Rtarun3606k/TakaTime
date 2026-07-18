package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	buildimg "github.com/Rtarun3606k/TakaTime/internal/BuildImg"
	dbqueryv2 "github.com/Rtarun3606k/TakaTime/internal/DBQueryV2"
	gogist "github.com/Rtarun3606k/TakaTime/internal/GoGist"
	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"github.com/Rtarun3606k/TakaTime/internal/types"
)

//go:embed FiraCodeNerdFontPropo-Retina.ttf
var fontData []byte

func main() {

	var themeFlag string
	var customBg, customText, customSubText, customBarBg string
	var customC1, customC2, customC3, customC4 string
	// add veresion flag
	var versionFlag bool

	//----------flags
	days := flag.Int("days", 0, "Number of past days to include (0 = today)")
	flag.StringVar(&themeFlag, "theme", "dark", "Base theme: dark, light, dracula, nord, gruvbox, monokai, cyberpunk")
	//parse version flag
	flag.BoolVar(&versionFlag, "version", false, "show version")
	// Custom Overrides
	flag.StringVar(&customBg, "bg", "", "Override Background Color")
	flag.StringVar(&customText, "text", "", "Override Text Color")
	flag.StringVar(&customSubText, "subtext", "", "Override Sub-Text Color")
	flag.StringVar(&customBarBg, "bar-bg", "", "Override Bar Background Color")

	flag.StringVar(&customC1, "c1", "", "Override Color 1 (Primary)")
	flag.StringVar(&customC2, "c2", "", "Override Color 2")
	flag.StringVar(&customC3, "c3", "", "Override Color 3")
	flag.StringVar(&customC4, "c4", "", "Override Color 4")

	flag.Parse() //  Parse flags first!
	log.Println("days as a flag is discontinued please stop adding it as a flag. remove it from your work flow file check out readme https://github.com/Rtarun3606k/Takatime for more information ", days)

	if versionFlag {
		fmt.Println(types.Version)
		return
	}
	// 2. Initialize the Final 'theme' Variable
	var theme types.ThemeConfig

	// Decide the base: Start with a full "Light" or "Dark" template
	theme = utils.ThemeSwitcher(themeFlag)

	// 3. Apply Overrides
	if customBg != "" {
		theme.BackgroundColor = customBg
	}
	if customText != "" {
		theme.TextColor = customText
	}
	if customSubText != "" {
		theme.SubTextColor = customSubText
	}
	if customBarBg != "" {
		theme.BarBackgroundColor = customBarBg
	}

	if customC1 != "" {
		theme.Color1 = customC1
	}
	if customC2 != "" {
		theme.Color2 = customC2
	}
	if customC3 != "" {
		theme.Color3 = customC3
	}
	if customC4 != "" {
		theme.Color4 = customC4
	}

	// Connect
	mongoURI := os.Getenv("MONGO_URI")
	// gistID := os.Getenv("GIST_ID")
	gistToken := os.Getenv("GIST_TOKEN")
	targetRepo := os.Getenv("TARGET_REPO")

	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	client, err := db.ConnectToDataBase(mongoURI)
	if err != nil {
		log.Printf("error in connecting to database : %v", err)
		return
	}

	defer client.Disconnect(context.TODO())

	if gistToken != "" && targetRepo != "" {
		fmt.Printf("Updating README for %s...\n", targetRepo)

		if err := dbqueryv2.RunMigrations(client); err != nil {
			log.Printf("Migration warning: %v", err)
			// Don't crash, just log it. The report can still run.
		}

		if err := dbqueryv2.RunLanguageMigration(client, "takatime", "logs"); err != nil {
			log.Printf("Migration warning: %v", err)

		}

		name := strings.Split(targetRepo, "/")

		// 1. Fetch Real Data for all time all time days = 0 else the count
		projects, err := dbqueryv2.GetListStats(client, "project", 5, theme, 0)
		if err != nil {
			log.Println("Proj Error:", err)
		}

		langs, err := dbqueryv2.GetListStats(client, "language", 5, theme, 0)
		if err != nil {
			log.Println("Lang Error:", err)
		}

		editors, err := dbqueryv2.GetListStats(client, "editor", 4, theme, 0)
		if err != nil {
			log.Println("Editor Error:", err)
		}

		osStats, err := dbqueryv2.GetListStats(client, "os", 3, theme, 0)
		if err != nil {
			log.Println("OS Error:", err)
		}

		timeStats, err := dbqueryv2.GetTimeStats(client)
		if err != nil {
			log.Println("Time Error:", err)
		}

		// past 30 days
		projects30, err := dbqueryv2.GetListStats(client, "project", 5, theme, 30)
		if err != nil {
			log.Println("Proj Error:", err)
		}

		langs30, err := dbqueryv2.GetListStats(client, "language", 5, theme, 30)
		if err != nil {
			log.Println("Lang Error:", err)
		}

		context, cancle := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancle()
		collection := client.Database("takatime").Collection("logs")

		_, _, _, dailyHistory, err := dbqueryv2.FetchStreakAndToday(context, collection)
		if err != nil {
			log.Printf("Error fetching streak: %v", err)
		}

		// job 1 : Languages
		utils.HandleImageJob("Top Languages - All Time", "public/taka-languages.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.DrawListCard("Top Languages - All Time", langs, fontData, time.Now(), theme, false)
		})

		// Job 2: Projects
		utils.HandleImageJob("Top Projects - All Time", "public/taka-projects.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.DrawListCard("Top Projects - All Time", projects, fontData, time.Now(), theme, true)
		})

		// Job 3: Time Grid (2x2 View)
		utils.HandleImageJob("Time Stats", "public/taka-time.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.DrawTimeCard(timeStats, fontData, time.Now(), theme, name[0])
		})

		// Job 4: Tech Stack (Editors Left / OS Right)
		utils.HandleImageJob("environment stats", "public/taka-tech.png", gistToken, targetRepo, func() (image.Image, error) {
			// Pass both lists to the split-view generator
			return buildimg.DrawTechCard(editors, osStats, fontData, time.Now(), theme)
		})

		// job 5 : Languages
		utils.HandleImageJob("Top Language - Last 30 Days", "public/taka-languages30.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.DrawListCard("Top Language - Last 30 Days", langs30, fontData, time.Now(), theme, true)
		})

		// Job 6: Projects
		utils.HandleImageJob("Top Projects - Last 30 Days", "public/taka-projects30.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.DrawListCard("Top Projects - Last 30 Days", projects30, fontData, time.Now(), theme, true)
		})

		// Job 7: heatmap generation
		utils.HandleImageJob("Heatmap", "public/taka-heatmap.png", gistToken, targetRepo, func() (image.Image, error) {
			return buildimg.HeatmapStatsImg(dailyHistory, 800, fontData, theme, 2.0)
		})

		content := utils.GenerateOutput()

		errr := gogist.UpdateReadMe(gistToken, targetRepo, content)
		if errr != nil {
			fmt.Println("Some error occured while updating readme ", errr)
		}
		fmt.Println("README Updated Successfully!")
	} else {
		fmt.Println("Skipping README update (GIST_TOKEN or TARGET_REPO missing)")
	}

	//
	// updateContentError := gogist.UpdateGist(gistToken, gistID, content)
	// if updateContentError != nil {
	// 	log.Fatalln("Some error occured gogist", updateContentError)
	// 	return
	// }
	//

}
