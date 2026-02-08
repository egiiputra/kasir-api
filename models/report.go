package models

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtySold    int    `json:"qty_terjual"`
}

type ReportResponse struct {
	TotalRevenue    int             `json:"total_revenue"`
	TotalTransaksi  int             `json:"total_transaksi"`
	ProdukTerlaris  ProdukTerlaris  `json:"produk_terlaris"`
}
