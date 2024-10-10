package vidfusion

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "runtime"
)

// runCommand 通用命令执行函数
func runCommand(name string, args ...string) error {
    // 获取调用者的函数名
    pc, _, _, ok := runtime.Caller(1)
    caller := "unknown"
    if ok {
        caller = runtime.FuncForPC(pc).Name()
    }
    fmt.Println("Caller:", caller)
    
    // 添加 -y 参数，以便在文件已存在时自动覆盖
    args = append([]string{"-y"}, args...)
    // 如果使用 H264Nvenc 编码器，则添加 -hwaccel cuda 参数进行解码硬件加速
    if Encoder == H264Nvenc {
        args = append([]string{"-hwaccel", "cuda"}, args...)
    }
    cmd := exec.Command(name, args...)
    fmt.Printf("Running command: %v\n", cmd.String())
    cmdOutput, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("command error: %v\noutput: %s", err, string(cmdOutput))
    }
    // fmt.Printf("Command output: %s\n", string(cmdOutput))
    return nil
}

// runCommandAndExtractFloat 从命令结果中提取 float 值
func runCommandAndExtractFloat(name string, args ...string) (float64, error) {
    cmd := exec.Command(name, args...)
    fmt.Printf("Running command: %v\n", cmd.String())
    output, err := cmd.Output()
    if err != nil {
        return 0, fmt.Errorf("command error: %v", err)
    }
    
    var result float64
    _, err = fmt.Sscanf(string(output), "%f", &result)
    if err != nil {
        return 0, err
    }
    return result, nil
}

// runCommandAndExtractDimensions 从命令中提取宽度和高度
func runCommandAndExtractDimensions(name string, args ...string) (int64, int64, error) {
    cmd := exec.Command(name, args...)
    fmt.Printf("Running command: %v\n", cmd.String())
    output, err := cmd.Output()
    if err != nil {
        return 0, 0, fmt.Errorf("command error: %v", err)
    }
    
    var width, height int
    fmt.Sscanf(string(output), "%dx%d", &width, &height)
    return int64(width), int64(height), nil
}

// copyFile 将源文件复制到目标文件
func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    // 使用 io.Copy 将源文件内容复制到目标文件
    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return err
    }
    
    // 确保复制操作完成后，刷新写入缓存
    err = destFile.Sync()
    if err != nil {
        return err
    }
    
    return nil
}
