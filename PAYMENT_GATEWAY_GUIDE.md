# 🚀 Xendit Payment Gateway Integration - Simple Guide

Integrasi sederhana Xendit Payment Gateway untuk memahami cara kerja payment gateway.

## 📋 Konsep Payment Gateway

```
User → Backend → Xendit → User Bayar → Webhook → Backend Update Status
```

## 🏗️ Struktur Project

```
app/testing_transaction/
├── domain/testing_transaction.go     # Model & Interface
├── repository/repository.go           # Database operations
├── usecase/usecase.go                 # Business logic
└── delivery/http/handler.go           # HTTP endpoints

pkg/xendit/xendit.go                   # Xendit SDK wrapper

database/
├── migration/..._table_testing_transaction.sql
└── seeder/..._seed_testing_transaction.sql
```

## 🔧 Setup

### 1. Install Dependencies
```bash
go get github.com/xendit/xendit-go/v6
```

### 2. Setup Xendit Account
1. Daftar di [https://dashboard.xendit.co](https://dashboard.xendit.co)
2. Ambil **API Key** (Test Mode) dari menu **Settings → Developers → API Keys**
3. Copy API key yang dimulai dengan `xnd_development_...`

### 3. Update .env
```env
XENDIT_API_KEY="xnd_development_your_actual_key_here"
XENDIT_WEBHOOK_TOKEN="optional_webhook_token"
XENDIT_SUCCESS_REDIRECT="http://localhost:3000/payment/success"
XENDIT_FAILURE_REDIRECT="http://localhost:3000/payment/failed"
```

### 4. Run Migration
```bash
make migration-up
make seed-up
```

### 5. Start Server
```bash
make run-http
# atau
go run main.go http
```

## 🌐 API Endpoints

### 1. Create Payment (Buat Invoice)
**POST** `/api/v1/public/payments`

Request:
```json
{
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "amount": 100000
}
```

Response:
```json
{
  "code": 200,
  "message": "Payment created successfully",
  "data": {
    "order_no": "ORDER-1234567890",
    "invoice_url": "https://checkout.xendit.co/web/xxxxx",
    "amount": 100000,
    "status": "pending"
  }
}
```

### 2. Get Payment by Order No
**GET** `/api/v1/public/payments/{order_no}`

Response:
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "order_no": "ORDER-1234567890",
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "amount": 100000,
    "status": "paid",
    "invoice_url": "https://checkout.xendit.co/web/xxxxx",
    "payment_method": "BANK_TRANSFER",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 3. Get All Payments
**GET** `/api/v1/public/payments`

Response:
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "order_no": "ORDER-1234567890",
      "customer_name": "John Doe",
      "amount": 100000,
      "status": "paid"
    }
  ]
}
```

### 4. Webhook (Otomatis dari Xendit)
**POST** `/api/v1/public/webhooks/xendit`

Xendit akan otomatis hit endpoint ini ketika payment berhasil.

## 🔄 Payment Flow

### 1. User Create Payment
```bash
curl -X POST http://localhost:8000/api/v1/public/payments \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "amount": 50000
  }'
```

Response akan berisi `invoice_url` yang bisa dibuka untuk bayar.

### 2. User Bayar
- Buka `invoice_url` di browser
- Pilih metode pembayaran (Virtual Account BCA/Mandiri/BNI, QRIS, E-wallet)
- Bayar sesuai instruksi

### 3. Xendit Webhook
Setelah user bayar, Xendit otomatis kirim webhook ke backend:
```
POST /api/v1/public/webhooks/xendit
```

Backend akan update status payment dari `pending` → `paid`

## 📝 Code Explanation

### 1. Xendit Client (pkg/xendit/xendit.go)
```go
// Wrapper sederhana untuk Xendit SDK
func (c *Client) CreateInvoice(ctx context.Context, orderNo, email string, amount float64, description string) (invoiceURL string, err error)
```

### 2. Usecase (app/testing_transaction/usecase/usecase.go)

**CreatePayment**:
1. Generate order number
2. Panggil Xendit API untuk buat invoice
3. Simpan ke database dengan status `pending`
4. Return invoice URL ke user

**HandleWebhook**:
1. Terima callback dari Xendit
2. Update status payment di database
3. Bisa ditambahkan: kirim email notifikasi, update inventory, dll

### 3. Handler (app/testing_transaction/delivery/http/handler.go)

HTTP handler untuk menerima request dari client dan webhook dari Xendit.

## 🧪 Testing

### Test dengan Postman/cURL

1. **Create Payment**
```bash
curl -X POST http://localhost:8000/api/v1/public/payments \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"Test User","customer_email":"test@example.com","amount":100000}'
```

2. **Buka invoice_url** yang dikembalikan untuk simulasi pembayaran

3. **Check status payment**
```bash
curl http://localhost:8000/api/v1/public/payments/ORDER-1234567890
```

### Test Webhook di Local

Karena Xendit tidak bisa hit localhost, gunakan salah satu tools:
- **ngrok**: `ngrok http 8000`
- **localtunnel**: `lt --port 8000`

Lalu set webhook URL di Xendit dashboard ke:
```
https://your-ngrok-url.ngrok.io/api/v1/public/webhooks/xendit
```

## 📊 Database Structure

```sql
CREATE TABLE testing_transaction (
    id SERIAL PRIMARY KEY,
    order_no VARCHAR(100) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, paid, expired
    invoice_url TEXT,
    payment_method VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

## 🎯 Key Points

1. **Simple & Clear**: Code dibuat sesederhana mungkin untuk mudah dipahami
2. **Official SDK**: Menggunakan `xendit-go` official SDK, bukan custom HTTP client
3. **Webhook Important**: Webhook adalah cara Xendit memberitahu kita kalau payment berhasil
4. **Async Process**: Payment adalah proses async - user bayar diluar sistem kita
5. **Test Mode**: Gunakan test API key untuk development (tidak ada biaya)

## 🔐 Security Notes

- Jangan commit API key ke git
- Gunakan `.env` untuk menyimpan credentials
- Di production, validate webhook signature dari Xendit
- Gunakan HTTPS untuk webhook endpoint

## 📚 Resources

- [Xendit API Documentation](https://developers.xendit.co/api-reference/)
- [Xendit Go SDK](https://github.com/xendit/xendit-go)
- [Xendit Dashboard](https://dashboard.xendit.co)

## ✅ Next Steps

Setelah memahami basic flow, bisa develop:
1. Integrate dengan event registration yang sudah ada
2. Tambah email notification setelah payment success
3. Implementasi refund/cancellation
4. Add pagination untuk list payments
5. Export payment report

---

**Happy Coding! 🚀**
