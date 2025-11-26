# Event Purchase Flow Documentation

Dokumentasi lengkap untuk flow pembelian event dari registrasi hingga pembayaran selesai.

---

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Flow Diagram](#flow-diagram)
4. [Step-by-Step Process](#step-by-step-process)
5. [API Endpoints](#api-endpoints)
6. [Payment Status](#payment-status)
7. [Error Handling](#error-handling)

---

## Overview

Flow pembelian event terdiri dari beberapa tahap:
1. User melihat daftar event yang tersedia
2. User memilih event dan membuat transaksi
3. Sistem membuat invoice pembayaran via Xendit
4. User melakukan pembayaran
5. Sistem memverifikasi status pembayaran
6. Registrasi event berhasil

---

## Prerequisites

### Authentication
Semua endpoint (kecuali public endpoints) memerlukan JWT token di header:
```
Authorization: Bearer <your_jwt_token>
```

### User Requirements
- User harus sudah register dan login
- User harus memiliki email yang valid (untuk invoice)

---

## Flow Diagram

```
┌──────────────┐
│ List Events  │
│ (Public)     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Get Event    │
│ Detail       │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│ Login/Register   │
│ (if not logged)  │
└──────┬───────────┘
       │
       ▼
┌──────────────────────┐
│ Create Transaction   │
│ POST /transactions   │
└──────┬───────────────┘
       │
       ├─── Free Event ──────► Success (auto-registered)
       │
       └─── Paid Event ──────┐
                             │
                             ▼
              ┌──────────────────────────┐
              │ Get Payment URL          │
              │ from Response            │
              └──────┬───────────────────┘
                     │
                     ▼
              ┌──────────────────────────┐
              │ User Pays via            │
              │ Xendit Payment URL       │
              └──────┬───────────────────┘
                     │
                     ▼
              ┌──────────────────────────┐
              │ Check Payment Status     │
              │ GET /transactions/       │
              │     {transaction_no}/    │
              │     status               │
              └──────┬───────────────────┘
                     │
                     ├─── Pending ─────► Wait & Retry
                     │
                     ├─── Paid ────────► Success (registered)
                     │
                     ├─── Expired ─────► Failed (re-create transaction)
                     │
                     └─── Failed ──────► Failed (re-create transaction)
```

---

## Step-by-Step Process

### Step 1: Browse Available Events

**Endpoint:** `GET /api/v1/public/events`

**Request:**
```bash
curl -X GET "http://localhost:8000/api/v1/public/events"
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "Workshop React Advanced",
      "description": "Workshop tentang React Advanced",
      "price": 100000,
      "date": "2025-12-01T10:00:00Z",
      "location": "Online",
      "quota": 50,
      "remaining_quota": 45
    }
  ]
}
```

---

### Step 2: Get Event Detail

**Endpoint:** `GET /api/v1/public/events/{id}`

**Request:**
```bash
curl -X GET "http://localhost:8000/api/v1/public/events/1"
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Workshop React Advanced",
    "description": "Workshop lengkap tentang React Advanced patterns",
    "price": 100000,
    "date": "2025-12-01T10:00:00Z",
    "location": "Online via Zoom",
    "quota": 50,
    "remaining_quota": 45,
    "image_url": "https://..."
  }
}
```

---

### Step 3: Login (if not authenticated)

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```bash
curl -X POST "http://localhost:8000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 2,
      "username": "johndoe",
      "email": "user@example.com",
      "role": "user"
    }
  }
}
```

**Save the token** untuk digunakan di request selanjutnya.

---

### Step 4: Create Transaction

**Endpoint:** `POST /api/v1/transactions`

**Request:**
```bash
curl -X POST "http://localhost:8000/api/v1/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "event_id": 1
  }'
```

**Response untuk Paid Event:**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "transaction_no": "TRX-20251108105507-24",
    "order_no": "TXE-1-251108-a3f2",
    "amount": 100000,
    "payment_url": "https://checkout.xendit.co/web/...",
    "status": "pending"
  }
}
```

**Response untuk Free Event:**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "transaction_no": "",
    "order_no": "TXE-1-251108-b4g3",
    "amount": 0,
    "payment_url": "",
    "status": "success"
  }
}
```

**Important Notes:**
- `transaction_no`: ID unik transaksi untuk tracking pembayaran
- `order_no`: ID unik registrasi event
- `payment_url`: URL untuk melakukan pembayaran (only for paid events)
- Jika event gratis (price = 0), status langsung "success" dan tidak perlu payment
- User tidak bisa create transaction 2x untuk event yang sama

---

### Step 5: Pay via Payment URL

**For Paid Events Only**

1. Buka `payment_url` dari response Step 4
2. Pilih metode pembayaran (Virtual Account, E-Wallet, Credit Card, dll)
3. Ikuti instruksi pembayaran dari Xendit
4. Selesaikan pembayaran

**Payment Methods Available:**
- Virtual Account (BCA, BNI, Mandiri, BRI, Permata)
- E-Wallet (OVO, DANA, LinkAja, ShopeePay)
- Credit Card
- Retail Outlet (Alfamart, Indomaret)
- QRIS

---

### Step 6: Check Payment Status

**Endpoint:** `GET /api/v1/transactions/{transaction_no}/status`

**Request:**
```bash
curl -X GET "http://localhost:8000/api/v1/transactions/TRX-20251108105507-24/status" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (Pending):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-20251108105507-24",
    "status": "pending",
    "paid_at": null,
    "payment_method": null
  }
}
```

**Response (Paid/Success):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "transaction_no": "TRX-20251108105507-24",
    "status": "paid",
    "paid_at": "2025-11-08T10:55:30Z",
    "payment_method": "BCA Virtual Account"
  }
}
```

**When to Check:**
- Poll setiap 5-10 detik setelah user melakukan pembayaran
- Atau gunakan webhook untuk auto-update (jika sudah disetup)

---

### Step 7: View My Registrations

**Endpoint:** `GET /api/v1/events/registrations`

**Request:**
```bash
curl -X GET "http://localhost:8000/api/v1/events/registrations" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "order_no": "TXE-1-251108-a3f2",
      "event": {
        "id": 1,
        "title": "Workshop React Advanced",
        "date": "2025-12-01T10:00:00Z"
      },
      "status": "success",
      "payment_date": "2025-11-08T10:55:30Z"
    }
  ]
}
```

---

## API Endpoints

### Public Endpoints (No Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/public/events` | List all events |
| GET | `/api/v1/public/events/{id}` | Get event detail |
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |

### Protected Endpoints (Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/transactions` | Create transaction for event |
| GET | `/api/v1/transactions/{transaction_no}/status` | Check payment status |
| GET | `/api/v1/events/registrations` | List user's registrations |

---

## Payment Status

### Transaction Status

| Status | Description | Next Action |
|--------|-------------|-------------|
| `pending` | Menunggu pembayaran | User perlu bayar via payment_url |
| `paid` | Pembayaran berhasil | Registrasi event berhasil |
| `expired` | Invoice expired | Create transaction baru |
| `failed` | Pembayaran gagal | Create transaction baru |

### Registration Status

| Status | Description |
|--------|-------------|
| `pending` | Menunggu pembayaran |
| `success` | Registrasi berhasil (paid atau free event) |
| `cancelled` | Pembayaran gagal |
| `expired` | Invoice expired |

### Status Mapping

```
Transaction Status → Registration Status
─────────────────────────────────────────
pending           → pending
paid              → success
expired           → expired
failed            → cancelled
```

---

## Error Handling

### Common Errors

#### 1. Event Not Found
```json
{
  "code": 404,
  "message": "Event not found",
  "data": null
}
```

#### 2. Already Registered
```json
{
  "code": 400,
  "message": "You have already registered for this event",
  "data": null
}
```

#### 3. Unauthorized
```json
{
  "code": 401,
  "message": "Unauthorized",
  "data": null
}
```

#### 4. Transaction Not Found
```json
{
  "code": 404,
  "message": "Transaction not found",
  "data": null
}
```

#### 5. Event Full (Quota Habis)
```json
{
  "code": 400,
  "message": "Event quota is full",
  "data": null
}
```

---

## Complete Flow Example

### Scenario: User membeli event berbayar

```bash
# 1. Get available events
curl -X GET "http://localhost:8000/api/v1/public/events"

# 2. Get event detail
curl -X GET "http://localhost:8000/api/v1/public/events/1"

# 3. Login
curl -X POST "http://localhost:8000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "putrasatria893@gmail.com",
    "password": "password123"
  }'

# Save token from response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 4. Create transaction
curl -X POST "http://localhost:8000/api/v1/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "event_id": 1
  }'

# Response will contain:
# - transaction_no: "TRX-20251108105507-24"
# - payment_url: "https://checkout.xendit.co/web/..."

# 5. User opens payment_url and pays

# 6. Check payment status (poll every 5-10 seconds)
curl -X GET "http://localhost:8000/api/v1/transactions/TRX-20251108105507-24/status" \
  -H "Authorization: Bearer $TOKEN"

# 7. When status becomes "paid", registration is complete!

# 8. View my registrations
curl -X GET "http://localhost:8000/api/v1/events/registrations" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Testing

### Test dengan Free Event

1. Create event dengan price = 0
2. Create transaction
3. Response akan langsung status "success"
4. Tidak perlu payment process

### Test dengan Paid Event

1. Create event dengan price > 0
2. Create transaction
3. Buka payment_url
4. Pay via Xendit test mode
5. Check status sampai "paid"

### Xendit Test Mode

Untuk testing, gunakan Xendit test credentials:
- API Key: `xnd_development_...`
- Test payments akan langsung berhasil tanpa real money

**Virtual Account Test:**
```
Bank: BCA
VA Number: Provided by Xendit
Amount: Sesuai invoice
```

Simulate payment via Xendit Dashboard → Simulations

---

## Notes

### Invoice Expiration
- Default: 24 jam setelah dibuat
- Setelah expired, user harus create transaction baru

### Double Registration Prevention
- System check apakah user sudah punya registrasi dengan status `pending` atau `success`
- Jika sudah ada, akan return error "already registered"

### Payment Method
- Tersimpan setelah payment berhasil
- Bisa dilihat di response check payment status

### Webhook (Optional)
- Xendit bisa mengirim webhook saat payment berhasil
- Auto-update status tanpa perlu polling
- Perlu setup webhook endpoint & verify signature

---

## Troubleshooting

### Empty Reply from Server

**Problem:** API return empty response

**Possible Causes:**
1. Server crash/panic
2. Database connection error
3. Middleware error

**Solution:**
- Check server logs
- Verify database connection
- Check if user exists in context

### Payment Status Not Updating

**Problem:** Status masih pending padahal sudah bayar

**Possible Causes:**
1. Xendit webhook tidak jalan
2. Invoice ID tidak valid

**Solution:**
- Manual check via endpoint status
- Verify invoice ID di database
- Check Xendit dashboard

---

## FAQ

**Q: Berapa lama invoice berlaku?**
A: Default 24 jam. Setelah itu expired dan perlu create transaction baru.

**Q: Apakah bisa bayar dengan metode lain setelah pilih Virtual Account?**
A: Ya, invoice URL mendukung multiple payment methods.

**Q: Bagaimana jika payment gagal?**
A: Create transaction baru untuk mendapatkan invoice baru.

**Q: Apakah bisa cancel transaksi?**
A: Tidak ada cancel endpoint. Biarkan invoice expired automatically.

**Q: Bagaimana cara refund?**
A: Harus manual via Xendit dashboard atau contact admin.

---

## Related Documentation

- [Payment Gateway Guide](./PAYMENT_GATEWAY_GUIDE.md)
- [Event Registration API](./API_EVENT_REGISTRATION.md)
- [Xendit Documentation](https://docs.xendit.co)

---

**Last Updated:** 2025-11-22
**Version:** 1.0.0
