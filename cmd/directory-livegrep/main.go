package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// IndexConfig is a config of indexes
type IndexConfig struct {
	Name         string       `json:"name"`
	Repositories []RepoConfig `json:"repositories"`
}

// RepoConfig is a config of repos
type RepoConfig struct {
	Path      string            `json:"path"`
	Name      string            `json:"name"`
	Revisions []string          `json:"revisions"`
	Metadata  map[string]string `json:"metadata"`
}

// loadRepos finds every git directory under current and returns an array of RepoConfig
// for every bare repository
func loadRepos(path string) ([]string, error) {
	var repos []string

	err := filepath.Walk(path,
		func(p string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !(f.IsDir()) {
				return nil
			}
			if strings.HasPrefix(f.Name(), ".") {
				return nil
			}
			_, err = exec.Command("git", "--git-dir", p, "rev-parse", "--is-bare-repository").Output()
			if err != nil {
				return nil
			}
			fmt.Println("Adding", p, "to livegrep index")
			repos = append(repos, p)

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return repos, nil
}

func buildConfig(basepath string, repos []string) ([]byte, error) {
	cfg := IndexConfig{
		Name: "livegrep-config",
	}

	for _, r := range repos {
		url, err := exec.Command("git", "--git-dir", r, "config", "--get", "remote.origin.url").Output()
		if err != nil {
			fmt.Println("No remote origin URL defined for", r, "skipping...")
			continue
		}
		urlString := strings.TrimSuffix(string(url), filepath.Ext(string(url)))
		pathString := strings.TrimPrefix(r, basepath)
		cfg.Repositories = append(cfg.Repositories, RepoConfig{
			Path:      path.Join("/data", pathString),
			Name:      path.Base(r),
			Revisions: []string{"HEAD"},
			Metadata: map[string]string{
				"github": urlString,
			},
		})
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func writeConfig(config []byte, file string) error {
	dir := path.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(file, config, 0644)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("You must pass an absolute directory path as argument (without trailing slash)")
	}
	var path = os.Args[1]
	repos, err := loadRepos(path)
	if err != nil {
		log.Fatalln(err.Error())
	}

	config, err := buildConfig(path, repos)
	if err != nil {
		log.Fatalln(err.Error())
	}
	configPath := path + "/livegrep.json"
	if err := writeConfig(config, configPath); err != nil {
		log.Fatalln(err.Error())
	}

	err = ioutil.WriteFile(path+"/docker-compose.yaml", []byte(dockerCompose), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("")
	fmt.Println("---- Instructions for LiveGrep ----")
	fmt.Println("To build index at " + path + "/livegrep.idx:")
	fmt.Println("  docker run -v " + path + ":/data livegrep/indexer /livegrep/bin/codesearch -index_only -dump_index /data/livegrep.idx /data/livegrep.json")
	fmt.Println("")
	fmt.Println("To launch livegrep:")
	fmt.Println("  docker-compose -f " + path + "/docker-compose.yaml up")
}

const dockerCompose = `
version: "3.3"
services:
  livegrep-backend-linux:
    image: "docker.io/livegrep/base:latest"
    command:
      - "/livegrep/bin/codesearch"
      - "-grpc=0.0.0.0:9898"
      - "-load_index=/data/livegrep.idx"
    ports:
      - "9898:9898"
    volumes:
      - .:/data
    restart: unless-stopped
    networks:
      - livegrep

  livegrep-frontend:
    image: "docker.io/livegrep/base:latest"
    command:
      - "/livegrep/bin/livegrep"
      - "-docroot"
      - "/livegrep/web/"
      - "-connect"
      - "livegrep-backend-linux:9898"
      - "-listen"
      - "0.0.0.0:8910"
    ports:
      - "8910:8910"
    restart: unless-stopped
    networks:
      - livegrep

networks:
  livegrep:
`
