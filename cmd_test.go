package vidfusion

import "testing"

// Test_runCommand 测试 runCommand 函数
func Test_runCommand(t *testing.T) {
    err := runCommand("ffmpeg", "-version")
    if err != nil {
        t.Errorf("runCommand error: %v", err)
    }
}

// Benchmark_runCommand 基准测试 runCommand 函数
func Benchmark_runCommand(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = runCommand("ffmpeg", "-version")
    }
}

// Test_runCommandAndExtractFloat 测试 runCommandAndExtractFloat 函数
func Test_runCommandAndExtractFloat(t *testing.T) {
    result, err := runCommandAndExtractFloat("ffmpeg", "-version")
    if err != nil {
        t.Errorf("runCommandAndExtractFloat error: %v", err)
    }
    t.Logf("ffprobe version: %f", result)
}

// Benchmark_runCommandAndExtractFloat 基准测试 runCommandAndExtractFloat 函数
func Benchmark_runCommandAndExtractFloat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = runCommandAndExtractFloat("ffmpeg", "-version")
    }
}
