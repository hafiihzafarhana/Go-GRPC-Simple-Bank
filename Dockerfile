# Agar ukuran image itu kecil, maka ambil binary file saja
# dengan cara menerapkan konsep multistage

# Stage: builder
# Untuk menentukan basic image
FROM golang:1.21-alpine3.18 AS builder
# Mendeklarasikan direktori kerja saat ini di dalam gambar
# Setelah dijalankan, Docker akan membuat direktori /app, sebagai tempat menjalankan perintah 
WORKDIR /app
# Copy semua yang ada pada direktori ini (pada titik pertama)
# Pada titik kedua adalah direktori kerja saat ini di dalam image yaitu /app
COPY . .
# Membangun aplikasi ke satu file biner executable
# output bernama main
# main.go menjadi entry utama
RUN go build -o main main.go

# Stage: Run
FROM alpine:3.18
WORKDIR /app
# hanya copy binernya saja
COPY --from=builder /app/main .
# Copy .env yang dipunya
COPY dev.env .

# Container berjalan pada port 8080
EXPOSE 8080
# Lakukan command setelah kontainer dimulai
CMD [ "/app/main" ]