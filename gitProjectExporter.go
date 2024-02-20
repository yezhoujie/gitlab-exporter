package main

// 创建一个方法调用gitlab的API
import (
	"flag"
	"fmt"
	"log"

	pd "github.com/elulcao/progress-bar/cmd"
	"github.com/xanzy/go-gitlab"
)

func main() {
	// 接收命令行参数
	filePath := flag.String("f", "config.yaml", "path of config file/配置文件路径")
	helpFlag := flag.Bool("h", false, "Display help")
	fromId := flag.Int("fromId", -1, "export project from project a id desc")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		println("example: ./gitlab-backup -f ~/config.yaml ")
		return
	}

	log.Println("gitlab backup starting.....")
	//load config
	log.Println("loading config file.....")

	if *fromId == -1 {
		log.Println("going to export all projects")
	} else {
		log.Printf("going to only export projects <= id: %v\n", *fromId)
	}

	config := loadConfig(*filePath)

	// 创建gitlab client
	if config.Gitlab.Url == "" {
		log.Fatal("gitlab url is empty")
	}
	if config.Gitlab.Token == "" {
		log.Fatal("gitlab token is empty")
	}

	// 创建gitlab client
	gitlabClient := getGitLabClent(config)

	// 获取所有的项目
	getAllProject(gitlabClient, *fromId)
	fmt.Printf("get %v projects to backup\n", len(ProjectList))
	bucket := getOssBucket(config)

	bar := pd.NewPBar()
	bar.Total = uint16(len(ProjectList))
	// 备份项目
	for i, project := range ProjectList {
		bar.RenderPBar(i + 1)
		// 将当前时间yyyymmddhhmm作为备份文件夹名称
		// folderName := time.Now().Format("20060102")
		// projectName := folderName + "/" + project.ProjectName
		// 获取oss签名URL
		// signUrl := getOssSign(bucket, "test")
		// log.Println("get oss preSign url: ", signUrl)
		// 备份项目
		log.Printf("start to backup project: %v, project path : %v, project id: %v\n", project.ProjectName, project.ProjectPath, project.ProjectId)
		backupProjectToOss(gitlabClient, config.KeepLocalBackup, project, bucket)
		// break
	}
	log.Println("gitlab backup finished.....")
}

func getGitLabClent(config Config) *gitlab.Client {
	gitlabClient, err := gitlab.NewClient(config.Gitlab.Token, gitlab.WithBaseURL(config.Gitlab.Url))
	if err != nil {
		log.Fatal("Failed to create gitlab client: ", err)
	}
	return gitlabClient
}
