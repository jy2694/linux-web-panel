package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const rootDir = "/" // 최상위 디렉토리(보안상 실제 서비스시 제한 필요)

// 디렉토리/파일 목록 조회
func ListFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = rootDir
	}
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path))
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "목록을 불러올 수 없어요."})
		return
	}
	var list []gin.H
	for _, f := range files {
		list = append(list, gin.H{
			"name":     f.Name(),
			"is_dir":   f.IsDir(),
			"size":     f.Size(),
			"mod_time": f.ModTime(),
		})
	}
	c.JSON(http.StatusOK, list)
}

// 파일/폴더 생성
func CreateFileOrDir(c *gin.Context) {
	var req struct {
		Path  string `json:"path"`
		IsDir bool   `json:"is_dir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청"})
		return
	}
	absPath := filepath.Join(rootDir, filepath.Clean("/"+req.Path))
	if req.IsDir {
		err := os.MkdirAll(absPath, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "폴더 생성 실패"})
			return
		}
	} else {
		f, err := os.Create(absPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 생성 실패"})
			return
		}
		f.Close()
	}
	c.JSON(http.StatusOK, gin.H{"message": "생성됨"})
}

// 파일/폴더 삭제
func DeleteFileOrDir(c *gin.Context) {
	path := c.Query("path")
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path))
	err := os.RemoveAll(absPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "삭제 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "삭제됨"})
}

// 이름 변경
func RenameFileOrDir(c *gin.Context) {
	var req struct {
		OldPath string `json:"old_path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청"})
		return
	}
	oldAbs := filepath.Join(rootDir, filepath.Clean("/"+req.OldPath))
	newAbs := filepath.Join(rootDir, filepath.Clean("/"+req.NewPath))
	err := os.Rename(oldAbs, newAbs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "이름 변경 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "이름 변경됨"})
}

// 파일 내용 읽기
func ReadFile(c *gin.Context) {
	path := c.Query("path")
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path))
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 읽기 실패"})
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
}

// 파일 내용 쓰기
func WriteFile(c *gin.Context) {
	path := c.Query("path")
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path))
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데이터"})
		return
	}
	err = ioutil.WriteFile(absPath, data, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 쓰기 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "저장됨"})
}

// 파일 업로드
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "업로드 파일 없음"})
		return
	}
	path := c.PostForm("path")
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path), file.Filename)
	if err := c.SaveUploadedFile(file, absPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "업로드 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "업로드됨"})
}

// 파일 다운로드
func DownloadFile(c *gin.Context) {
	path := c.Query("path")
	absPath := filepath.Join(rootDir, filepath.Clean("/"+path))
	c.FileAttachment(absPath, filepath.Base(absPath))
}
