package storage

import (
	"log"

	"github.com/gofiber/storage/redis/v3"
)

var RediStorage *redis.Storage

func InitRedis() {
	RediStorage = redis.New(redis.Config{
		Host:      "127.0.0.1", // Alamat server Redis
		Port:      6379,        // Port default Redis
		Username:  "",          // Gunakan username jika ada
		Password:  "",          // Gunakan password jika ada
		Database:  0,           // Database Redis
		Reset:     false,       // Jika true, semua data Redis akan dihapus saat inisialisasi
		TLSConfig: nil,         // Config TLS jika Redis menggunakan TLS
		PoolSize:  10,          // Ukuran pool koneksi Redis
		// IdleTimeout: 5 * time.Minute, // Timeout untuk idle connections
	})

	// Try saving and retrieving data to make sure Redis is working
	err := RediStorage.Set("test-key", []byte("Hello Redis!"), 0)
	if err != nil {
		log.Fatalf("Error save data in Redis: %v", err)
	}

	log.Println("Redis berhasil diinisialisasi dan terhubung!")
}
