# Ngrok Setup Guide untuk Xendit Webhook Testing

## 1. Install Ngrok

### macOS (Homebrew)
```bash
brew install ngrok/ngrok/ngrok
```

### macOS (Manual)
```bash
# Download
curl -O https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-darwin-amd64.zip

# Extract
unzip ngrok-v3-stable-darwin-amd64.zip

# Move to PATH
sudo mv ngrok /usr/local/bin/
```

### Linux
```bash
# Download
wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz

# Extract
tar xvzf ngrok-v3-stable-linux-amd64.tgz

# Move to PATH
sudo mv ngrok /usr/local/bin/
```

### Windows
1. Download dari https://ngrok.com/download
2. Extract ZIP file
3. Jalankan `ngrok.exe` dari folder

## 2. Setup Ngrok Account (Optional tapi Recommended)

### Kenapa perlu account?
- Free tier sudah cukup untuk development
- Dapat custom subdomain (tidak berubah-ubah)
- Session tidak expire
- Lebih stable

### Cara daftar:
1. Buka https://dashboard.ngrok.com/signup
2. Sign up dengan email/GitHub/Google
3. Setelah login, copy authtoken dari dashboard
4. Add authtoken:

```bash
ngrok config add-authtoken YOUR_AUTH_TOKEN_HERE
```

## 3. Jalankan Ngrok

### Step 1: Jalankan Aplikasi Anda
```bash
# Di terminal pertama
go run main.go http

# Atau jika sudah build
./lms-be http
```

Pastikan aplikasi berjalan di port yang benar (misal: `:8000`)

### Step 2: Jalankan Ngrok
```bash
# Di terminal kedua
ngrok http 8000
```

### Output Ngrok:
```
ngrok

Session Status                online
Account                       your-email@example.com (Plan: Free)
Version                       3.5.0
Region                        Asia Pacific (ap)
Latency                       45ms
Web Interface                 http://127.0.0.1:4040
Forwarding                    https://abc123.ngrok-free.app -> http://localhost:8000

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

### Step 3: Copy URL Publik
URL publik Anda: `https://abc123.ngrok-free.app`

Webhook endpoint: `https://abc123.ngrok-free.app/webhooks/xendit`

## 4. Setup di Xendit Dashboard

### Development (Sandbox)
1. Login ke https://dashboard.xendit.co
2. Switch ke **Test Mode** (toggle di pojok kanan atas)
3. Settings → **Webhooks** (di sidebar kiri)
4. Klik **Add Webhook URL**
5. Input:
   - URL: `https://abc123.ngrok-free.app/webhooks/xendit`
   - Environment: **Test**
6. Select Events:
   - ✅ `invoice.paid`
   - ✅ `invoice.expired`
   - ✅ `fva.payment` (jika pakai Fixed VA)
7. Klik **Save**

### Production (Live)
1. Switch ke **Live Mode**
2. Ulangi langkah yang sama dengan URL production Anda

## 5. Test Webhook

### Cara 1: Test dengan Xendit Webhook Tester
1. Di Xendit Dashboard → Webhooks
2. Klik **Test Webhook** di sebelah webhook URL
3. Pilih event type: `invoice.paid`
4. Klik **Send Test**

### Cara 2: Buat Transaksi Real
```bash
# 1. Login
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 2. Buat transaksi (gunakan token dari response login)
curl -X POST http://localhost:8000/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "event_id": 1
  }'

# 3. Buka payment_url dari response
# 4. Lakukan pembayaran test di sandbox Xendit
```

### Cara 3: Manual Test via Ngrok URL
```bash
curl -X POST https://abc123.ngrok-free.app/webhooks/xendit \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-invoice-123",
    "external_id": "TRX-20250126120000-1",
    "status": "PAID",
    "amount": 80000,
    "paid_at": "2025-01-26T12:00:00.000Z",
    "payment_channel": "BANK_TRANSFER"
  }'
```

## 6. Monitor Webhook Requests

### Ngrok Web Interface
```bash
# Buka di browser
http://localhost:4040
```

Di dashboard ini Anda bisa:
- ✅ Lihat semua HTTP requests yang masuk
- ✅ Inspect request headers & body
- ✅ Replay requests untuk debugging
- ✅ Lihat response dari aplikasi

### Application Logs
Cek logs di terminal aplikasi Anda:
```
[INFO] Received Xendit webhook: {"id":"test-invoice-123",...}
[INFO] Webhook processed successfully for transaction TRX-20250126120000-1
```

## 7. Tips & Best Practices

### Keep URL Consistent (Paid Feature)
```bash
# Free tier: URL berubah setiap restart
ngrok http 8000

# Paid tier: Custom subdomain (tetap sama)
ngrok http 8000 --subdomain=my-lms-webhook
# URL: https://my-lms-webhook.ngrok.io
```

### Run in Background
```bash
# Jalankan ngrok di background
ngrok http 8000 > /dev/null &

# Atau gunakan screen/tmux
screen -S ngrok
ngrok http 8000
# Ctrl+A, D untuk detach
```

### Custom Config File
Buat file `ngrok.yml`:
```yaml
version: "2"
authtoken: YOUR_AUTH_TOKEN
tunnels:
  lms-webhook:
    proto: http
    addr: 8000
    bind_tls: true
```

Jalankan:
```bash
ngrok start lms-webhook
```

## 8. Troubleshooting

### Error: "ngrok: command not found"
```bash
# Pastikan ngrok sudah di PATH
which ngrok

# Jika tidak ada, install ulang atau add ke PATH
export PATH=$PATH:/usr/local/bin
```

### Error: "Failed to listen on localhost:4040"
Ngrok web interface sudah jalan. Kill process:
```bash
pkill ngrok
ngrok http 8000
```

### Webhook tidak sampai ke aplikasi
1. Cek ngrok masih jalan: `http://localhost:4040`
2. Cek aplikasi masih jalan: `curl http://localhost:8000/health`
3. Cek URL di Xendit dashboard benar
4. Cek logs aplikasi untuk error

### URL berubah setiap restart
- Normal untuk free tier
- Upgrade ke paid ($8/month) untuk custom subdomain
- Atau update URL di Xendit setiap restart

## 9. Alternative: Cloudflare Tunnel (Free & Permanent)

Jika tidak mau install ngrok, bisa pakai Cloudflare Tunnel:

```bash
# Install
brew install cloudflared

# Jalankan
cloudflared tunnel --url http://localhost:8000
```

URL yang dihasilkan permanent dan tidak berubah.

## 10. Clean Up

Setelah selesai testing:

```bash
# Stop ngrok
# Tekan Ctrl+C di terminal ngrok

# Atau kill process
pkill ngrok

# Remove webhook di Xendit dashboard (optional)
# Settings → Webhooks → Delete
```

## Ready to Test! 🚀

1. ✅ Jalankan aplikasi: `go run main.go http`
2. ✅ Jalankan ngrok: `ngrok http 8000`
3. ✅ Copy URL ngrok
4. ✅ Setup di Xendit dashboard
5. ✅ Test dengan transaksi atau webhook tester
6. ✅ Monitor di http://localhost:4040

Selamat mencoba!
