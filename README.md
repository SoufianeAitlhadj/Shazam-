# Shazoom - a shazam clone 
## Audio Fingerprinting & Recognition System

## 📌 Overview

This project is a simplified implementation of a Shazam-like audio recognition system built in Go.

The system:

- Extracts spectral fingerprints from `.wav` audio files using FFT  
- Stores fingerprints in PostgreSQL  
- Uses a simplified offset-voting algorithm for matching  
- Returns the **top 3 ranked matching songs**  
- Displays:
  - Title  
  - Artist  
  - Album  
  - Embedded YouTube video  
  - Processing time (in milliseconds)  

This implementation is a **simplified tweak of the real Shazam algorithm**, as it extracts **only one dominant frequency bin per window** instead of multiple peak pairs.

---
## 🎬 Demo

Click here: https://youtu.be/B8SyA4M7y04

## 🧠 How It Works

### 1️⃣ Audio Processing Pipeline

1. Read `.wav` file  
2. Split signal into overlapping windows  
3. Apply FFT to each window  
4. Extract the dominant frequency bin per window  
5. Store fingerprints as: ## 🧠 How It Works

### 1️⃣ Audio Processing Pipeline

1. Read `.wav` file  
2. Split signal into overlapping windows  
3. Apply FFT to each window  
4. Extract the dominant frequency bin per window  
5. Store fingerprints as: (song_id, hash, time_offset)


---

### 2️⃣ Matching Algorithm (Simplified Shazam)

Instead of raw overlap counting, this implementation uses:

- Offset voting (temporal alignment)
- For each matching hash:  
---

### 2️⃣ Matching Algorithm (Simplified Shazam)

Instead of raw overlap counting, this implementation uses:

- Offset voting (temporal alignment)
- For each matching hash: delta = db_time_offset - query_time_offset


- Votes are counted per `(song_id, delta)`
- The most consistent temporal alignment wins

The system returns:

- Top 3 ranked songs
- Ranked by highest temporal vote score

---

## ⚙ Requirements

### 🖥 Software

- Go 1.20+
- Docker
- PostgreSQL (runs inside Docker)

### 📦 Go Dependencies

Installed automatically via:

```bash
go mod tidy
``` 

- Votes are counted per `(song_id, delta)`
- The most consistent temporal alignment wins

The system returns:

- Top 3 ranked songs
- Ranked by highest temporal vote score

---

## ⚙ Requirements

### 🖥 Software

- Go 1.20+
- Docker
- PostgreSQL (runs inside Docker)

### 📦 Go Dependencies

Installed automatically via:

```bash
go mod tidy
```

Main external libraries used:

- github.com/go-audio/wav
- gonum.org/v1/gonum/dsp/fourier
- github.com/lib/pq

## Getting started:
1) Start PostgresSQL (Docker) 
```bash
docker run --name shazam-postgres-container \
-e POSTGRES_USER=postgres \
-e POSTGRES_DB=songs_db \
-e POSTGRES_HOST_AUTH_METHOD=trust \
-p 5432:5432 \
-d postgres:16
```

2) Build offline database (among which you can search)
Just uncomment the file *cmd/server/offline_database_build.go* and run it.

This:
- Reads .wav files from the /songs folder
- Extracts fingerprints
- Populates the PostgreSQL database
**An initial dataset of songs is provided for MVP demonstration.**

3) Start the web server with : 
```bash
go run cmd/server/main.go"
```

and open: 
```bash
http://localhost:8080
```
and upload a .wav file to identify it.

## Output :
The system returns:
    - top 3 ranked matches
    - their scores (vote count)
    - processing time (in ms) in the end of the page
    - embedded youtube video

## 🧩 Simplifications Compared to Real Shazam

| Real Shazam | This Implementation |
|-------------|--------------------|
| Multiple spectral peaks | Single dominant frequency |
| Peak pairing | No peak pairing |
| Complex hash structure | Single-bin hash |
| Industrial-scale search | Local PostgreSQL |

This was done intentionally to:

- Keep the system understandable  
- Allow fast implementation  
- Demonstrate core fingerprinting logic  

## 📚 Resources

- **Shazam Algorithm Explanation (Toptal Article)**  
  https://www.toptal.com/developers/algorithms/shazam-it-music-processing-fingerprinting-and-recognition  

- **Original Shazam Paper**  
  *An Industrial-Strength Audio Search Algorithm*  
  Avery Li-Chun Wang  

- **Reference GitHub Repository**  
  https://github.com/Danztee/shazam-build  

- **Helpful YouTube Video**  
  https://www.youtube.com/watch?v=a0CVCcb0RJM  


by *Soufiane AIT LHADJ*
