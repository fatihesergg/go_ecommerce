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

<details>
<summary>Proje Yapısı (Tıklayarak Göster/Gizle)</summary>

```plaintext
├── cmd                           // Ana uygulamanın bulunduğu klasör.
│   └── go_ecommerce              // Uygulamanın çalıştırılabilir dosyalarının yer aldığı klasör.
│       └── main.go               // Uygulamanın giriş noktası (main fonksiyonu).
├── docs                          // Swagger dokümantasyon dosyaları.
│   ├── docs.go                   // Dokümantasyon yapılandırma dosyası.
│   ├── swagger.json              // JSON formatında Swagger dokümantasyonu.
│   └── swagger.yaml              // YAML formatında Swagger dokümantasyonu.
├── go.mod                        // Proje modül dosyası (bağımlılıklar).
├── go.sum                        // Modül bağımlılıklarının kontrol toplamları.
├── internal                      // Uygulamanın çekirdek işlevselliği burada yer alır.
│   ├── dto                      // Data Transfer Object'ler (veri aktarım yapıları).
│   │   ├── auth.go              // Kimlik doğrulama ile ilgili DTO'lar.
│   │   ├── category.go          // Kategori işlemleri için DTO'lar.
│   │   ├── order.go             // Sipariş işlemleri için DTO'lar.
│   │   ├── product.go           // Ürün işlemleri için DTO'lar.
│   │   └── review.go            // Ürün inceleme/yorum işlemleri için DTO'lar.
│   ├── handler                  // API endpoint handler'ları (HTTP isteklerini işleyen fonksiyonlar).
│   │   ├── auth.go              // Authentication ile ilgili endpoint'ler.
│   │   ├── category.go          // Kategori ile ilgili endpoint'ler.
│   │   ├── order.go             // Sipariş ile ilgili endpoint'ler.
│   │   ├── payment.go           // Ödeme işlemleri (Stripe) ile ilgili endpoint'ler.
│   │   ├── product.go           // Ürün ile ilgili endpoint'ler.
│   │   └── review.go            // İnceleme/yorum ile ilgili endpoint'ler.
│   ├── middleware               // HTTP middleware'leri (ör. JWT doğrulama, loglama).
│   │   └── middleware.go        // Ortak middleware fonksiyonları.
│   ├── model                    // GORM modelleri (veritabanı şema tanımlamaları).
│   │   └── model.go             // Tüm model tanımlamaları.
│   ├── service                  // İş mantığı ve servis katmanı.
│   │   └── service.go           // İş mantığına ait fonksiyonlar ve işlemler.
│   ├── storage                  // Repository Pattern implementasyonu (veritabanı işlemleri).
│   │   └── storage.go           // CRUD işlemleri ve veritabanı etkileşimleri.
│   └── util                     // Yardımcı fonksiyonlar ve genel araçlar.
│       ├── errors.go            // Hata yönetimi ve özel hata mesajları.
│       └── util.go              // Genel yardımcı fonksiyonlar.
```
</details>

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
```
API Temel URL: http://localhost:8080
Swagger Dokümantasyonu: http://localhost:8080/swagger/index.html
