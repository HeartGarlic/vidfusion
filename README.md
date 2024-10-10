# vid-fusion

## 简介
这是一个使用 `ffmpeg` 和 `ffprobe` 实现的 Golang 视频处理 SDK。该 SDK 提供了一系列方便的接口，用于裁剪视频、翻转视频、添加背景音乐、添加图片覆盖、合并视频、以及处理字幕和音频。

## 安装

你可以使用 `go get` 安装该 SDK：

```bash
go get github.com/HeartGarlic/vidfusion
```

## DEMO
```go
package main

import (
    "github.com/HeartGarlic/vidfusion"
    "log"
    "time"
)

func main() {
    // 记录开始时间
    startTime := time.Now().Unix()
    
    wavFile := ".\\" + "mp3.wav"
    srtFile := "srt.srt"
    outputFile := ".\\" + "结婚纪念日那天.mp4"
    
    basePath := ".\\videos\\"
    musicFile := ".\\music\\" + "music.flac"
    overlayFile := ".\\images\\" + "overlay.png"
    
    sdk := vidfusion.NewVideoSDKV2("")
    vidfusion.Encoder = vidfusion.Libx264
    // 获取mp3文件时长
    duration, err := sdk.GetMP3Duration(wavFile)
    if err != nil {
        log.Fatalf("获取mp3文件时长失败: %v", err)
    }
    options := vidfusion.ProcessVideosOptions{
        VideoDuration: duration,
        Width:         720,
        Height:        1280,
        VideosOptions: []vidfusion.VideosOptions{
            {VideoFile: basePath + "1.mp4", Process: vidfusion.FlipVideo},
            {VideoFile: basePath + "2.mp4", Process: vidfusion.SpeedUpVideo, Params: 1.6},
            {VideoFile: basePath + "3.mp4", Process: vidfusion.ScaleUpVideo, Params: 2},
            {VideoFile: basePath + "5.mp4", Process: vidfusion.SpeedUpVideo, Params: 2},
            {VideoFile: basePath + "6.mp4", Process: vidfusion.FlipVideo},
            {VideoFile: basePath + "7.mp4", Process: vidfusion.SpeedUpVideo, Params: 3},
            {VideoFile: basePath + "8.mp4", Process: vidfusion.ScaleUpVideo, Params: 1.5},
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
        log.Fatalf("获取视频尺寸失败: %v", err)
    }
    sdk.AddImageOverlay(vidfusion.OverlayOptions{
        ImageWidth:  width,
        ImageHeight: height,
        XPosition:   0,
        YPosition:   0,
        ImageFile:   overlayFile,
    })
    // 添加字幕 字幕要放在视频中间
    sdk.AddSubtitles(
        srtFile,
        vidfusion.SubtitleOptions{
            FontSize:  9,
            FontColor: "00FFFFFF",
            XPosition: 0,
            YPosition: 150,
            Alignment: 2,
            Font:      "Arial",
        },
    )
    sdk.Finalize(outputFile)
    
    // 记录结束时间
    endTime := time.Now().Unix()
    log.Printf("视频处理完成, 耗时: %d 秒", endTime-startTime)
}

```
