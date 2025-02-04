# E-commerce API Projesi

Bu proje, Golang kullanılarak geliştirilmiş bir E-commerce API'sidir. Projede JWT tabanlı kimlik doğrulama, Stripe ödeme entegrasyonu, Repository Pattern, GORM ile Postgresql veritabanı işlemleri, Zap structured logger ile loglama ve Swagger dokümantasyonu (swaggo/swag) gibi teknolojiler kullanılmıştır.

## Özellikler

- **JWT Tabanlı Kimlik Doğrulama:** Kullanıcı girişi ve yetkilendirme işlemleri için JWT token kullanımı.
- **Stripe Entegrasyonu:** Ödeme işlemleri için Stripe API entegrasyonu.
- **Repository Pattern:** Katmanlı mimari ve test edilebilirlik için repository pattern kullanımı.
- **GORM & Postgresql:** ORM aracı GORM ile Postgresql veritabanı yönetimi.
- **Zap Structured Logger:** Performanslı ve yapılandırılmış loglama.
- **Swagger Dokümantasyonu:** API endpoint'lerinin dokümantasyonu ve test edilebilirliği için swaggo/swag entegrasyonu.

## Gereksinimler

- Golang
- Postgresql
- Stripe API Key
- Git

## Kurulum

1. **Repository’i Klonlayın:**

    ```bash
    git clone https://github.com/kullanici_adi/ecommerce-api.git
    cd ecommerce-api
    ```

2. **Gerekli Go Paketlerini Yükleyin:**

    ```bash
    go mod download
    ```

3. **Environment Değişkenleri Ayarlayın:**

    ```bash
    export STRIPE_API=your_stripe_api_key
    export JWT_SECRET=your_jwt_secret
    ```

4. **Veritabanı ve Migrasyon İşlemleri:** 
    [main.go](https://github.com/fatihesergg/go_ecommerce/blob/main/cmd/go_ecommerce/main.go) Dosyasındaki "dsn" satırını kendi database bağlantınızla değiştirin.

## Çalıştırma

Projeyi yerel makinenizde çalıştırmak için:

```bash
go run cmd/go_ecommerce/main.go
