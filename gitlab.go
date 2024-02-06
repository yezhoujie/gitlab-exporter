package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/xanzy/go-gitlab"
)

/**
 * The project struct used to backup
 * 用于备份的项目结构体
 */
type BackUpProject struct {
	ProjectId   int    // 项目ID
	ProjectName string // 项目名称
	ProjectPath string // 项目路径
}

var ProjectList []BackUpProject

func getAllProject(gitlabClient *gitlab.Client, fromId int) {
	log.Println("listing all projects.....")

	// 获取所有的项目
	var options gitlab.ListProjectsOptions = gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10000,
			OrderBy: "id",
			Page:    1,
		},
	}

	for fetchProjects(gitlabClient, &options, fromId) {
		options.Page++
	}

}

/**
* 分页获取项目
* @return bool 是否有下一页
 */
func fetchProjects(gitlabClient *gitlab.Client, options *gitlab.ListProjectsOptions, fromId int) bool {
	// 获取所有的项目
	projects, response, err := gitlabClient.Projects.ListProjects(options, nil)
	if err != nil {
		log.Fatal("Failed to list projects: ", err)
	}
	if response.StatusCode != 200 {
		log.Fatal("Failed to list projects: ", response.Status)
	}
	fmt.Printf("totalProject：%v, perPage： %v, totalPages：%v, currentPage：%v, nextPage: %v\n", response.TotalItems, response.ItemsPerPage, response.TotalPages, response.CurrentPage, response.NextPage)

	for _, project := range projects {
		if fromId == -1 || project.ID <= fromId {
			ProjectList = append(ProjectList, BackUpProject{
				ProjectId:   project.ID,
				ProjectName: project.PathWithNamespace,
				ProjectPath: project.WebURL,
			})
		} else {
			log.Printf("skip project: %v, project id: %v\n", project.PathWithNamespace, project.ID)
		}
	}

	return response.NextPage != 0
	// return false
}

func backupProjectToOss(gitlabClient *gitlab.Client, keepLocalBackup bool, project BackUpProject, bucket *oss.Bucket) {

	// 调用gitlab 接口 直传oss
	res, err := gitlabClient.ProjectImportExport.ScheduleExport(project.ProjectId, nil)

	if err != nil {
		log.Fatal("Failed to schedule export: ", err)
	}
	if res.StatusCode != 202 {
		log.Fatal("Failed to schedule export: ", res.Status)
	}
	log.Println(res.Status)

	// 查看导出状态
	for !isFinished(gitlabClient, project.ProjectId) {
		time.Sleep(5 * time.Second)
	}
	time.Sleep(5 * time.Second)
	DownloadThenToOss(gitlabClient, project, bucket, keepLocalBackup)
}

func DownloadThenToOss(gitlabClient *gitlab.Client, project BackUpProject, bucket *oss.Bucket, keepLocalBackup bool) {

	// 调用gitlab 接口 直传oss
	data, res, err := gitlabClient.ProjectImportExport.ExportDownload(project.ProjectId)
	if err != nil {
		log.Fatal("Failed to export project: ", err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Failed to export project: ", res.Status)
	}
	// 把types[] data 存入文件
	folderName := "export/" + time.Now().Format("20060102") + "/" + strings.Split(project.ProjectName, "/")[0]
	projectName := folderName + "/" + strings.Split(project.ProjectName, "/")[1]
	exportFile := fmt.Sprintf("%v.tar.gz", projectName+"_export")
	err = os.MkdirAll(folderName, 0755)
	if err != nil {
		log.Fatal("Failed to create folder: ", err)
	}
	dataFile, err := os.Create(exportFile)
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	defer dataFile.Close()
	_, err = dataFile.Write(data)
	if err != nil {
		log.Fatal("Failed to write to file: ", err)
	}

	log.Printf("export download finished at: %v\n", exportFile)

	log.Println("start to upload to oss")
	upload(bucket, exportFile)
	log.Println("upload to oss finished")

	//删除本地备份
	if !keepLocalBackup {
		deleteLocalFile(exportFile)
	}
}

func deleteLocalFile(exportFile string) {
	err := os.Remove(exportFile)
	if err != nil {
		log.Fatal("Failed to delete file: ", err)
	}
}

func isFinished(gitlabClient *gitlab.Client, projectId int) bool {
	// 查看导出状态
	exportStatus, res, err := gitlabClient.ProjectImportExport.ExportStatus(projectId)
	if err != nil {
		log.Fatal("Failed to get export status: ", err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Failed to get export status: ", res.Status)
	}
	// log.Println(exportStatus)
	return exportStatus.ExportStatus == "finished"
}
