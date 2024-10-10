package vidfusion

import (
    "os"
    "testing"
)

// TestVideoSDK_GetMP3Duration 测试获取 MP3 文件时长
func TestVideoSDK_GetMP3Duration(t *testing.T) {
    sdk := NewVideoSDK()
    duration, err := sdk.GetMP3Duration("D:\\work\\go\\vid-fusion-demo\\mp3.wav")
    if err != nil {
        t.Error(err)
    }
    t.Log(duration)
}

// BenchmarkVideoSDK_GetMP3Duration 基准测试获取 MP3 文件时长
func BenchmarkVideoSDK_GetMP3Duration(b *testing.B) {
    sdk := NewVideoSDK()
    for i := 0; i < b.N; i++ {
        _, _ = sdk.GetMP3Duration("D:\\work\\go\\vid-fusion-demo\\mp3.wav")
    }
}

// TestVideoSDK_GetVideoDuration 测试获取视频时长
func TestVideoSDK_GetVideoDuration(t *testing.T) {
    sdk := NewVideoSDK()
    duration, err := sdk.GetVideoDuration("C:\\Users\\Administrator\\Desktop\\视频\\1.mp4")
    if err != nil {
        t.Error(err)
    }
    t.Log(duration)
}

// TestVideoSDK_CropVideoTimeline 测试裁剪视频时间线
func TestVideoSDK_CropVideoTimeline(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.CropVideoTimeline("C:\\Users\\Administrator\\Desktop\\视频\\1.mp4", "C:\\Users\\Administrator\\Desktop\\视频\\1_copy.mp4", 0, 5)
    if err != nil {
        t.Error(err)
    }
}

// TestVideoSDK_CropVideo 测试裁剪视频
func TestVideoSDK_CropVideo(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.CropVideo("C:\\Users\\Administrator\\Desktop\\视频\\1.mp4", "C:\\Users\\Administrator\\Desktop\\视频\\1_copy.mp4", 720, 1280)
    if err != nil {
        t.Error(err)
    }
}

// TestVideoSDK_FlipVideo 测试翻转视频
func TestVideoSDK_FlipVideo(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.FlipVideo("C:\\Users\\Administrator\\Desktop\\视频\\1.mp4", "C:\\Users\\Administrator\\Desktop\\视频\\1_copy.mp4")
    if err != nil {
        t.Error(err)
    }
}

// TestVideoSDK_AddBackgroundMusic 测试添加背景音乐
func TestVideoSDK_AddBackgroundMusic(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.AddBackgroundMusic(
        "C:\\Users\\Administrator\\Desktop\\视频\\1.mp4",
        "C:\\Users\\Administrator\\Desktop\\视频\\music.flac",
        "C:\\Users\\Administrator\\Desktop\\视频\\1_copy.mp4", 1)
    if err != nil {
        t.Error(err)
    }
}

// TestVideoSDK_AddImageOverlay 测试添加图片覆盖
func TestVideoSDK_AddImageOverlay(t *testing.T) {
    sdk := NewVideoSDK()
    width, height, err := sdk.GetVideoDimensions("C:\\Users\\Administrator\\Desktop\\视频\\1.mp4")
    if err != nil {
        t.Error(err)
    }
    err = sdk.AddImageOverlay(OverlayOptions{
        ImageWidth:  width,
        ImageHeight: height,
        XPosition:   0,
        YPosition:   0,
        VideoFile:   "C:\\Users\\Administrator\\Desktop\\视频\\1.mp4",
        ImageFile:   "C:\\Users\\Administrator\\Desktop\\视频\\overlay.png",
        OutputFile:  "C:\\Users\\Administrator\\Desktop\\视频\\1_copy.mp4",
    })
    if err != nil {
        t.Error(err)
    }
}

// TestVideoSDK_ConcatenateVideos 测试合并多个视频
func TestVideoSDK_ConcatenateVideos(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.ConcatenateVideos([]string{
        "C:\\Users\\Administrator\\Desktop\\视频\\1.mp4",
        "C:\\Users\\Administrator\\Desktop\\视频\\2.mp4",
    }, "C:\\Users\\Administrator\\Desktop\\视频\\output.mp4", 720, 1280)
    if err != nil {
        t.Error(err)
    }
}

// 添加字幕
func TestVideoSDK_AddSubtitles(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.AddSubtitles(
        "C:\\Users\\Administrator\\Desktop\\视频\\output6.mp4",
        "srt.srt",
        "C:\\Users\\Administrator\\Desktop\\视频\\4_copy.mp4",
        SubtitleOptions{
            FontSize:  12,
            FontColor: "00FFFFFF",
            XPosition: 0,
            YPosition: 0,
            Alignment: 2,
            Font:      "Arial",
        },
    )
    if err != nil {
        t.Error(err)
    }
}

// 关闭原生并添加静音轨道
func TestVideoSDK_MuteVideo(t *testing.T) {
    sdk := NewVideoSDK()
    err := sdk.MuteVideo(
        "C:\\Users\\Administrator\\Desktop\\视频\\5.mp4",
        "C:\\Users\\Administrator\\Desktop\\视频\\5_copy.mp4")
    if err != nil {
        t.Error(err)
    }
    err = sdk.MuteTrack(
        "C:\\Users\\Administrator\\Desktop\\视频\\5_copy.mp4",
        "C:\\Users\\Administrator\\Desktop\\视频\\5_copy_copy.mp4")
    if err != nil {
        t.Error(err)
    }
    // 添加背景音乐
    err = sdk.AddBackgroundMusic(
        "C:\\Users\\Administrator\\Desktop\\视频\\5_copy_copy.mp4",
        "C:\\Users\\Administrator\\Desktop\\视频\\music.flac",
        "C:\\Users\\Administrator\\Desktop\\视频\\5_copy_copy_5_copy_copy.mp4", 1)
}

// 完整的调用示例
func TestVideoSDK_CompleteExample(t *testing.T) {
    basePath := "C:\\Users\\Administrator\\Desktop\\视频\\"
    var clearFile = make([]string, 0)
    // 创建 VideoSDK 实例
    sdk := NewVideoSDK()
    // 获取mp3文件时长
    duration, err := sdk.GetMP3Duration(basePath + "mp3.wav")
    if err != nil {
        t.Error(err)
    }
    // 获取视频文件时长
    videos, err := sdk.GetRandomVideos(basePath, 10)
    if err != nil {
        return
    }
    var execVideos []string
    // 获取视频时长 累计时长大于mp3时长 就退出循环
    // 修改循环 使视频时长累计大于mp3时长
    var totalDuration float64
    for totalDuration < duration {
        for _, video := range videos {
            // 获取视频时长
            videoDuration, err := sdk.GetVideoDuration(video)
            if err != nil {
                continue
            }
            totalDuration += videoDuration
            execVideos = append(execVideos, video)
            if totalDuration >= duration {
                break
            }
        }
    }
    // 合并视频
    err = sdk.ConcatenateVideos(
        execVideos,
        basePath+"output1.mp4",
        720, 1280)
    if err != nil {
        t.Error(err)
        return
    }
    clearFile = append(clearFile, basePath+"output1.mp4")
    // 裁剪视频
    err = sdk.CropVideoTimeline(
        basePath+"output1.mp4",
        basePath+"output2.mp4",
        0,
        duration)
    if err != nil {
        t.Error(err)
    }
    clearFile = append(clearFile, basePath+"output2.mp4")
    // 关闭视频原声
    err = sdk.MuteVideo(
        basePath+"output2.mp4",
        basePath+"output3.mp4")
    if err != nil {
        t.Error(err)
    }
    clearFile = append(clearFile, basePath+"output3.mp4")
    // 添加静音轨道
    err = sdk.MuteTrack(
        basePath+"output3.mp4",
        basePath+"output4.mp4")
    if err != nil {
        t.Error(err)
    }
    clearFile = append(clearFile, basePath+"output4.mp4")
    // 添加背景音乐
    err = sdk.AddBackgroundMusic(
        basePath+"output4.mp4",
        basePath+"music.flac",
        basePath+"output5.mp4", 0.4)
    if err != nil {
        t.Error(err)
        return
    }
    clearFile = append(clearFile, basePath+"output5.mp4")
    // 添加解说
    err = sdk.AddBackgroundMusic(
        basePath+"output5.mp4",
        basePath+"mp3.wav",
        basePath+"output6.mp4",
        2)
    clearFile = append(clearFile, basePath+"output6.mp4")
    // 添加图片覆盖
    width, height, err := sdk.GetVideoDimensions(basePath + "output6.mp4")
    if err != nil {
        t.Error(err)
        return
    }
    err = sdk.AddImageOverlay(OverlayOptions{
        ImageWidth:  width,
        ImageHeight: height,
        XPosition:   0,
        YPosition:   0,
        VideoFile:   basePath + "output6.mp4",
        ImageFile:   basePath + "overlay.png",
        OutputFile:  basePath + "output7.mp4",
    })
    if err != nil {
        t.Error(err)
        return
    }
    clearFile = append(clearFile, basePath+"output7.mp4")
    // 字幕要放在视频中间
    // 添加字幕
    err = sdk.AddSubtitles(
        basePath+"output7.mp4",
        "srt.srt",
        basePath+"output8.mp4",
        SubtitleOptions{
            FontSize:  9,
            FontColor: "00FFFFFF",
            XPosition: 0,
            YPosition: 150,
            Alignment: 2,
            Font:      "Arial",
        },
    )
    if err != nil {
        t.Error(err)
        return
    }
    // 清理临时文件
    defer func() {
        for _, file := range clearFile {
            _ = os.Remove(file)
        }
    }()
    t.Log("视频处理完成")
}
