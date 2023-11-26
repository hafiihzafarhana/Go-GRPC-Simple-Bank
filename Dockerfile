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
# Instalasi curl
RUN apk add curl
# jalankan migrasi databasenya
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Stage: Run
FROM alpine:3.18
WORKDIR /app
# hanya copy binernya saja
COPY --from=builder /app/main .
# Copy packaga migrasi (/app/migrate merujuk pada package go migrate)
COPY --from=builder /app/migrate ./migrate
# Copy .env yang dipunya
COPY dev.env .
# COPY start.sh
COPY start.sh .
COPY wait-for.sh .
# Copy db/migration yang berisi file migrasi ke folder migration
COPY db/migration ./migration

# Container berjalan pada port 8080
EXPOSE 8080
# Lakukan command setelah kontainer dimulai
CMD [ "/app/main" ]
# Untuk menjalankan start.sh
# Jadi nanti CMD [ "/app/main" ] akan diteruskan di entrypoint ini
ENTRYPOINT ["/app/start.sh"]