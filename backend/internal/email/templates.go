package email

import "fmt"

// Templates E01-E05 per Section 17 PRD. Kept as plain inline HTML (no template
// engine) since the email set is small and unlikely to grow much before a proper
// design pass in a later phase.

func WelcomeEmail(name string) (subject, html string) {
	subject = "Selamat datang di KickPick!"
	html = fmt.Sprintf(`<p>Hai %s,</p><p>Terima kasih sudah bergabung dengan KickPick. Sekarang kamu bisa membandingkan harga sepatu dari berbagai brand, simpan wishlist, dan pasang alert restock.</p>`, name)
	return
}

func VerifyEmail(appURL, token string) (subject, html string) {
	subject = "Verifikasi email kamu"
	link := fmt.Sprintf("%s/id/verifikasi?token=%s", appURL, token)
	html = fmt.Sprintf(`<p>Klik link berikut untuk verifikasi email kamu (berlaku 24 jam):</p><p><a href="%s">%s</a></p>`, link, link)
	return
}

func ResetPasswordEmail(appURL, token string) (subject, html string) {
	subject = "Reset password KickPick kamu"
	link := fmt.Sprintf("%s/id/reset-password?token=%s", appURL, token)
	html = fmt.Sprintf(`<p>Klik link berikut untuk reset password kamu (berlaku 1 jam, hanya bisa dipakai sekali):</p><p><a href="%s">%s</a></p><p>Kalau kamu tidak meminta reset password, abaikan email ini.</p>`, link, link)
	return
}

func RestockAlertEmail(productName, productURL string) (subject, html string) {
	subject = fmt.Sprintf("%s baru saja restock!", productName)
	html = fmt.Sprintf(`<p><strong>%s</strong> baru saja tersedia lagi.</p><p><a href="%s">Lihat sekarang</a></p>`, productName, productURL)
	return
}

func PriceDropAlertEmail(productName, productURL string) (subject, html string) {
	subject = fmt.Sprintf("Harga %s turun sekarang", productName)
	html = fmt.Sprintf(`<p>Harga <strong>%s</strong> baru saja turun.</p><p><a href="%s">Lihat sekarang</a></p>`, productName, productURL)
	return
}
