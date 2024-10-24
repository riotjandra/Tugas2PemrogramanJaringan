package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Menghubungkan ke server
	conn, err := net.Dial("tcp", "localhost:8080") // Ganti dengan alamat dan port server Anda jika berbeda
	if err != nil {
		fmt.Println("Gagal menghubungkan ke server:", err)
		return
	}
	defer conn.Close()

	// Membaca username dari pengguna
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan username: ")
	username, _ := reader.ReadString('\n')

	// Mengirim username ke server
	fmt.Fprintf(conn, username)

	// Membaca pesan dari server dalam goroutine terpisah
	go readMessages(conn)

	// Mengirim pesan ke server
	for {
		fmt.Print("Masukkan pesan: ")
		message, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, message)
	}
}

// Fungsi untuk membaca pesan dari server
func readMessages(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Koneksi terputus:", err)
			return
		}
		fmt.Print(message)
	}
}
