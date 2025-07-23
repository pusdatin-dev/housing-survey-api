package shared

var (
	StatusDraft           = "Draf"
	StatusWaitingBalai    = "Laporan Menunggu Verifikasi Balai"    // Report Waiting for Balai Verification
	StatusWaitingEselon1  = "Laporan Menunggu Verifikasi Eselon 1" // Report Waiting for Eselon 1 Verification
	StatusVerified        = "Laporan Terverifikasi"                // Report Verified
	StatusRejectedBalai   = "Laporan Reject (Balai)"               // Report Rejected
	StatusRejectedEselon1 = "Laporan Reject (Eselon 1)"            // Report Rejected by Eselon 1

	Pending  = "Pending"  // General Pending Status
	Approved = "Approved" // General Approved Status
	Rejected = "Rejected" // General Rejected Status

	StatusRealProses  = "Proses Pembangunan"
	StatusRealSelesai = "Sudah Realisasi"

	Create = "create"
	Update = "update"

	TagNegara       = "Negara"
	TagPengembang   = "Pengembang"
	TagSwadaya      = "Swadaya"
	TagGotongRoyong = "Gotong Royong"

	ListTagResource = []string{TagNegara, TagPengembang, TagSwadaya, TagGotongRoyong}

	StatusResolved   = "Resolved"   // Status for resolved comments
	StatusUnresolved = "Unresolved" // Status for unresolved comments

	PICSurvey = map[string]bool{
		"Surveyor":             true,
		"Admin Balai":          true,
		"Verificator Balai":    true,
		"Admin Eselon 1":       true,
		"Verificator Eselon 1": true,
		"Super Admin":          true,
	}
)
