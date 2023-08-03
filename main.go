package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var version = "0.1.0"

func main() {
	var hostname string
	var httpMethod string
	var showVersion bool
	var authBasic string
	var authFile string
	var jsonKey string
	var warning string
	var critical string
	var timeout int32
	var header string
	var debug bool
	var insecure bool

	var rootCmd = &cobra.Command{
		Use:   "check_rest",
		Short: "CLI application to check API status",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("check_rest version %s\n", version)
			} else {
				resp, err := http.Get(hostname)
				if err != nil {
					fmt.Printf("CRITICAL - error while calling the url: %s\n", err)
					os.Exit(2)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("CRITICAL - Error reading body: %s", err)
					os.Exit(2)
				}
				result := gjson.ParseBytes(body)
				count := result.Get(jsonKey)

				if !count.Exists() {
					fmt.Printf("UNKNOWN - The key '%s' not found\n", jsonKey)
					os.Exit(3)
				}

				fmt.Printf("%s = '%s'\n", jsonKey, count.String())
				os.Exit(0)
			}
		},
	}

	rootCmd.Flags().BoolVarP(&showVersion, "version", "V", false, "Print version information")
	rootCmd.Flags().StringVarP(&authBasic, "auth-basic", "b", "", "Uses HTTP Basic authentication with provided <username>:<password>")
	rootCmd.Flags().StringVarP(&authFile, "auth-file", "f", "", "Uses HTTP Basic authentication and takes the required 'username:password' from the file provided. This file should only have one line that contains text in the format <username>:<password>")
	rootCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "The hostname or IP address of the API you want to check")
	rootCmd.Flags().StringVarP(&httpMethod, "http-method", "m", "GET", "The HTTP method for the API call. Supported methods: GET, POST, PUT")
	rootCmd.Flags().StringVarP(&jsonKey, "key", "K", "", "The json key to check. If not provided check_rest will check the HTTP status code.")
	rootCmd.Flags().StringVarP(&warning, "warning", "w", "", "Warning thresholds")
	rootCmd.Flags().StringVarP(&critical, "critical", "c", "", "Critical threshholds")
	rootCmd.Flags().Int32VarP(&timeout, "timeout", "t", 10, "Seconds before connection times out")
	rootCmd.Flags().StringVarP(&header, "header", "D", "", "HTTP header")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&insecure, "insecure", "k", false, "Disables checking SSL certificate (if using SSL/HTTPS)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
