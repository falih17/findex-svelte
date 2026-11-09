# Findex - Enterprise Multi-Tenant Attendance SaaS

Aplikasi SaaS B2B Presensi Pegawai & HR Management dengan verifikasi **Biometrik Wajah (Anti-Spoofing & 512-D Embedding)**, **Geofencing Lokasi Kerja (PostGIS)**, dan **Isolasi Data Multi-Tenant (PostgreSQL Row Level Security)**.

---

## 🏛️ Struktur Monorepo Proyek Findex

Arsitektur sistem dibangun dengan pola **Multi-Service Monorepo**:

```
findex-svelte/
├── docker/                              # Infrastruktur Database & Layanan Pendukung
│   ├── postgres/
│   │   ├── Dockerfile                   # PostgreSQL 16 + PostGIS 3.4 + pgvector 0.8.0
│   │   └── init/
│   │       ├── 01-extensions.sql        # Inisialisasi ekstensi PostGIS & pgvector
│   │       └── 02-schema.sql            # DDL, Partisi Tabel, RLS, & Seed Data Demo
│   └── minio/
│       └── init.sh                      # Inisialisasi bucket storage S3 foto presensi
│
├── services/
│   ├── backend-core/                    # [Core API Service] - Golang 1.23 Clean Architecture
│   │   ├── cmd/api/main.go              # Entrypoint server & Dependency Injection
│   │   ├── internal/
│   │   │   ├── domain/                  # Entities & Interface kontrak repository/usecase
│   │   │   │   ├── company.go           # Tenant & Company Settings
│   │   │   │   ├── branch.go            # Lokasi Kantor & Geofence (Radius)
│   │   │   │   ├── employee.go          # Data Pegawai & Master Biometrik
│   │   │   │   ├── attendance.go        # Transaksi Presensi & Status
│   │   │   │   └── shift.go             # Shift Kerja & Pengajuan Cuti
│   │   │   ├── repository/              # PostgreSQL + PostGIS + pgvector implementations
│   │   │   │   ├── branch_repository.go # Query spasial: ST_DWithin & ST_Distance
│   │   │   │   ├── employee_repository.go # Query biometrik: 1 - (face_embedding <=> $1)
│   │   │   │   └── attendance_repository.go # Insert ke tabel partisi bulanan
│   │   │   ├── usecase/                 # Business logic presensi, auth, dan geofence
│   │   │   │   ├── attendance_usecase.go
│   │   │   │   ├── auth_usecase.go
│   │   │   │   └── branch_usecase.go
│   │   │   ├── delivery/http/           # Fiber router & HTTP handler
│   │   │   ├── middleware/              # JWT Auth & Tenant RLS context injector
│   │   │   └── pkg/                     # Driver database, standard response & AI client
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── ai-face-service/                 # [AI Biometrics Worker] - Python 3.11 FastAPI
│   │   ├── app/
│   │   │   ├── core/config.py           # Threshold anti-spoofing & model parameters
│   │   │   ├── models/
│   │   │   │   ├── detector.py          # Deteksi wajah & cropping foreground
│   │   │   │   ├── liveness.py          # Anti-spoofing heuristic (Laplacian & Saturation)
│   │   │   │   └── extractor.py         # Ekstraksi vector embedding 512 dimensi
│   │   │   ├── schemas/face.py          # Pydantic request & response schemas
│   │   │   ├── api/v1/endpoints/face.py # Endpoint /extract-face & /compare-vectors
│   │   │   └── main.py                  # Entrypoint FastAPI
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   │
│   └── frontend/                        # [UI Client & PWA] - SvelteKit & TailwindCSS
│       ├── src/
│       │   ├── lib/
│       │   │   ├── stores/auth.svelte.ts # Svelte 5 Runes Auth State Management
│       │   │   ├── components/Navbar.svelte # Top navigation & portal switch
│       │   │   └── api/client.ts        # Typed HTTP API client
│       │   └── routes/
│       │       ├── +page.svelte         # Clock-in Presensi (WebCam + GPS Live)
│       │       ├── (auth)/login/        # Login Form (Pre-filled demo credentials)
│       │       ├── (employee)/
│       │       │   ├── clock-in/        # Halaman presensi mobile PWA
│       │       │   └── history/         # Riwayat presensi pegawai
│       │       └── (admin)/
│       │           ├── dashboard/       # HRD Executive Monitoring & Live Log
│       │           └── branches/        # Pengaturan Lokasi Absen & Slider Radius
│       ├── static/manifest.json         # Manifest PWA (Add to Home Screen)
│       ├── Dockerfile
│       └── package.json
│
├── docker-compose.yml                   # Orkestrator multi-kontainer Docker
├── .env.example                         # Contoh variabel lingkungan
└── README.md
```

---

## ⚡ Fitur Utama Findex

1. **Biometric Face Verification & Anti-Spoofing**:
   * Menghitung skor ketajaman tepi (*Laplacian variance*) dan konsistensi warna untuk mendeteksi apakah wajah yang menghadap kamera adalah orang asli atau foto hasil cetak/layar ponsel (*replay attack*).
   * Menghasilkan representasi wajah berupa vektor 512 dimensi yang ternormalisasi (*L2-normalized*).

2. **Perhitungan Geofencing Presisi Tinggi (PostGIS)**:
   * Setiap perusahaan dapat mendaftarkan kantor pusat, cabang, atau proyek lapangan.
   * Toleransi radius presensi dapat diatur per cabang (misal: 50m, 100m, 200m).
   * Validasi jarak menggunakan fungsi bawaan PostGIS:
     ```sql
     ST_DWithin(location_geom, ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography, radius_meters)
     ```

3. **Pencocokan Biometrik Instan di PostgreSQL (pgvector)**:
   * Menghitung kemiripan kosinus (*Cosine Similarity*) langsung di level database:
     ```sql
     SELECT 1 - (face_embedding <=> $1::vector) AS similarity 
     FROM employee_biometrics WHERE employee_id = $2;
     ```

4. **Multi-Tenant Data Isolation (PostgreSQL RLS)**:
   * Setiap tabel transaksi memiliki `company_id`.
   * Row-Level Security (RLS) menjamin data antar perusahaan tidak saling tumpang tindih.

5. **HR Portal Dashboard & PWA Mobile Client**:
   * Karyawan dapat membuka aplikasi di browser ponsel atau memasangnya sebagai aplikasi mandiri (PWA).
   * HRD memiliki dasbor pemantauan kehadiran, anomali jarak, dan pengaturan lokasi kerja.

---

## 🚀 Menjalankan Aplikasi dengan Docker

```bash
# 1. Salin konfigurasi environment
cp .env.example .env

# 2. Build & jalankan seluruh kontainer
docker compose up -d --build

# 3. Pantau log
docker compose logs -f
```

### Akses Port:
* **Frontend PWA & HR Portal**: [http://localhost:3000](http://localhost:3000)
* **Backend Core API (Golang)**: [http://localhost:8080/health](http://localhost:8080/health)
* **AI Face Biometric Service (Python)**: [http://localhost:8000/health](http://localhost:8000/health)
* **PostgreSQL (PostGIS + pgvector)**: `localhost:5432` (`findex_user` / `findex_secret_pwd`)
* **MinIO Object Storage Console**: [http://localhost:9001](http://localhost:9001) (`minioadmin` / `minioadminsecret`)

---

## 📱 Build APK Android dengan Capacitor

Aplikasi frontend Findex telah dikonfigurasi penuh dengan **Capacitor 6** agar dapat dikompilasi menjadi file APK mandiri (.apk) yang dapat diunduh dan dipasang oleh karyawan perusahaan:

### Langkah Build APK:
```bash
# 1. Jalankan script otomatisasi
./scripts/build-apk.sh debug

# Atau build release untuk produksi:
./scripts/build-apk.sh release
```

File APK akan dihasilkan di:
* `dist/findex-attendance-debug.apk`

### Halaman Unduh untuk Karyawan:
Karyawan dapat langsung membuka **[http://localhost:3000/download](http://localhost:3000/download)** di browser ponsel mereka untuk mengunduh APK secara langsung tanpa harus melalui Google Play Store.
