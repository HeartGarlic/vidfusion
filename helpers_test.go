package vidfusion

import "testing"

func Test_getRandomFiles(t *testing.T) {
    files, err := getRandomFiles("C:\\Users\\Administrator\\Desktop\\视频", "*.mp4", 3)
    if err != nil {
        t.Errorf("getRandomFiles error: %v", err)
    }
    t.Logf("random files: %v", files)
}
