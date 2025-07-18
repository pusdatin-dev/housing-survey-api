package seed

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func ProgramSeed(db *gorm.DB) {
	fmt.Println("Running Program Seeder...")

	program := []models.Program{
		// Anggaran Negara - Pembangunan Baru
		{Name: "Rumah Susun", Detail: "Anggaran PKP", ResourceID: 1},
		{Name: "Rumah Khusus", Detail: "Anggaran Non PKP", ResourceID: 1},
		{Name: "DAK Tematik PPKT", Detail: "Anggaran Non PKP", ResourceID: 1},
		{Name: "APBD - PB", Detail: "Anggaran Non PKP", ResourceID: 1},

		// Anggaran Pengembang - Pembangunan Baru
		{Name: "Sederhana Non FLPP", Detail: "", ResourceID: 2},
		{Name: "Menengah Kebawah", Detail: "", ResourceID: 2},
		{Name: "Operasionalisasi BP3", Detail: "", ResourceID: 2},
		{Name: "Peng. Usaha Modular", Detail: "", ResourceID: 2},
		{Name: "Fasilitasi Pemanfaatan Tanah", Detail: "", ResourceID: 2},

		// Anggaran Swadaya - Pembangunan Baru
		{Name: "Swadaya Masyarakat PBG (PB)", Detail: "", ResourceID: 3},
		{Name: "Swadaya Masyarakat Non PBG (PB)", Detail: "", ResourceID: 3},

		// Anggaran Gotong Royong - Pembangunan Baru
		{Name: "CSR (PB)", Detail: "", ResourceID: 4},

		// Anggaran Negara - Peningkatan Kualitas
		{Name: "Penanganan Kumuh", Detail: "Anggaran PKP", ResourceID: 5},
		{Name: "BSPS", Detail: "Anggaran PKP", ResourceID: 5},
		{Name: "Dana Desa", Detail: "Anggaran Non PKP", ResourceID: 5},
		{Name: "RTLH Kemensos", Detail: "Anggaran Non PKP", ResourceID: 5},
		{Name: "APBD - PK", Detail: "Anggaran Non PKP", ResourceID: 5},
		{Name: "DAK Tematik PPKT", Detail: "Anggaran Non PKP", ResourceID: 5},

		// Anggaran Pembiayaan - Peningkatan Kualitas
		{Name: "Tabungan Kontrak Perumahan", Detail: "", ResourceID: 6},

		// Anggaran Swadaya - Peningkatan Kualitas
		{Name: "Swadaya Masyarakat PBG (PK)", Detail: "", ResourceID: 7},
		{Name: "Swadaya Masyarakat Non PBG (PK)", Detail: "", ResourceID: 7},

		// Anggaran Gotong Royong - Peningkatan Kualitas
		{Name: "CSR (PK)", Detail: "", ResourceID: 8},

		// Anggaran Pengembang - Pembangunan Baru-Upaya Eksternal
		{Name: "FLPP 2025", Detail: "", ResourceID: 9},
		{Name: "FLPP Tambahan 2025", Detail: "", ResourceID: 9},
		{Name: "Pelonggaran GWM BI", Detail: "", ResourceID: 9},
		{Name: "Tambahan Mikro / Milenial", Detail: "", ResourceID: 9},

		// Anggaran Investasi - Pembangunan Baru-Upaya Eksternal
		{Name: "Investasi Luar Negeri", Detail: "", ResourceID: 10},

		// Anggaran Gotong Royong - Pembangunan Baru-Upaya Eksternal
		{Name: "CSR Inisiasi Menteri PKP", Detail: "", ResourceID: 11},
	}
	for _, p := range program {
		if err := db.FirstOrCreate(&p, models.Program{Name: p.Name, ResourceID: p.ResourceID}).Error; err != nil {
			log.Printf("Error seeding Program: %v", err)
		}
	}
	fmt.Println("Finished Program Seeder...")
}
