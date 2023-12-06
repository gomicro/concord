package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gomicro/concord/config"
	"github.com/gomicro/trust"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func init() {
	rootCmd.AddCommand(NewAuthCmd(os.Stdout, openBrowser))
}

const (
	state = "5bbeb48d-8f22-407a-9586-6da897ebffac"
)

var (
	clientID     string
	clientSecret string
)

func NewAuthCmd(out io.Writer, browserFunc func(string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth with github",
		Long:  `authorize concord against github`,
		RunE:  authRun(browserFunc),
	}

	cmd.Flags().BoolP("reauth", "r", false, "force concord to reauth")
	cmd.SetOut(out)

	return cmd
}

func authRun(browserFunc func(string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("auth: %w", err)
		}

		port := listener.Addr().(*net.TCPAddr).Port

		pool := trust.New()

		certs, err := pool.CACerts()
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("auth: %w", err)
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certs},
			},
		}

		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes: []string{
				"repo",      // Grants full access to public and private repositories including read and write access to code, commit statuses, repository invitations, collaborators, deployment statuses, and repository webhooks.
				"admin:org", // Fully manage the organization and its teams, projects, and memberships.
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			RedirectURL: fmt.Sprintf("http://localhost:%v/auth", port),
		}

		token := make(chan string)

		go startServer(ctx, listener, conf, token)

		var opts []oauth2.AuthCodeOption
		if strings.EqualFold(cmd.Flag("reauth").Value.String(), "true") {
			opts = []oauth2.AuthCodeOption{oauth2.AccessTypeOffline, oauth2.ApprovalForce}
		} else {
			opts = []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
		}

		url := conf.AuthCodeURL(state, opts...)

		err = browserFunc(url)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("auth: %w", err)
		}

		tkn := <-token
		close(token)

		c, err := config.ParseFromFile()
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("auth: %w", err)
		}

		c.Github.Token = tkn

		err = c.WriteToFile()
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("auth: %w", err)
		}

		return nil
	}
}

func startServer(ctx context.Context, listener net.Listener, conf *oauth2.Config, token chan string) {
	http.HandleFunc("/auth", authHandler(ctx, conf, token))

	srv := &http.Server{}

	go func() {
		<-ctx.Done()
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Printf("Error shutting down server: %s", err)
			os.Exit(1)
		}
	}()

	err := srv.Serve(listener)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func authHandler(ctx context.Context, conf *oauth2.Config, token chan string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		rstate := req.URL.Query().Get("state")

		if rstate != state {
			fmt.Println("bad response from oauth server")
			os.Exit(1)
		}

		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			fmt.Printf("errored exchanging token: %s", err)
			os.Exit(1)
		}

		body := `<html>
	<body>
		<h1>Config file updated, you can close this window.</h1>
	</body>
</html>`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body)) //nolint
		token <- tok.AccessToken
	}
}

func openBrowser(url string) error {
	ctx := context.Background()

	var perr error
	switch runtime.GOOS {
	case "linux":
		err := exec.CommandContext(ctx, "xdg-open", url).Start()
		if err != nil {
			perr = fmt.Errorf("open browser: linux: %w", err)
		}

	case "windows":
		err := exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", url).Start()
		if perr != nil {
			perr = fmt.Errorf("open browser: windows: %w", err)
		}

	case "darwin":
		err := exec.CommandContext(ctx, "open", url).Start()
		if perr != nil {
			perr = fmt.Errorf("open browser: darwin: %w", err)
		}

	default:
		perr = fmt.Errorf("open browser: unsupported platform")
	}

	return perr
}
