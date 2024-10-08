package profile

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/liamg/github-profile-terminal-action/config"
	"github.com/liamg/github-profile-terminal-action/terminal"
	"github.com/liamg/github-profile-terminal-action/theme"
	"github.com/nfnt/resize"
)

//go:embed *.png
var embedded embed.FS

const (
	Width = 830
)

type Profile struct {
	dir    string
	config *config.Config
	gh     *github.Client
	theme  theme.Theme
	term   *terminal.Terminal
	stats  *Stats
}

func New(conf *config.Config) *Profile {
	return &Profile{
		config: conf,
		gh:     newGithubClient(conf),
		theme:  theme.ByName(conf.Theme),
	}
}

func (p *Profile) Generate(ctx context.Context, dir string) error {
	p.dir = dir

	repoBits := strings.Split(p.config.Context.Repository, "/")
	if len(repoBits) != 2 {
		return fmt.Errorf("invalid repository: %s", p.config.Context.Repository)
	}
	if repoBits[0] != repoBits[1] {
		return fmt.Errorf(
			"this action should be run on a special profile repository e.g. 'github.com/visitor/visitor', "+
				"'%s' does not appear to be such a repo", p.config.Context.Repository,
		)
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	p.term = terminal.New(Width, 600, nil, p.theme)

	if err := p.boot(); err != nil {
		return err
	}

	if err := p.login(ctx); err != nil {
		return err
	}

	if err := p.showStats(ctx); err != nil {
		return err
	}

	if err := p.term.ToGif(filepath.Join(p.dir, "os.gif"), true); err != nil {
		return err
	}

	return nil
}

func (p *Profile) boot() error {
	term := p.term

	f, err := embedded.Open("gh.png")
	if err != nil {
		return err
	}

	gh, err := png.Decode(f)
	if err != nil {
		return err
	}

	ratio := float64(gh.Bounds().Max.X) / float64(gh.Bounds().Max.Y)

	iw := 0.333 * float64(Width)
	ih := iw / ratio
	newImage := resize.Resize(uint(iw), uint(ih), gh, resize.Lanczos3)
	term.DrawImage(image.Rect(Width-(int(iw)+30), 30, Width, int(ih)+30), newImage)

	term.Frame(40)
	term.ShowCursor(false)

	term.Println("Release gifBIOS v7.3.4 - Build date 20/01/2031")
	term.Frame(20)
	term.Println("(C) 2022 GIF Systems Inc.\n\n\n")
	term.Frame(20)
	term.Println("GH Profile 0 Rev 1000")
	term.Frame(20)

	term.CursorToLastRow()
	term.Print("Press DEL to enter setup, ESC to skip memory test")
	term.CursorToRow(7)
	term.Frame(20)

	term.Println("Main Processor: GifCPU - 100Hz")
	term.Frame(20)
	for i := 0; i < 0x40000000; i += 0x4ffffff {
		term.ClearLine()
		term.CursorToHome()
		term.Print(fmt.Sprintf("Memory Check: %d", i))
		term.Frame(1)
	}
	term.ClearLine()
	term.CursorToHome()
	term.Println("Memory Check: 1048576K + 1024K Shared Memory\n")
	term.Frame(20)
	term.Println("WAIT...")
	term.Frame(100)

	term.Clear()
	term.Frame(100)

	term.Print("Starting GifOS...")
	term.Frame(150)
	term.Clear()

	return nil
}

func (p *Profile) login(_ context.Context) error {

	p.term.Clear()
	p.term.Println("GifOS v0.1.0 tty1")
	p.term.Println("")
	p.term.Print("login: ")
	p.term.ShowCursor(true)
	p.term.Frame(150)
	p.term.Type("visitor\n", terminal.Fast)

	p.term.Print("password: ")
	p.term.ShowCursor(true)
	p.term.Frame(200)
	p.term.Println("")
	p.term.Println("")
	p.term.Println(fmt.Sprintf("Last login %s on tty1", time.Now().Add(time.Hour*-24).Format(time.RFC1123)))
	p.term.Print(`Welcome to GifOS v0.1.0

  * Documentation: https://github.com/liamg/github-profile-terminal-action

0 packages can be updated.
0 updates are security updates.
`)
	p.term.Frame(50)

	return nil
}

func (p *Profile) prompt() {
	p.term.Print("\nvisitor@github:~$ ")
	p.term.ShowCursor(true)
	p.term.Frame(75)
}

func (p *Profile) showStats(ctx context.Context) error {

	p.prompt()

	p.term.Frame(75)

	p.term.Type("ls -la\n", terminal.Fast)
	p.term.Print(`drwxr-xr-x 10  visitor visitor 4.0K Mar 14 06:33 .
drwxr-xr-x  35 visitor visitor 4.0K Mar 11 06:17 ..
-rw-------   1 visitor visitor 2.7K Mar 14 09:03 .bash_history
-rw-r--r--   1 visitor visitor   21 Nov 21 19:31 .bash_logout
-rw-r--r--   1 visitor visitor   78 Dec 17 13:06 .bash_profile
-rw-r--r--   1 visitor visitor  609 Mar 14 20:19 .bashrc
drwxr-xr-x   3 visitor visitor 4.0K Jan 12 20:49 .bundle
drwxr-xr-x  21 visitor visitor 4.0K Mar  8 20:26 .cache
drwxr-xr-x   4 visitor visitor 4.0K Dec 17 13:08 .cargo
drwx------  27 visitor visitor 4.0K Mar 14 21:49 .config
drwxr-xr-x   2 visitor visitor 4.0K Dec 24 18:59 Desktop
drwxr-xr-x   2 visitor visitor 4.0K Mar 11 21:11 Downloads
-rw-r--r--   1 visitor visitor  398 Mar  3 21:43 .gitconfig
-rwx------   1 visitor visitor  239 Mar 11 22:58 ghlookup
-rw-r--r--   1 visitor visitor   14 Jan  3 10:51 .gitignore
drwx------   5 visitor visitor 4.0K Mar 14 21:56 .gnupg
drwxr-xr-x   4 visitor visitor 4.0K Dec 27 20:14 go
-rw-------   1 visitor visitor  15K Mar 15 20:06 .histfile
-rw-------   1 visitor visitor   20 Mar 14 21:52 .lesshst
drwx------   3 visitor visitor 4.0K Dec 17 19:15 .local
drwx------   3 visitor visitor 4.0K Dec 17 20:06 .pki
-rw-r--r--   1 visitor visitor   21 Dec 17 13:06 .profile
drwxr-xr-x   6 visitor visitor 4.0K Dec 17 13:07 .rustup
drwx------   2 visitor visitor 4.0K Mar  8 20:02 .ssh
drwxr-xr-x   2 visitor visitor 4.0K Jan  2 13:12 .vim
-rw-------   1 visitor visitor  16K Jan  3 15:12 .viminfo
drwxr-xr-x   3 visitor visitor 4.0K Dec 27 19:08 .vscode-oss
-rw-r--r--   1 visitor visitor  426 Jan  8 19:56 .zprofile
-rw-r--r--   1 visitor visitor  877 Mar 14 20:27 .zshaliases
-rw-r--r--   1 visitor visitor  600 Jan  3 18:29 .zshenv
-rw-------   1 visitor visitor  212 Jan  3 19:55 .zsh_history
-rw-r--r--   1 visitor visitor 1.2K Mar 14 20:37 .zshrc
`)
	p.prompt()
	p.term.Frame(75)

	stats, err := p.Stats(ctx)
	if err != nil {
		return err
	}

	user := stats.User

	p.term.Type(fmt.Sprintf("./ghlookup -u %s\n", user.GetLogin()), terminal.Fast)
	p.term.ShowCursor(false)

	p.term.Print("\nConnecting...")
	p.term.Type("............", terminal.VeryFast)
	p.term.Print("\nSending query...")
	p.term.Type("............................\n", terminal.VeryFast)

	p.term.Println("")
	p.term.Println("")
	p.term.Println(" -- User Details")
	p.term.Print("   User ID:     ")
	p.term.SetHighlight(true)
	p.term.Println(fmt.Sprintf("%d", user.GetID()))
	p.term.SetHighlight(false)
	p.term.Print("   Username:    ")
	p.term.SetHighlight(true)
	p.term.Println(user.GetLogin())
	p.term.SetHighlight(false)
	p.term.Print("   Real name:   ")
	p.term.SetHighlight(true)
	p.term.Println(user.GetName())
	p.term.SetHighlight(false)
	p.term.Print("   Location:    ")
	p.term.SetHighlight(true)
	p.term.Println(user.GetLocation())
	p.term.SetHighlight(false)
	p.term.Println("")
	p.term.Println("")
	p.term.Println("")
	p.term.Println("")
	p.term.Println("")

	p.term.Println(" -- Statistics")
	p.term.Print("   Total Stars:     ")
	p.term.SetHighlight(true)
	p.term.Println(fmt.Sprintf("%d", stats.TotalStars))
	p.term.SetHighlight(false)
	p.term.Print("   Total Followers: ")
	p.term.SetHighlight(true)
	p.term.Println(fmt.Sprintf("%d", stats.TotalFollowers))
	p.term.SetHighlight(false)
	p.term.Println("")
	p.term.Println("")
	p.term.Println("")

	resp, err := http.Get(stats.User.GetAvatarURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		img, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
	}

	ratio := float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y)
	iw := 0.25 * float64(Width)
	ih := iw / ratio
	resized := resize.Resize(uint(iw), uint(ih), img, resize.Lanczos3)
	p.term.DrawImage(image.Rect(Width-(int(iw)+30), 600-(int(ih)+30), Width, 600), resized)
	p.term.Frame(500)
	return nil
}
