Summary of the key `ffmpeg` and `ffprobe` commands:

---

**1. Create a file list for concatenation**
```
file 'video1.mp4'
file 'video2.mp4'
file 'video3.mp4'
```
(Save as `filelist.txt`)

---

**2. Concatenate MP4 files without re-encoding (fast, if compatible)**
```
ffmpeg -f concat -safe 0 -i filelist.txt -c copy output.mp4
```

---

**3. Concatenate and normalize audio (with re-encoding)**
```
ffmpeg -f concat -safe 0 -i filelist.txt -c:v libx264 -af loudnorm -c:a aac output.mp4
```

---

**4. Concatenate, normalize, and control output quality using CRF**
```
ffmpeg -f concat -safe 0 -i filelist.txt -c:v libx264 -crf 18 -af loudnorm -c:a aac -b:a 192k output.mp4
```
- `-crf 18` for high video quality (lower is better quality)
- `-b:a 192k` for audio bitrate

---

**5. Concatenate, normalize, and match original video/audio bitrate**
```
ffmpeg -f concat -safe 0 -i filelist.txt -c:v libx264 -b:v 20340k -af loudnorm -c:a aac -b:a 192k output.mp4
```
- Replace `20340k` and `192k` with your actual bitrates

---

**6. Extract video bitrate with ffprobe**
```
ffprobe -v error -select_streams v:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 input.mp4
```

---

**7. Extract audio bitrate with ffprobe**
```
ffprobe -v error -select_streams a:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 input.mp4
```

---
