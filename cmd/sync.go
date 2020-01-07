package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"kingcity.app/tools/sync-mvn-deps/utils"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	all             bool
	projectsPattern string
	projects        []string
	dependency      string
	version         string
	shouldReplace   bool
	shouldShowDiff  bool
	shouldGitAddPom bool
	shouldMvnUpdate bool
	shouldBuild     bool
)

func replace(filePath string, newVersion string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	content := string(bytes)

	pattern := fmt.Sprintf(`(<artifactId>%v</artifactId>[\n\s]+<version>).+(</version>)`, dependency)
	reg := regexp.MustCompile(pattern)
	newContent := reg.ReplaceAllString(content, "${1}"+newVersion+"${2}")
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}
}

func showDiff(servicePath string) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = servicePath
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(string(out))
}

func gitAddPom(servicePath string) {
	cmd := exec.Command("git", "add", "pom.xml")
	cmd.Dir = servicePath
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(string(out))
}

func mvnUpdate(servicePath string) {
	cmd := exec.Command("mvn", "-U", "clean", "install")
	cmd.Dir = servicePath
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(string(out))
}

func build(servicePath string) {
	cmd := exec.Command("./build.sh")
	cmd.Dir = servicePath
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(string(out))
}

func action(projects []string) {
	for _, project := range projects {
		pomPath := project + "/pom.xml"
		if utils.Exists(pomPath) {
			println("=========================================" + pomPath)
			if shouldReplace {
				replace(pomPath, version)
			}
			if shouldShowDiff {
				showDiff(project)
			}
			if shouldGitAddPom {
				gitAddPom(project)
			}
			if shouldMvnUpdate {
				mvnUpdate(project)
			}
			if shouldBuild {
				build(project)
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "sync-mvn-deps",
	Short: "Sync specified dependency to specified version",
	Long:  `Sync specified dependency to specified version`,
	Run: func(cmd *cobra.Command, args []string) {
		if shouldReplace && (len(dependency) == 0 || len(version) == 0) {
			_ = cmd.Help()
			return
		}
		if !all && len(projects) == 0 {
			_ = cmd.Help()
			return
		}
		if all && len(projectsPattern) == 0 {
			_ = cmd.Help()
			return
		}
		if all {
			var err error
			projects, err = filepath.Glob(projectsPattern)
			if err != nil {
				panic(err)
			}
		}
		action(projects)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&projectsPattern, "pattern", "p", "", "path pattern of target projects")
	rootCmd.Flags().StringArrayVarP(&projects, "projects", "j", []string{}, "specified projects path to sync")
	rootCmd.Flags().StringVarP(&dependency, "dependency", "d", "", "dependency name")
	rootCmd.Flags().StringVarP(&version, "version", "v", "", "new version of dependency to sync")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "replace all services or not")

	rootCmd.Flags().BoolVarP(&shouldReplace, "replace", "r", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldShowDiff, "showdiff", "s", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldGitAddPom, "gitaddpom", "g", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldMvnUpdate, "mvnupdate", "u", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldBuild, "build", "b", false, "replace all services or not")
}
