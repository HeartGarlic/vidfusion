package vidfusion

import (
    "testing"
)

// baseDir 是测试视频文件的基础目录
var baseDir = "C:\\Users\\Administrator\\Desktop\\视频\\"

// TestNewVideoSDKV2 测试 NewVideoSDKV2 函数
func TestNewVideoSDKV2(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    if sdk.CurrentFile != baseDir+"1.mp4" {
        t.Errorf("CurrentFile should be test.mp4, but got %s", sdk.CurrentFile)
    }
    if len(sdk.tempFiles) != 1 {
        t.Errorf("tempFiles should have 1 element, but got %d", len(sdk.tempFiles))
    }
    t.Logf("tempFiles: %v", sdk.tempFiles)
}

// TestGetNextTempFile 测试 getNextTempFile 函数
func TestGetNextTempFile(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    tempFile := sdk.getNextTempFile()
    if len(sdk.tempFiles) != 2 {
        t.Errorf("tempFiles should have 2 elements, but got %d", len(sdk.tempFiles))
    }
    t.Logf("tempFile: %s", tempFile)
    t.Logf("tempFiles: %v", sdk.tempFiles)
}

// TestCleanup 测试 Cleanup 函数
func TestCleanup(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    tempFile := sdk.getNextTempFile()
    t.Logf("tempFile: %s", tempFile)
    sdk.Cleanup()
    if len(sdk.tempFiles) != 1 {
        t.Errorf("tempFiles should have 1 element, but got %d", len(sdk.tempFiles))
    }
    t.Logf("tempFiles: %v", sdk.tempFiles)
}

// TestCropVideoTimeline 测试 CropVideoTimeline 函数
func TestCropVideoTimeline(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.CropVideoTimeline(0, 5)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// TestCropVideo 测试 CropVideo 函数
func TestCropVideo(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.CropVideo(640, 360)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// FlipVideo 翻转视频
func TestFlipVideo(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.FlipVideo()
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// SpeedUpVideo 加速视频 @TODO 视频整体时长没有变?
func TestSpeedUpVideo(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.SpeedUpVideo(2)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// ScaleUpVideo 放大视频
func TestScaleUpVideo(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.ScaleUpVideo(2)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// AddBackgroundMusic 添加背景音乐
func TestAddBackgroundMusic(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.AddBackgroundMusic(baseDir+"music.flac", 0.5)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// AddImageOverlay 添加图片覆盖
func TestAddImageOverlay(t *testing.T) {
    // 获取视频的宽度和高度
    width, height, _ := NewVideoSDKV2("").GetVideoDimensions(baseDir + "1.mp4")
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.AddImageOverlay(OverlayOptions{
        ImageWidth:  width,
        ImageHeight: height,
        XPosition:   0,
        YPosition:   0,
        ImageFile:   baseDir + "overlay.png",
    })
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// Mute 去除声音
func TestMute(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.Mute()
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// MuteTrack
func TestMuteTrack(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.MuteTrack()
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// AddSubtitles 添加字幕
func TestAddSubtitles(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.AddSubtitles("srt.srt", SubtitleOptions{
        FontSize:  9,
        FontColor: "00FFFFFF",
        XPosition: 0,
        YPosition: 150,
        Alignment: 2,
        Font:      "Arial",
    })
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// ConcatenateVideos 合并多个视频
func TestConcatenateVideos(t *testing.T) {
    sdk := NewVideoSDKV2(baseDir + "1.mp4")
    sdk.ConcatenateVideos([]string{
        baseDir + "1.mp4",
        baseDir + "2.mp4",
    }, 720, 1280)
    sdk.Finalize(baseDir + "output.mp4")
    t.Logf("CurrentFile: %s", sdk.CurrentFile)
}

// Demo
func TestDemo(t *testing.T) {
    wavFile := ".\\" + "结婚纪念日那天.wav"
    srtFile := "结婚纪念日那天.srt"
    outputFile := ".\\" + "结婚纪念日那天.mp4"
    
    basePath := ".\\videos\\"
    musicFile := ".\\music\\" + "music.flac"
    overlayFile := ".\\images\\" + "overlay.png"
    
    sdk := NewVideoSDKV2("")
    // 获取mp3文件时长
    duration, err := sdk.GetMP3Duration(wavFile)
    if err != nil {
        t.Fatal(err)
    }
    options := ProcessVideosOptions{
        VideoDuration: duration,
        Width:         720,
        Height:        1280,
        VideosOptions: []VideosOptions{
            {VideoFile: basePath + "1.mp4", Process: FlipVideo},
            {VideoFile: basePath + "2.mp4", Process: SpeedUpVideo, Params: 1.6},
            {VideoFile: basePath + "3.mp4", Process: ScaleUpVideo, Params: 2},
            {VideoFile: basePath + "5.mp4", Process: SpeedUpVideo, Params: 2},
            {VideoFile: basePath + "6.mp4", Process: FlipVideo},
            {VideoFile: basePath + "7.mp4", Process: SpeedUpVideo, Params: 3},
            {VideoFile: basePath + "8.mp4", Process: ScaleUpVideo, Params: 1.5},
        },
    }
    sdk.ProcessVideos(options)
    // 裁剪视频时长
    sdk.CropVideoTimeline(0, duration)
    // 关闭视频原声
    sdk.Mute()
    // 添加静音轨道
    sdk.MuteTrack()
    // 添加背景音乐
    sdk.AddBackgroundMusic(musicFile, 0.4)
    // 添加解说
    sdk.AddBackgroundMusic(wavFile, 2)
    // 添加图片覆盖
    width, height, err := sdk.GetVideoDimensions(sdk.CurrentFile)
    if err != nil {
        t.Fatal(err)
    }
    sdk.AddImageOverlay(OverlayOptions{
        ImageWidth:  width,
        ImageHeight: height,
        XPosition:   0,
        YPosition:   0,
        ImageFile:   overlayFile,
    })
    // 添加字幕 字幕要放在视频中间
    sdk.AddSubtitles(
        srtFile,
        SubtitleOptions{
            FontSize:  9,
            FontColor: "00FFFFFF",
            XPosition: 0,
            YPosition: 150,
            Alignment: 2,
            Font:      "Arial",
        },
    )
    sdk.Finalize(outputFile)
    t.Log("视频处理完成")
}
