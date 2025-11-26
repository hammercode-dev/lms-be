# API Documentation - Event Registration & Payment

Dokumentasi lengkap untuk flow pembelian/registrasi event di Hammercode LMS menggunakan Xendit Payment Gateway.

## Table of Contents
1. [Flow Overview](#flow-overview)
2. [Authentication](#authentication)
3. [API Endpoints](#api-endpoints)
4. [Status Codes](#status-codes)
5. [Examples](#examples)

---

## Flow Overview

Proses registrasi event menggunakan Xendit Payment Gateway:

```
1. User melihat daftar event (GET /api/v1/public/events)
2. User melihat detail event (GET /api/v1/public/events/{id})
3. User membuat transaksi & registrasi event (POST /api/v1/transactions) - Requires Auth
   → Sistem otomatis:
     - Buat registrasi event
     - Generate invoice Xendit
     - Return payment URL
4. User membayar melalui Xendit (eksternal)
5. Xendit kirim webhook ke sistem (POST /api/v1/public/webhooks/xendit)
6. User cek status pembayaran (GET /api/v1/transactions/{transaction_no}/status)
```

---

## Authentication

Untuk endpoint yang memerlukan autentikasi, gunakan Bearer Token di header:

```
Authorization: Bearer <your_jwt_token>
```

Token didapat setelah login melalui endpoint `/api/v1/auth/login`

---

## API Endpoints

### 1. Get List Events (Public)

Mendapatkan daftar semua event yang tersedia.

**Endpoint:** `GET /api/v1/public/events`

**Headers:**
```
Content-Type: application/json
```

**Query Parameters:**
- `page` (optional): Page number, default = 1
- `limit` (optional): Items per page, default = 10
- `type` (optional): Filter by event type (workshop, webinar, conference)
- `status` (optional): Filter by status (open, closed, full)

**Response Success (200):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "Tech Talk: Modern Web Development",
      "description": "Belajar tentang teknologi web terkini",
      "author": "john_doe",
      "image_event": "http://localhost:8000/api/v1/public/storage/images/abc123.png",
      "date_event": "2025-12-01T10:00:00Z",
      "type": "webinar",
      "price": 50000,
      "location": "Online - Zoom",
      "duration": "2 hours",
      "capacity": 100,
      "registration_link": "",
      "session_type": "online",
      "status": "open"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_data": 1,
    "total_page": 1
  }
}
```

---

### 2. Get Event Detail (Public)

Mendapatkan detail lengkap sebuah event.

**Endpoint:** `GET /api/v1/public/events/{id}`

**Headers:**
```
Content-Type: application/json
```

**Path Parameters:**
- `id` (required): Event ID

**Response Success (200):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Tech Talk: Modern Web Development",
    "description": "Belajar tentang teknologi web terkini...",
    "author": "john_doe",
    "image_event": "http://localhost:8000/api/v1/public/storage/images/abc123.png",
    "date_event": "2025-12-01T10:00:00Z",
    "type": "webinar",
    "price": 50000,
    "location": "Online - Zoom",
    "duration": "2 hours",
    "capacity": 100,
    "registration_link": "",
    "session_type": "online",
    "status": "open"
  }
}
```

---

### 3. Create Transaction & Register Event (Protected)

Membuat transaksi pembayaran dan registrasi event sekaligus. Sistem akan otomatis:
1. Membuat registrasi event untuk user
2. Generate invoice Xendit
3. Return payment URL untuk pembayaran

**Endpoint:** `POST /api/v1/transactions`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "event_id": 1
}
```

**Field Descriptions:**
- `event_id` (uint, required): ID event yang ingin diikuti

**Note:** Data user (name, email, phone_number) akan diambil otomatis dari user yang sedang login melalui JWT token.

**Response Success (201):**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "order_no": "ORD-20251107-ABC123",
    "amount": 50000,
    "payment_url": "https://checkout.xendit.co/web/abc123def456",
    "status": "pending"
  }
}
```

**Field Response:**
- `transaction_no`: Nomor transaksi untuk tracking
- `order_no`: Nomor order registrasi event
- `amount`: Total yang harus dibayar
- `payment_url`: URL Xendit untuk melakukan pembayaran
- `status`: Status transaksi (pending)

**Response Error (400 - Event Full):**
```json
{
  "code": 400,
  "message": "Event sudah penuh"
}
```

**Response Error (400 - Already Registered):**
```json
{
  "code": 400,
  "message": "User sudah terdaftar untuk event ini"
}
```

**Response Error (500 - Xendit Error):**
```json
{
  "code": 500,
  "message": "Failed to create payment invoice"
}
```

---

### 4. Check Payment Status (Protected)

Cek status pembayaran berdasarkan transaction number.

**Endpoint:** `GET /api/v1/transactions/{transaction_no}/status`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `transaction_no` (required): Transaction number dari response create transaction

**Response Success (200 - Pending):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "status": "pending",
    "paid_at": null,
    "payment_method": null
  }
}
```

**Response Success (200 - Paid):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "status": "paid",
    "paid_at": "2025-11-07T15:30:00Z",
    "payment_method": "CREDIT_CARD"
  }
}
```

**Status Values:**
- `pending`: Menunggu pembayaran
- `paid`: Sudah dibayar
- `expired`: Invoice expired (biasanya 24 jam)

**Response Error (404):**
```json
{
  "code": 404,
  "message": "Transaction not found"
}
```

---

### 5. Get Registration List (Protected - User)

Melihat daftar registrasi event milik user yang login.

**Endpoint:** `GET /api/v1/admin/events/registrations`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response Success (200):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "order_no": "ORD-20251107-ABC123",
      "event_id": 1,
      "user_id": "user123",
      "image_proof_payment": "http://localhost:8000/api/v1/public/storage/images/payment_proof_123.jpg",
      "payment_date": "2025-11-07T10:00:00Z",
      "status": "pending",
      "created_at": "2025-11-07T09:00:00Z",
      "event_detail": {
        "id": 1,
        "title": "Tech Talk: Modern Web Development",
        "image": "http://localhost:8000/api/v1/public/storage/images/abc123.png",
        "date": "2025-12-01T10:00:00Z",
        "price": 50000
      },
      "user_detail": {
        "id": 123,
        "username": "john_doe",
        "email": "john@example.com"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_data": 1,
    "total_page": 1
  }
}
```

---

### 6. Xendit Webhook (Internal - Called by Xendit)

**Note:** Endpoint ini dipanggil otomatis oleh Xendit ketika status pembayaran berubah. Tidak perlu dipanggil manual oleh user/frontend.

**Endpoint:** `POST /api/v1/public/webhooks/xendit`

**Headers:**
```
Content-Type: application/json
x-callback-token: <xendit_webhook_token>
```

**Request Body (from Xendit):**
```json
{
  "external_id": "TRX-1730987654",
  "status": "PAID",
  "payment_method": "CREDIT_CARD",
  "paid_at": "2025-11-07T15:30:00Z"
}
```

**Response Success (200):**
```json
{
  "code": 200,
  "message": "Webhook processed successfully"
}
```

---

## Status Codes

### Transaction Status
- `pending` - Transaksi dibuat, menunggu pembayaran
- `paid` - Pembayaran berhasil, user terdaftar sebagai peserta
- `expired` - Invoice expired (tidak dibayar dalam waktu yang ditentukan)
- `failed` - Pembayaran gagal

### Event Status
- `open` - Event masih menerima pendaftaran
- `closed` - Event sudah tutup pendaftaran
- `full` - Kapasitas event sudah penuh
- `cancelled` - Event dibatalkan

---

## Examples

### Complete Flow Example

#### Step 1: Login
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user123",
      "email": "john@example.com",
      "username": "john_doe"
    }
  }
}
```

#### Step 2: Get Events List
```bash
curl -X GET http://localhost:8000/api/v1/public/events
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "Tech Talk: Modern Web Development",
      "price": 50000,
      "status": "open",
      "capacity": 100
    }
  ]
}
```

#### Step 3: Get Event Detail
```bash
curl -X GET http://localhost:8000/api/v1/public/events/1
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Tech Talk: Modern Web Development",
    "description": "Learn modern web development...",
    "price": 50000,
    "date_event": "2025-12-01T10:00:00Z",
    "capacity": 100,
    "status": "open"
  }
}
```

#### Step 4: Create Transaction & Register
```bash
curl -X POST http://localhost:8000/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "event_id": 1
  }'
```

Response:
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "order_no": "ORD-20251107-XYZ789",
    "amount": 50000,
    "payment_url": "https://checkout.xendit.co/web/abc123def456",
    "status": "pending"
  }
}
```

**Important:** Simpan `transaction_no` dan redirect user ke `payment_url` untuk melakukan pembayaran!

#### Step 5: User Pays via Xendit
User akan membuka `payment_url` di browser dan melakukan pembayaran melalui berbagai metode:
- Credit Card
- Bank Transfer
- E-Wallet (OVO, Dana, GoPay, dll)
- QRIS
- Retail Outlet (Alfamart, Indomaret)

#### Step 6: Check Payment Status
Setelah user melakukan pembayaran, cek status:

```bash
curl -X GET http://localhost:8000/api/v1/transactions/TRX-1730987654/status \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Response (if still pending):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "status": "pending",
    "paid_at": null,
    "payment_method": null
  }
}
```

Response (if paid):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-1730987654",
    "status": "paid",
    "paid_at": "2025-11-07T15:30:00Z",
    "payment_method": "CREDIT_CARD"
  }
}
```

#### Step 7: Check Registration List
```bash
curl -X GET http://localhost:8000/api/v1/admin/events/registrations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "order_no": "ORD-20251107-XYZ789",
      "status": "paid",
      "event_detail": {
        "id": 1,
        "title": "Tech Talk: Modern Web Development"
      }
    }
  ]
}
```

---

## Notes

1. **Authentication**: Endpoint transaction memerlukan JWT token yang didapat dari login
2. **Payment Gateway**: Menggunakan Xendit untuk payment processing
3. **Transaction Number**: Simpan `transaction_no` untuk cek status pembayaran
4. **Payment URL**: Redirect user ke `payment_url` dari response untuk melakukan pembayaran
5. **Status Flow**: pending → paid (otomatis update via webhook Xendit)
6. **Free Events**: Jika event gratis (`price` = 0), tetap akan generate transaction tapi langsung status `paid`
7. **Webhook**: Xendit akan otomatis memanggil webhook endpoint ketika payment berhasil/gagal
8. **Invoice Expiry**: Invoice Xendit biasanya expire dalam 24 jam jika tidak dibayar

---

## Error Handling

Semua error response mengikuti format:
```json
{
  "code": 400,
  "message": "Error message here"
}
```

Common HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error
