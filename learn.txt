Perbedaan Wait Group, Mutex, Channels?
sebenarnya ke tiga hal ini digunakan untuk handling goroutine, tetapi ada beberapa hal yang membedakan:

1) Wait Group
    -   Antara Goroutine utama/main/code utama dengan goroutine fungsi lain tidak berjalan secara
        konkurensi secara langsung, harus diatur sedemikian rupa untuk membuat konkurensi
        seperti penerapan: `Add`, `Done`, dan `Wait`
    
    -   Tepatnya wait group ini dikhusukan untuk menunggu kelompok goroutine selesai untuk dijalankan
        , setelah itu baru menjalankan atau melanjutkan ke kode utama   

2) Mutex
    -   Menerapkan langkah penguncian (lock) code atau data bersama untuk mencegah konkuren yang konflik
    
    -   Digunakan untuk menghindari race condition

    -   Cocok untuk mengamankan akses data bersama, untuk menghindari konflik seperti race condition

    -   Bisa saja tidak bersifat konkurensi, tetapi saling tunggu untuk menyelesaikan goroutine yang ada
        dan melanjutkan kode yang sama

3) Channels
    -   Digunakan untuk komunikasi antar go routine

    -   memungkinkan goroutine untuk mengirim dan menerima data secara aman, 
        menjaga sinkronisasi antar goroutine secara bersamaan

    -   Sudah pasti konkurensi

Tentang Go Routine?
Jadi, sebenarnya goroutine ada go routine utama (fucn main) dan juga ada go routine yang dibuat
seperti menambahkan fungsi dengan awalan "go". 

Jika ada 2000 goroutine, bisa saja goruntime melakukan scheduling agar hanya beberapa goroutine dijalankan.
Tiap goroutine akan berjalan secara konkuren dengan goroutine lainya, termasuk goroutine utama.

Pada bagian scheduling, ada kemungkinan ada 1 goroutine lainya yang berjalan pada jalur yang sama,
jika goroutine yang pertama tersebut berhenti sejenak atau istirahat (ini bisa terjadi karena beberapa faktor)

Goroutine akan selalu berjalan, tetapi jika ada panic() atau return akan berhenti

Perbedaan channels tanpa buffer dengan menggunakan buffer?

Tanpa Buffer:
    -   Hanya menampung 1 data yang akan dikirim oleh goroutine pengirim. goroutine pengirim akan
        berhenti sejenak sampai data diterima oleh goroutine penerima

Menggunakan Buffer:
    -   Bisa menampung banyak data, pada saat channel masih bisa menampung data. goroutine pengirim
        masih bisa berjalan, tetapi pada saat sudah penuh, maka goroutine pengirim akan berhenti sampai
        ada goroutine penerima


=====================================================================================

Database Transaction
Langkah membuat data selalu persistance dengan pembaharuan yang ada sampai data yang telah
diperbaharui tersebut melakukan commit (jika sukses) atau rollback (jika gagal)

Deadlock
Biasanya jika melakukan Database Transaction, ada kemungkinan Deadlock jika tidak di-handling dengan
baik. Kondisi Deadlock ini terjadi ketika beberapa transaksi dalam Database berusaha
mengunci baris atau tabel dalam urutan yang berbeda dan saling menunggu satu sama lain

=====================================================================================

// Goroutine 1
	go func() {
		fmt.Println("Goroutine 1: Trying to lock muA")
		muA.Lock()
		defer muA.Unlock()

		fmt.Println("Goroutine 1: Locked muA, trying to lock muB")
		muB.Lock()
		defer muB.Unlock()

		fmt.Println("Goroutine 1: Locked muB")
	}()

	// Goroutine 2
	go func() {
		fmt.Println("Goroutine 2: Trying to lock muB")
		muB.Lock()
		defer muB.Unlock()

		fmt.Println("Goroutine 2: Locked muB, trying to lock muA")
		muA.Lock()
		defer muA.Unlock()

		fmt.Println("Goroutine 2: Locked muA")
	}()

	// Goroutine 3
	go func() {
		// Goroutine 3 menunggu Goroutine 1 dan Goroutine 2 selesai
		// sebelum mencoba mengunci mutex apa pun
		// Ini memastikan urutan yang telah Anda tentukan
		fmt.Println("Goroutine 3: Waiting for Goroutine 1 and 2 to finish")
	}()

	select {}

    -   Pada kode di atas goroutine 1 dan 2 melakukan lock untuk menuntaskan tugasnya sehingga sampai Unlock
    -   Dan goroutine tidak saling lock tugas satu sama lain yang dapat membuat Deadlock
    -   jadi jika goroutine 1 ini lock muA, maka goroutine 2 akan lock tugas yang lainya yaitu muB
    -   Lalu, untu goRoutine 3 ini akan pada posisi istirahat sejenak atau tidur, jika memang tidak adape kerjaan

	===============================================================================================

	Dalam DB Transaction harus mencakup 4 hal yaitu ACID:
	1) Atomicity

	2)

	3) Isolation
	   Properti ini menjamin bahwa satu transaksi tidak akan terpengaruh oleh 
	   transaksi lain yang sedang berlangsung secara bersamaan.

	4)