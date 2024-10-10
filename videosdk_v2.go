package vidfusion

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

// 编码方式 h264_nvenc
const (
    Libx264   = "libx264"
    H264Nvenc = "h264_nvenc"
)

// Encoder 视频编码器
var Encoder = Libx264

const (
    FlipVideo    = "FlipVideo"    // 翻转视频
    SpeedUpVideo = "SpeedUpVideo" // 加速视频
    ScaleUpVideo = "ScaleUpVideo" // 放大视频
)

// VideoSDKV2 核心结构体
type VideoSDKV2 struct {
    CurrentFile string
    tempFiles   []string
    uniqueID    string // 唯一ID
}

// NewVideoSDKV2 创建 VideoSDKV2 实例
func NewVideoSDKV2(inputFile string) *VideoSDKV2 {
    return &VideoSDKV2{
        CurrentFile: inputFile,
        tempFiles:   []string{},
        uniqueID:    generateUniqueID(),
    }
}

// getNextTempFile 生成下一个临时文件路径
func (sdk *VideoSDKV2) getNextTempFile() string {
    tempFile, err := ioutil.TempFile("", sdk.uniqueID+"_video_*.mp4")
    if err != nil {
        panic(fmt.Sprintf("failed to create temp file: %v", err))
    }
    sdk.tempFiles = append(sdk.tempFiles, tempFile.Name())
    return tempFile.Name()
}

// Cleanup 清理临时文件
func (sdk *VideoSDKV2) Cleanup() {
    for _, file := range sdk.tempFiles {
        _ = os.Remove(file)
        fmt.Println("Removed temp file:", file)
    }
}

// CropVideoTimeline 裁剪视频时间线
func (sdk *VideoSDKV2) CropVideoTimeline(start, end float64) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-ss", fmt.Sprintf("%.2f", start), "-to", fmt.Sprintf("%.2f", end), outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to crop video: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// CropVideo 裁剪视频
func (sdk *VideoSDKV2) CropVideo(width, height int64) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-vf", fmt.Sprintf("scale=%d:%d", width, height), outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to crop video: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// FlipVideo 翻转视频
func (sdk *VideoSDKV2) FlipVideo() *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-vf", "hflip", "-c:v", Encoder, "-c:a", "copy", outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to flip video: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// SpeedUpVideo 加速视频
func (sdk *VideoSDKV2) SpeedUpVideo(speed float64) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-filter:v", fmt.Sprintf("setpts=%f*PTS", 1/speed), outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to speed up video: %v", err))
    }
    sdk.CurrentFile = outputFile
    duration, err := sdk.GetVideoDuration(outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to get video duration: %v", err))
    }
    // 调用裁剪视频时间线方法 缩减视频时长
    if speed > 0 {
        sdk.CropVideoTimeline(0, duration/speed)
    }
    return sdk
}

// ScaleUpVideo 放大视频
func (sdk *VideoSDKV2) ScaleUpVideo(scale float64) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-vf", fmt.Sprintf("scale=iw*%f:ih*%f", scale, scale), outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to scale up video: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// AddBackgroundMusic 添加背景音乐
func (sdk *VideoSDKV2) AddBackgroundMusic(audioFile string, volume float64) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg",
        "-i", sdk.CurrentFile,
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
    if err != nil {
        panic(fmt.Sprintf("failed to add background music: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// AddImageOverlay 添加图片水印
func (sdk *VideoSDKV2) AddImageOverlay(options OverlayOptions) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    // 使用传入的宽度和高度来缩放图片，并在指定位置进行覆盖
    filterComplex := fmt.Sprintf("[1:v]scale=%d:%d[img];[0:v][img]overlay=%d:%d", options.ImageWidth, options.ImageHeight, options.XPosition, options.YPosition)
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-i", options.ImageFile,
        "-filter_complex", filterComplex,
        "-c:v", Encoder, "-preset", "slow", "-crf", "23", "-c:a", "copy", outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to add image overlay: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// Mute 关闭视频原声
func (sdk *VideoSDKV2) Mute() *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand("ffmpeg", "-i", sdk.CurrentFile, "-an", "-c:v", "copy", "-c:a", "aac", outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to mute video: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// MuteTrack 添加静音轨道
func (sdk *VideoSDKV2) MuteTrack() *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
    err := runCommand(
        "ffmpeg", "-i", sdk.CurrentFile,
        "-f", "lavfi", "-i", "anullsrc=r=44100:cl=stereo", // 添加静音音轨
        "-c:v", "copy", // 保持视频编码不变
        "-c:a", "aac", // 使用 AAC 音频编码
        "-shortest", // 确保音轨与视频长度一致
        outputFile,
    )
    if err != nil {
        panic(fmt.Sprintf("failed to mute track: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// ConcatenateVideos 合并多个视频
func (sdk *VideoSDKV2) ConcatenateVideos(videoList []string, targetWidth, targetHeight int64) *VideoSDKV2 {
    tempFile, err := ioutil.TempFile("", "videos_*.txt")
    if err != nil {
        panic(fmt.Sprintf("failed to create temp file: %v", err))
    }
    defer os.Remove(tempFile.Name())
    defer tempFile.Close()
    
    tempDir, err := ioutil.TempDir("", "scaled_videos")
    if err != nil {
        panic(fmt.Sprintf("failed to create temp dir: %v", err))
    }
    defer os.RemoveAll(tempDir)
    
    for _, video := range videoList {
        scaledVideo := filepath.Join(tempDir, filepath.Base(video))
        scaleFilter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", targetWidth, targetHeight)
        err := runCommand("ffmpeg", "-i", video, "-vf", scaleFilter, "-r", "30", "-c:a", "copy", scaledVideo)
        if err != nil {
            panic(fmt.Sprintf("failed to scale video %s: %v", video, err))
        }
        _, err = tempFile.WriteString(fmt.Sprintf("file '%s'\n", scaledVideo))
        if err != nil {
            panic(fmt.Sprintf("failed to write to temp file: %v", err))
        }
    }
    
    outputFile := sdk.getNextTempFile()
    err = runCommand("ffmpeg", "-f", "concat", "-safe", "0", "-i", tempFile.Name(), "-c:v", Encoder, outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to concatenate videos: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// Finalize 最终生成文件，将临时文件复制到最终输出路径
func (sdk *VideoSDKV2) Finalize(outputFile string) string {
    tempFile := sdk.CurrentFile
    // 确保最终文件路径的目录存在
    if err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm); err != nil {
        panic(fmt.Sprintf("failed to create output directory: %v", err))
    }
    // 执行文件复制
    err := copyFile(tempFile, outputFile)
    if err != nil {
        panic(fmt.Sprintf("failed to finalize video: %v", err))
    }
    sdk.Cleanup()
    return outputFile
}

// GetMP3Duration 获取 MP3 文件时长
func (sdk *VideoSDKV2) GetMP3Duration(mp3File string) (float64, error) {
    return runCommandAndExtractFloat("ffprobe", "-i", mp3File, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
}

// GetVideoDuration 获取视频时长
func (sdk *VideoSDKV2) GetVideoDuration(videoFile string) (float64, error) {
    return runCommandAndExtractFloat("ffprobe", "-i", videoFile, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
}

// GetVideoDimensions 获取视频尺寸
func (sdk *VideoSDKV2) GetVideoDimensions(videoFile string) (int64, int64, error) {
    return runCommandAndExtractDimensions("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", videoFile)
}

// AddSubtitles 添加字幕并应用样式
func (sdk *VideoSDKV2) AddSubtitles(subtitleFile string, options SubtitleOptions) *VideoSDKV2 {
    outputFile := sdk.getNextTempFile()
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
    err := runCommand("ffmpeg",
        "-i", sdk.CurrentFile,
        "-i", subtitleFile,
        "-vf", fmt.Sprintf("subtitles='%s':force_style='%s'", subtitleFile, style), // 使用双引号
        "-c:v", Encoder,
        "-c:a", "aac",
        "-b:a", "192k",
        "-shortest",
        outputFile,
    )
    if err != nil {
        panic(fmt.Sprintf("failed to add subtitles: %v", err))
    }
    sdk.CurrentFile = outputFile
    return sdk
}

// GetRandomVideos 随机选择多个视频文件
func (sdk *VideoSDKV2) GetRandomVideos(directory string, numVideos int) ([]string, error) {
    return getRandomFiles(directory, "*.mp4", numVideos)
}

// VideosOptions 视频选项
type VideosOptions struct {
    VideoFile string
    Process   string
    Params    float64 // 放大视频的倍数 或者 加速视频的倍数
}

// ProcessVideosOptions 视频处理选项
type ProcessVideosOptions struct {
    VideoDuration float64         // 视频总时长
    Width         int64           // 视频宽度
    Height        int64           // 视频高度
    VideosOptions []VideosOptions // 视频选项
}

// ProcessVideos 封装方法 传入多个视频 时长 + 每个视频的处理方法 然后合并视频返回
// 视频处理方法 FlipVideo 翻转视频 SpeedUpVideo 加速视频 ScaleUpVideo 放大视频
func (sdk *VideoSDKV2) ProcessVideos(options ProcessVideosOptions) *VideoSDKV2 {
    var tempFiles []string
    // 根据视频选项, 先处理视频
    for _, videoOption := range options.VideosOptions {
        // 设置要处理的视频
        sdk.CurrentFile = videoOption.VideoFile
        // 按照设置的处理方法处理视频
        switch videoOption.Process {
        case "FlipVideo":
            sdk.FlipVideo()
        case "SpeedUpVideo":
            sdk.SpeedUpVideo(videoOption.Params)
        case "ScaleUpVideo":
            sdk.ScaleUpVideo(videoOption.Params)
        default:
            sdk.CurrentFile = videoOption.VideoFile
        }
        // 处理后的视频加入到临时文件列表
        tempFiles = append(tempFiles, sdk.CurrentFile)
    }
    // 开始拼合视频
    var execVideos []string
    var totalDuration float64
    for totalDuration < options.VideoDuration {
        for _, video := range tempFiles {
            videoDuration, err := sdk.GetVideoDuration(video)
            if err != nil {
                continue
            }
            totalDuration += videoDuration
            // 裁剪视频到指定尺寸
            sdk.CurrentFile = video
            sdk.CropVideo(options.Width, options.Height)
            execVideos = append(execVideos, sdk.CurrentFile)
            if totalDuration >= options.VideoDuration {
                break
            }
        }
    }
    sdk.ConcatenateVideos(execVideos, options.Width, options.Height)
    sdk.CropVideoTimeline(0, options.VideoDuration)
    return sdk
}
