package vidfusion

import (
    "math/rand"
    "path/filepath"
    "strings"
    "time"
)

// 获取随机视频文件列表
func getRandomFiles(directory, pattern string, numFiles int) ([]string, error) {
    rand.Seed(time.Now().UnixNano())
    files, err := filepath.Glob(filepath.Join(directory, pattern))
    if err != nil {
        return nil, err
    }
    
    if len(files) < numFiles {
        numFiles = len(files)
    }
    
    rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
    return files[:numFiles], nil
}

// escapeFilePath 转义文件路径，确保路径正确传递给 ffmpeg
func escapeFilePath(filePath string) string {
    // 将 Windows 反斜杠替换为正斜杠，或者使用 filepath.ToSlash()
    return strings.Replace(filePath, "\\", "/", -1)
}

// randomString 生成指定长度的随机字符串
func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// generateUniqueID 生成唯一 ID
func generateUniqueID() string {
    return time.Now().Format("20060102150405") + "_" + randomString(5)
}
