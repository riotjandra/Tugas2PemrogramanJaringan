package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Menyimpan semua koneksi client yang terhubung
var (
	clients    []net.Conn
	clientsMux sync.Mutex // Mutex untuk mengamankan akses ke slice clients
)

// Fungsi untuk menangani setiap koneksi client
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Menambahkan client ke daftar client
	clientsMux.Lock()
	clients = append(clients, conn)
	clientsMux.Unlock()

	// Membaca username dari client
	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Gagal membaca username:", err)
		return
	}
	username = strings.TrimSpace(username)
	fmt.Printf("Client terhubung: %s\n", username)

	// Mengirim pesan ke semua client yang terhubung
	for {
		// Membaca pesan dari client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Koneksi terputus dari %s\n", username)
			removeClient(conn) // Hapus client yang terputus dari daftar
			return
		}
		message = strings.TrimSpace(message)

		// Jika klien mengirim "/quit", tutup koneksi
		if message == "/quit" {
			fmt.Printf("%s telah keluar.\n", username)
			removeClient(conn)
			return
		}

		// Kirim pesan ke semua client yang terhubung
		broadcastMessage(fmt.Sprintf("%s: %s", username, message), conn)
	}
}

// Fungsi untuk mengirim pesan ke semua client yang terhubung
func broadcastMessage(message string, sender net.Conn) {
	clientsMux.Lock() // Kunci akses ke clients
	defer clientsMux.Unlock()

	for _, client := range clients {
		// Jangan kirim kembali ke pengirimnya sendiri
		if client != sender {
			_, err := fmt.Fprintln(client, message)
			if err != nil {
				fmt.Println("Gagal mengirim pesan:", err)
			}
		}
	}
}

// Fungsi untuk menghapus client dari daftar client
func removeClient(conn net.Conn) {
	clientsMux.Lock() // Kunci akses ke clients
	defer clientsMux.Unlock()

	for i, client := range clients {
		if client == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func main() {
	// Membuat listener server pada port 8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Gagal membuka server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server berjalan di localhost:8080...")

	// Menunggu koneksi dari client
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Gagal menerima koneksi:", err)
			continue
		}

		// Membuat goroutine untuk menangani koneksi client
		go handleConnection(conn)
	}
}
