# Engineer Competency API

Aplikasi Go untuk layanan autentikasi internal dengan fitur manajemen wallet dan invoice menggunakan framework Gin dan database PostgreSQL.

## Deskripsi

Proyek ini adalah API Gateway untuk layanan autentikasi internal yang menyediakan fitur-fitur seperti:

- OAuth token generation dan verifikasi
- Manajemen wallet (create, top-up, get balance)
- Manajemen invoice (create, get, update status)

## Instalasi

1. Clone repository ini:

   ```bash
   git clone <repository-url>
   cd engineer-competency
   ```

2. Setup environment variables. Copy `.env.example` ke `.env` atau gunakan `.env.docker` untuk konfigurasi Docker, lalu isi nilai-nilai yang diperlukan:
   - `AUTHBASIC_USERNAME`: Username untuk basic auth
   - `AUTHBASIC_PASSWORD`: Password untuk basic auth
   - `DB_PGSQL_PROXY_URL`: URL database PostgreSQL
   - `DB_PGSQL_PROXY_HOST`: Host database
   - `DB_PGSQL_PROXY_PORT`: Port database
   - `DB_PGSQL_PROXY_DWH_USERNAME`: Username database
   - `DB_PGSQL_PROXY_DWH_PASSWORD`: Password database
   - `DB_PGSQL_PROXY_DWH_DATABASE`: Nama database
   - `MAIN_PORT`: Port server (default: 8080)

3. Untuk development lokal, pastikan PostgreSQL berjalan. Untuk Docker, gunakan docker-compose.

## Menjalankan Aplikasi

### Menggunakan Go langsung:

```bash
cd app
go mod tidy
go run main.go
```

### Menggunakan Docker:

```bash
docker-compose up --build
```

Aplikasi akan berjalan di port 8080 (lokal) atau 8100 (Docker).

## Testing

Proyek ini menggunakan unit testing dengan framework `testify` untuk Go.

### Menjalankan Unit Tests

Untuk menjalankan semua unit tests:

```bash
cd app
go test ./controllers/...
```

Untuk menjalankan test dengan coverage:

```bash
go test ./controllers/... -cover
```

Untuk menjalankan test spesifik:

```bash
go test -run TestInvoiceController_CreateInvoice_Success ./controllers/
```

### Struktur Test

- `invoice_controller_test.go`: Test untuk InvoiceController
- `wallet_controller_test.go`: Test untuk WalletController
- `oauth_controller_test.go`: Test untuk OauthController

Test menggunakan mock untuk usecase layer agar fokus pada logika controller tanpa dependency database.

## Swagger API Docs

Aplikasi menggunakan port dari environment variable `MAIN_PORT`. Untuk `.env.docker`, port default adalah `8100`.

Setelah aplikasi berjalan, buka Swagger UI di:

```bash
http://localhost:8100/swagger/index.html
```

Jika Anda menjalankan tanpa Docker dan `MAIN_PORT` diatur lain, ganti `8100` dengan nilai `MAIN_PORT`.

Swagger akan menampilkan endpoint API berikut:

- `POST /api/v1/oauth/token`
- `POST /api/v1/oauth/verify-login-token`
- `GET /api/v1/protected-resource`
- `POST /api/v1/wallet/create`
- `POST /api/v1/wallet/top-up`
- `GET /api/v1/wallet/{id}`
- `GET /api/v1/wallet/user/{user_id}`
- `GET /api/v1/wallet/balance/{user_id}`
- `POST /api/v1/invoice/create`
- `GET /api/v1/invoice/{id}`
- `GET /api/v1/invoice/all`
- `PUT /api/v1/invoice/update-status`

Swagger docs di-generate ke folder `app/docs` oleh tool `swag`, dan file definisi Swagger JSON tersedia sebagai:

- `app/docs/swagger.json`
- `http://localhost:8100/swagger/doc.json`

- Response: `{"status": "OK"}`

### Version

- **GET** `/api/v1/version`
  - Deskripsi: Mendapatkan informasi versi API
  - Response: `{"status": "initialized", "message": "API Gateway Internal Auth Service is running"}`

### OAuth (dengan Basic Auth)

Endpoint ini memerlukan Basic Authentication dengan username dan password dari environment variables.

- **POST** `/api/v1/oauth/token`
  - Deskripsi: Generate OAuth token
  - Body: (sesuai dengan implementasi GenerateToken)

- **POST** `/api/v1/oauth/verify-login-token`
  - Deskripsi: Verifikasi dan validasi login token
  - Body: (sesuai dengan implementasi VerifyAndValidateLoginToken)

### Protected Resources (dengan OAuth Middleware)

Endpoint ini memerlukan OAuth token yang valid.

- **GET** `/api/v1/protected-resource`
  - Deskripsi: Akses resource yang dilindungi
  - Headers: `Authorization: Bearer <token>`
  - Response: `{"status": "success", "message": "You have accessed a protected resource"}`

### Wallet Management

- **POST** `/api/v1/wallet/create`
  - Deskripsi: Membuat wallet baru
  - Headers: `Authorization: Bearer <token>`
  - Body: (sesuai dengan model wallet)

- **POST** `/api/v1/wallet/top-up`
  - Deskripsi: Top up saldo wallet
  - Headers: `Authorization: Bearer <token>`
  - Body: (sesuai dengan model top-up)

- **GET** `/api/v1/wallet/:id`
  - Deskripsi: Mendapatkan wallet berdasarkan ID
  - Headers: `Authorization: Bearer <token>`
  - Params: `id` (wallet ID)

- **GET** `/api/v1/wallet/user/:user_id`
  - Deskripsi: Mendapatkan wallet berdasarkan user ID
  - Headers: `Authorization: Bearer <token>`
  - Params: `user_id`

- **GET** `/api/v1/wallet/balance/:user_id`
  - Deskripsi: Mendapatkan saldo wallet berdasarkan user ID
  - Headers: `Authorization: Bearer <token>`
  - Params: `user_id`

### Invoice Management

- **POST** `/api/v1/invoice/create`
  - Deskripsi: Membuat invoice baru
  - Headers: `Authorization: Bearer <token>`
  - Body: (sesuai dengan model invoice)

- **GET** `/api/v1/invoice/:id`
  - Deskripsi: Mendapatkan invoice berdasarkan ID
  - Headers: `Authorization: Bearer <token>`
  - Params: `id` (invoice ID)

- **GET** `/api/v1/invoice/all`
  - Deskripsi: Mendapatkan semua invoice
  - Headers: `Authorization: Bearer <token>`

- **PUT** `/api/v1/invoice/update-status`
  - Deskripsi: Update status invoice
  - Headers: `Authorization: Bearer <token>`
  - Body: (sesuai dengan model update status)

## Teknologi yang Digunakan

- Go (Gin framework)
- PostgreSQL
- Docker & Docker Compose
