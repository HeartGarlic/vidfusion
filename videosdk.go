package vidfusion

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

// VideoSDK 核心结构体
type VideoSDK struct{}

// NewVideoSDK 创建 VideoSDK 实例
func NewVideoSDK() *VideoSDK {
    return &VideoSDK{}
}

// GetMP3Duration 获取 MP3 文件时长
func (sdk *VideoSDK) GetMP3Duration(mp3File string) (float64, error) {
    return runCommandAndExtractFloat("ffprobe", "-i", mp3File, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
}

// GetVideoDuration 获取视频时长
func (sdk *VideoSDK) GetVideoDuration(videoFile string) (float64, error) {
    return runCommandAndExtractFloat("ffprobe", "-i", videoFile, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
}

// CropVideoTimeline 裁剪视频时间线
func (sdk *VideoSDK) CropVideoTimeline(inputFile, outputFile string, start, end float64) error {
    return runCommand("ffmpeg", "-i", inputFile, "-ss", fmt.Sprintf("%.2f", start), "-to", fmt.Sprintf("%.2f", end), outputFile)
}

// CropVideo 裁剪视频
func (sdk *VideoSDK) CropVideo(inputFile, outputFile string, width, height int64) error {
    return runCommand("ffmpeg", "-i", inputFile, "-vf", fmt.Sprintf("scale=%d:%d", width, height), outputFile)
}

// FlipVideo 翻转视频
func (sdk *VideoSDK) FlipVideo(inputFile, outputFile string) error {
    return runCommand("ffmpeg", "-i", inputFile, "-vf", "hflip", "-c:v", "libx264", "-c:a", "copy", outputFile)
}

// SpeedUpVideo 加速视频
func (sdk *VideoSDK) SpeedUpVideo(inputFile, outputFile string, speed float64) error {
    return runCommand("ffmpeg", "-i", inputFile, "-filter:v", fmt.Sprintf("setpts=%f*PTS", 1/speed), outputFile)
}

// ScaleUpVideo 放大视频
func (sdk *VideoSDK) ScaleUpVideo(inputFile, outputFile string, scale float64) error {
    return runCommand("ffmpeg", "-i", inputFile, "-vf", fmt.Sprintf("scale=iw*%f:ih*%f", scale, scale), outputFile)
}

// AddBackgroundMusic 添加背景音乐
func (sdk *VideoSDK) AddBackgroundMusic(videoFile, audioFile, outputFile string, volume float64) error {
    // 构建 ffmpeg 命令
    return runCommand("ffmpeg",
        "-i", videoFile,
        "-i", audioFile,
        "-filter_complex", fmt.Sprintf("[1:a]volume=%.1f[a1];[0:a][a1]amix=inputs=2:duration=first:dropout_transition=2[a]", volume),
        "-map", "0:v", // 选择视频流
        "-map", "[a]", // 选择混合后的音频流
        "-c:v", "copy", // 复制视频流
        "-c:a", "aac", // 使用 AAC 编码音频
        "-b:a", "192k", // 设置音频比特率
        "-shortest", // 输出文件最短持续时间
        outputFile,
    )
}

// GetVideoDimensions 获取视频的宽度和高度
func (sdk *VideoSDK) GetVideoDimensions(videoFile string) (int64, int64, error) {
    return runCommandAndExtractDimensions("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", videoFile)
}

// OverlayOptions 用于配置图片覆盖选项
type OverlayOptions struct {
    ImageWidth  int64  // 图片宽度
    ImageHeight int64  // 图片高度
    XPosition   int64  // 图片起始的 X 坐标
    YPosition   int64  // 图片起始的 Y 坐标
    VideoFile   string // 视频文件路径
    ImageFile   string // 图片文件路径
    OutputFile  string // 输出文件路径
}

// AddImageOverlay 添加图片覆盖，支持设置图片大小、起始位置
func (sdk *VideoSDK) AddImageOverlay(options OverlayOptions) error {
    // 使用传入的宽度和高度来缩放图片，并在指定位置进行覆盖
    filterComplex := fmt.Sprintf("[1:v]scale=%d:%d[img];[0:v][img]overlay=%d:%d",
        options.ImageWidth, options.ImageHeight, options.XPosition, options.YPosition)
    
    return runCommand("ffmpeg", "-i", options.VideoFile, "-i", options.ImageFile,
        "-filter_complex", filterComplex,
        "-c:v", "libx264", "-preset", "slow", "-crf", "23", "-c:a", "copy", options.OutputFile)
}

// ConcatenateVideos 合并多个视频，在视频尺寸不一致时，强制合并
func (sdk *VideoSDK) ConcatenateVideos(videoList []string, outputFile string, targetWidth, targetHeight int) error {
    // 创建一个唯一的临时文件
    tempFile, err := ioutil.TempFile("", "videos_*.txt")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %v", err)
    }
    defer os.Remove(tempFile.Name()) // 确保在函数结束时删除临时文件
    defer tempFile.Close()
    
    // 创建一个临时目录存储缩放后的视频
    tempDir, err := ioutil.TempDir("", "scaled_videos")
    if err != nil {
        return fmt.Errorf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // 确保临时目录在完成后被删除
    
    // 逐个视频进行缩放处理
    for _, video := range videoList {
        // 为每个视频生成一个新的输出文件名
        scaledVideo := filepath.Join(tempDir, filepath.Base(video))
        
        // 使用 ffmpeg 缩放视频到指定尺寸并统一帧率
        // 如果视频尺寸大于目标尺寸，则裁剪
        scaleFilter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", targetWidth, targetHeight)
        err := runCommand(
            "ffmpeg",
            "-i", video,
            "-vf", scaleFilter,
            "-r", "30",
            "-c:a", "copy",
            scaledVideo,
        )
        if err != nil {
            return fmt.Errorf("failed to scale video %s: %v", video, err)
        }
        // 将处理后的视频路径写入临时文件
        _, err = tempFile.WriteString(fmt.Sprintf("file '%s'\n", scaledVideo))
        if err != nil {
            return fmt.Errorf("failed to write to temp file: %v", err)
        }
    }
    
    // 使用缩放后的视频进行合并，不带入视频原声，强制编码
    return runCommand("ffmpeg", "-f", "concat", "-safe", "0", "-i", tempFile.Name(), "-c:v", "libx264", outputFile)
}

// MuteVideo 关闭视频原声
func (sdk *VideoSDK) MuteVideo(videoFile, outputFile string) error {
    return runCommand(
        "ffmpeg",
        "-i", videoFile,
        "-an", // 移除视频中的所有音轨
        //"-f", "lavfi", "-t", "1", "-i", "anullsrc=r=44100:cl=stereo", // 添加静音音轨
        "-shortest",    // 确保音频与视频长度匹配
        "-c:v", "copy", // 保持视频编码不变
        "-c:a", "aac", // 使用 AAC 音频编码
        outputFile,
    )
}

// MuteTrack 添加静音音轨
func (sdk *VideoSDK) MuteTrack(videoFile, outputFile string) error {
    return runCommand(
        "ffmpeg",
        "-i", videoFile,
        "-f", "lavfi", "-i", "anullsrc=r=44100:cl=stereo", // 添加静音音轨
        "-c:v", "copy", // 保持视频编码不变
        "-c:a", "aac", // 使用 AAC 音频编码
        "-shortest", // 确保音轨与视频长度一致
        outputFile,
    )
}

// SubtitleOptions 用于配置字幕样式的选项
type SubtitleOptions struct {
    FontSize  int64  // 字号
    XPosition int64  // 左右位置（负值表示靠右）
    YPosition int64  // 上下位置（负值表示靠下）
    Font      string // 字体
    FontColor string // 字体颜色
    Alignment int64  // 对齐方式（1 左对齐，2 居中对齐，3 右对齐）
}

// AddSubtitles 添加字幕并应用样式
func (sdk *VideoSDK) AddSubtitles(videoFile, subtitleFile, outputFile string, options SubtitleOptions) error {
    // 构建 ffmpeg 中的 force_style 字符串，转义必要字符
    style := fmt.Sprintf(
        "Alignment=%d,Fontsize=%d,PrimaryColour=&H%s&,FontName=%s,MarginL=%d,MarginR=%d,MarginV=%d",
        options.Alignment,
        options.FontSize,
        options.FontColor,
        options.Font,
        options.XPosition,
        options.XPosition, // 如果想要两个边距相同，可以重复使用
        options.YPosition,
    )
    
    // 使用转义后的文件路径和样式
    return runCommand("ffmpeg",
        "-i", videoFile,
        "-i", subtitleFile,
        "-vf", fmt.Sprintf("subtitles='%s':force_style='%s'", subtitleFile, style), // 使用双引号
        "-c:v", "libx264",
        "-c:a", "aac",
        "-b:a", "192k",
        "-shortest",
        outputFile,
    )
}

// GetRandomVideos 随机选择多个视频文件
func (sdk *VideoSDK) GetRandomVideos(directory string, numVideos int) ([]string, error) {
    return getRandomFiles(directory, "*.mp4", numVideos)
}
