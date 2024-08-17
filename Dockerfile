# Gunakan image dasar golang versi alpine untuk ukuran image yang lebih kecil
FROM golang:1.22.5

# Set working directory di dalam container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy semua file dari direktori lokal ke dalam working directory di container
COPY . .

# Download semua dependency yang diperlukan
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Tentukan command yang dijalankan saat container berjalan
CMD ["./main"]
