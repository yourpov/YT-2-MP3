package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"YT2MP3/utils"
)

const (
	reset = "\x1b[0m"
)

var option string

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Banner() string {
	return utils.Gradient(`
YT To MP3 ©

                 ┌───────────────────────────┐
                 │ ┓┏ ┏┳┓  ┏┳┓┏┓   ┳┳┓ ┏┓ ┏┓ │
                 │ ┗┫  ┃    ┃ ┃┃   ┃┃┃ ┣┛  ┫ │
                 │ ┗┛  ┻    ┻ ┗┛   ┛ ┗ ┃  ┗┛ │
                 │        Version 2.0        │
                 └───────────────────────────┘

 [1] Download Beat
 [2] Open Tunebat
 [3] Show Credits
─────────────────────────────────────────────────────────────────────
`, utils.MintyFresh) + reset
}

func download(link, path string) error {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return fmt.Errorf("yt-dlp is not installed. install it from https://github.com/yt-dlp/yt-dlp and place the /bin in /utils/ffmpeg")
	}

	ffmpegPath := filepath.Join("utils", "ffmpeg", "ffmpeg.exe")
	absFfmpegPath, err := filepath.Abs(ffmpegPath)
	if err != nil {
		return fmt.Errorf("could not find ffmpeg: %v", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("could not find downloads folder: %v", err)
	}

	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--ffmpeg-location", absFfmpegPath, "-o", absPath, link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	title := "Downloaded from: " + link + " by YT2MP3"
	metaCmd := exec.Command(absFfmpegPath, "-i", absPath, "-metadata", "title="+title, "-codec", "copy", absPath+"_temp.mp3")
	metaCmd.Stdout = os.Stdout
	metaCmd.Stderr = os.Stderr
	if err := metaCmd.Run(); err != nil {
		return err
	}

	if err := os.Rename(absPath+"_temp.mp3", absPath); err != nil {
		return err
	}

	return nil
}

func fetchRootFolder() string {
	dir, err := os.Getwd()
	if err != nil {
		return "Downloaded"
	}
	return filepath.Join(dir, "downloads")
}

func createFolder() string {
	downloadsPath := fetchRootFolder()
	if err := os.MkdirAll(downloadsPath, os.ModePerm); err != nil {
		fmt.Println(utils.Gradient("[!] Failed to create downloads folder.", utils.Candy))
		return "Downloaded"
	}
	return downloadsPath
}

func Credits() {
	Clear()
	fmt.Println(utils.Gradient(`
YT To MP3 ©

                 ┌───────────────────────────┐
                 │ ┓┏ ┏┳┓  ┏┳┓┏┓   ┳┳┓ ┏┓ ┏┓ │
                 │ ┗┫  ┃    ┃ ┃┃   ┃┃┃ ┣┛  ┫ │
                 │ ┗┛  ┻    ┻ ┗┛   ┛ ┗ ┃  ┗┛ │
                 │        Version 2.1        │
                 └───────────────────────────┘

Discord   : ItsJusNix.
Instagram : https://instagram.com/VanityVillian/
Telegram  : https://t.me/ItsJusNix
─────────────────────────────────────────────────────────────────────
`, utils.MintyFresh) + reset)
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(Banner())
		fmt.Print(utils.Gradient("Choose an option: ", utils.Candy) + " ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {

		case "1":
			Clear()
			fmt.Print(utils.Gradient("Enter YouTube link: ", utils.MintyFresh) + " ")
			link, _ := reader.ReadString('\n')
			link = strings.TrimSpace(link)

			fmt.Print(utils.Gradient("Save file as (without extension): ", utils.MintyFresh) + " ")

			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			downloadsPath := createFolder()
			savePath := filepath.Join(downloadsPath, fileName+".mp3")

			fmt.Println(utils.Gradient("[*] Downloading and converting to MP3...", utils.MintyFresh))

			if err := download(link, savePath); err != nil {
				fmt.Println(utils.Gradient(fmt.Sprintf("Error: %v", err), utils.Error))
				fmt.Println(utils.Gradient("Press ENTER to go back", utils.MintyFresh))
				fmt.Scanln()
				continue
			}

			fmt.Println(utils.Gradient("\n[✓] Download complete!", utils.Success))
			fmt.Println(utils.Gradient("Saved to: ", utils.MintyFresh) + savePath)

			if err := exec.Command("explorer", downloadsPath).Start(); err != nil {
				fmt.Println(utils.Gradient("[!] Failed to open folder.", utils.Error))
			}

			fmt.Print(utils.Gradient("Press ENTER to go back", utils.MintyFresh))
			fmt.Scanln()
			Clear()
			continue

		case "2":
			fmt.Println(utils.Gradient("Press ENTER to go back", utils.MintyFresh))
			if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "https://tunebat.com/Analyzer").Start(); err != nil {
				fmt.Println(utils.Gradient("[!] Failed to open Tunebat.", utils.Error))
			}
			fmt.Scanln()

		case "3":
			Clear()
			Credits()
			fmt.Println(utils.Gradient("Press ENTER to go back", utils.MintyFresh))
			fmt.Scanln()

		default:
			Clear()
			fmt.Println(utils.Gradient("Invalid Option.", utils.Error))
			fmt.Println(utils.Gradient("Press ENTER to go back", utils.MintyFresh))
			fmt.Scanln()
		}
		Clear()
	}
}
