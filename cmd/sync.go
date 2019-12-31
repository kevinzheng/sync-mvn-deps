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
	project         string
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

	pattern := fmt.Sprintf(`(<artifactId>%v</artifactId>[\n\s]+<version>)[0-9\.-]+(</version>)`, dependency)
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

func action(service string) {
	pomPath := service + "/pom.xml"
	if utils.Exists(pomPath) {
		println("=========================================" + pomPath)
		if shouldReplace {
			replace(pomPath, version)
		}
		if shouldShowDiff {
			showDiff(service)
		}
		if shouldGitAddPom {
			gitAddPom(service)
		}
		if shouldMvnUpdate {
			mvnUpdate(service)
		}
		if shouldBuild {
			build(service)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "sync-mvn-deps",
	Short: "Sync specified dependency to specified version",
	Long:  `Sync specified dependency to specified version`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(dependency) == 0 || len(version) == 0 {
			_ = cmd.Help()
			return
		}
		if !all && len(project) == 0 {
			_ = cmd.Help()
			return
		}
		if all && len(projectsPattern) == 0 {
			_ = cmd.Help()
			return
		}
		if all {
			services, err := filepath.Glob(projectsPattern)
			if err != nil {
				panic(err)
			}
			for _, service := range services {
				action(service)
			}
		} else {
			action(project)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&projectsPattern, "pattern", "p", "../../../*-service", "path pattern of target projects")
	rootCmd.Flags().StringVarP(&project, "project", "j", "", "specified project path to sync")
	rootCmd.Flags().StringVarP(&dependency, "dependency", "d", "", "dependency name")
	rootCmd.Flags().StringVarP(&version, "version", "v", "", "new version of dependency to sync")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "replace all services or not")

	rootCmd.Flags().BoolVarP(&shouldReplace, "replace", "r", true, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldShowDiff, "showdiff", "s", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldGitAddPom, "gitaddpom", "g", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldMvnUpdate, "mvnupdate", "u", false, "replace all services or not")
	rootCmd.Flags().BoolVarP(&shouldBuild, "build", "b", false, "replace all services or not")
}
