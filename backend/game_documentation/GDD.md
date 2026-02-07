# Game Design Document (GDD)
Proje: `My2DHobieMMORPGGame`  
Dosya: `game_documentation/GDD.md`  
Yazar: Asistan (referans taslak) — senin okuman ve üzerinde değişiklik yapman için hazırlandı.

---

## 1. Yüksek Seviye Özet
- Tür: 2D isometrik, medieval fantasy, MMORPG (hobi amaçlı, lightweight).
- Tema / Hikaye: Daemonic varlıklar dünyayı ele geçiriyor; NPC'leri öldürüyor veya onları "daemon" haline getiriyor. Oyuncu (kahraman) bu istilayı durdurmaya çalışıyor.
- Hedef platform: Web (Phaser ile render), ileride isteğe bağlı native wrapper (Electron / Tauri) olabilir.
- Teknoloji yığını (senin seçimin):
  - Backend: `Golang`, `Fiber`, `Gorm` (başlangıç), `Sqlite` (başlangıç), `Testify`
  - Headless physics: kendi yazdığın basit, server-otoriter fizik motoru (koda özel)
  - Frontend: `Phaser` (render, kamera, input), `Angular` (UI container) — fakat fizik tamamen backend'de
  - Level / map: `ldtk-ts` (https://github.com/jprochazk/ldtk-ts) ile LDtk harita dosyalarını frontendde parse etme
  - DevOps: `Docker`, `github cli` vb.

Not: README dosyasındaki teknoloji ağacını referans olarak ekledim: `Readme.md`.

---

## 2. Oyun Tasarımının Temel İlkeleri
- Tam otoriter sunucu: Tüm oyun mantığı, fizik, çarpışma, combat resolution sunucuda çalışır. Client salt görüntüleyici ve input göndericidir.
- Lightweight hedef: Basit veri paketleri, delta-snapshot, AOI (area of interest) ile gereksiz ağ trafiğini azalt.
- Deterministik kesinlik değil ama tutarlılık: Single authoritative server çalıştığı sürece deterministik şart şart değil; fixed timestep ile tutarlılık sağlanır.
- Güvenlik odaklı: Hile ihtimalini minimize etmek için doğrulama, limitler, rate-limiting, input sequencing.

---

## 3. Core Gameplay Loop
1. Oyuncu input gönderir (örn. hareket yönü, yetki kullanımı, etkileşim) — her input `seq` numarası taşır.  
2. Sunucu input'u alır, validasyon yapar (speed cap, pozisyon geçerliliği vb.), fizik simülasyonuna uygular (fixed tick).  
3. Sunucu her tick üretir: snapshot veya delta prepare edilir. Sadece ilgili (AOI içindeki) oyunculara gönderilir.  
4. Client snapshot alır, interpolation / reconciliation uygular (client-side prediction optional).  
5. Combat / loot / quest gibi olaylar sunucu tarafından resolve edilir; sonuçlar persist edilir.

---

## 4. Dünya ve Haritalar
- Harita formatı: LDtk ile seviyeler tasarlanır. Frontend `ldtk-ts` ile haritayı parse eder; tile atlas'ı ve collider layer'ları backend ile eşleştirilir.
- World structure:
  - Global world (persistent zones)
  - Instance / dungeon (geçici; oyuncu seti için izole)
  - Region partition: sunucu ölçeklenebilirliği için world bölgelere ayrılır (ör. 128x128 tile chunk'lar)
- Map metadata: her harita LDtk içinde collision layer, spawn noktaları, navmesh veya waypoint graph (AI için) içermeli.

Referans: `ldtk-ts` repo — https://github.com/jprochazk/ldtk-ts

Örnek görsel referanslar (ilham için — bunları inceleyip beğenirsen asset yolunu planla):
- Unsplash "fantasy" arama sayfası: https://unsplash.com/s/photos/fantasy  
- Pixabay "fantasy" arama sayfası: https://pixabay.com/images/search/fantasy/  
- ArtStation "fantasy" arama sayfası: https://www.artstation.com/search?q=fantasy

(Bunlar doğrudan konsept galerileri; ticari kullanım veya asset alımı öncesi lisans kontrolü yap.)

---

## 5. Entity ve Veri Modelleri (yüksek seviye)
- `Player`:
  - id, username, authToken, position (x,y), velocity (vx,vy), heading, hp, mp, stats, inventoryRef, activeEffects, seqNumber
- `NPC`:
  - id, type, position, state (idle, aggro, fleeing), hp, aiState
- `Item`:
  - id, type, stackCount, position (world veya inventory), attributes
- `Projectile` / `AOE`:
  - id, ownerId, position, velocity, lifetime, effect
- `WorldObject`:
  - id, type (static obstacle, door, chest), position, collisionShape

Not: Serileştirme için ikili format (Protobuf ya da MessagePack) tercih et — JSON çok fazla trafik üretebilir.

Örnek mesaj alanları (yarım-codish, açıklayıcı):
- Input paket: `{ playerId, seq, timestamp, inputs: [{type:'move', dir:[x,y]}, {type:'attack', targetId}] }`
- Snapshot paket: `{ tick, entities: [{id, pos:{x,y}, state}, ...], events: [...] }`

(Tam şema istersen `.proto` örneği hazırlayabilirim.)

---

## 6. Headless Physics Motor - Tasarım Detayları
Sana özel, basit ve hafif bir headless fizik motoru taslağı:

1. Fixed timestep
   - Sabit güncelleme aralığı: örn. 20ms (50Hz) veya 16ms (60Hz). Tüm fizik hesaplamaları bu tikte yapılır.
   - Sunucu tick sayacı: `tick++` her update'te snapshot oluşturma için referans.

2. Temel hareket modeli
   - Kinematic hareket: pozisyon = pozisyon + velocity * dt
   - Input sadece hedef yön/komut belirtir (örn. move toward vector), hız ve acceleration sunucuda clamplenir.
   - Max speed / acceleration / friction parametreleri ile çok hızlı hileleri engelle.

3. Çarpışma modeli
   - Basit şekiller: `circle` veya `AABB` (axis-aligned bounding box). Her ikisini desteklemek yeterli.
   - Collision detection: spatial hashing (grid) veya quadtree ile geniş ölçekli hızlandırma.
   - Çarpışma çözümü: pozisyon separasyonu (push-out) + küçük impulse düzeltmesi. Karmaşık gerçeğe dayalı reaksiyonlar gerekmiyorsa bu yeterli.
   - Statik colliders: tile-based collision layer (LDtk collider layer) sunucuda yüklenir. Dinamik colliders: oyuncular, NPC'ler, projeler.

4. Deterministiklik ve sayısal kararlılık
   - Eğer multi-instance veya snapshot replay planlıyorsan fixed-point (integer) veya double + deterministik rounding stratejisi uygula.
   - Başlangıç için `float64` + sabit timestep yeterli; farklı platformlarda bir sunucuda şart.

5. Spatial Partitioning / Interest Management
   - AOI: her oyuncu için görünür yarıçap (ör. 10 ekran çapı). Sunucu yalnızca bu bölge içindeki entity'leri gönderir.
   - Implementation: grid-based buckets (ör. 64x64 world units). Her entity bucket'a kayıtlı, her tick güncelleme bucket tabanlı sorgu ile alınır.

6. Oyun mantığı ile entegrasyon
   - Collision callback'leri: onEnter, onStay, onExit (ör. trigger bölgeleri için)
   - Damage / combat: collision detection sonucuna göre trigger ile combat system'e event oluşturulur — combat resolve sunucuda yapılır.

7. Optimization / scaling
   - Update frequency: NPC ve uzak entity'ler için daha düşük update rate uygula.
   - Snapshot delta: tam entity listesi yerine değişen alanları gönder.
   - Use Redis pub/sub veya internal job queue (NATS gibi) eğer birden fazla process / instance arasında state paylaşımı gerekiyorsa.

---

## 7. Combat Sistemi (örnek)
- Temel: hit-based, skill cooldown'ları, resource (mana/stamina) yönetimi.
- Saldırı akışı:
  - Client input => `cast skill` (skillId, targetId veya direction)
  - Server validation (menzilde mi, cooldown hazır mı, mana yeterli mi)
  - Eğer valid: skill effect spawn (projeler/aoe) -> physics hesapla -> hit sonuçlarını uygulayıp event publish et
- Damage formula (basit öneri): `damage = baseDamage * (1 + atk - defModifier)` (daha sofistike scaling isteğe bağlı)
- Crit / resist mekanikleri: ek katman olarak eklenebilir.

---

## 8. NPC / AI
- Basit state machine: `Idle` -> `Patrol` -> `Alert` -> `Chase` -> `Attack` -> `Flee` -> `Dead`
- Pathfinding: grid-based A* veya waypoint graph (LDtk'de oluşturulan nav nodes). Sunucuda hafif pathfinding implementasyonu yeterli.
- Ağır yol bulma işlemleri için bucketlı görev yöneticisi (batch olarak pathfinding).

---

## 9. Persistence ve Veri Tabanı
- Başlangıç: `Sqlite` (kolay, local dev için ideal). Ancak yazma çakışmaları ve yüksek eşzamanlılık için darboğaz olabilir.
- Üretim/ölçekleme planı: `Postgres`'e geçiş planı hazırla (migrations ile). Ayrıca ephemeral state / cache için `Redis`.
- Örnek persistans:
  - Accounts, player profiles, inventory, persistent world state, quest progress -> RDBMS
  - Realtime ephemeral state (player position, transient NPCs) -> in-memory (server process) + optionally Redis for cross-process

---

## 10. Ağ Protokolü ve Veri Serileştirme
- Transport: WebSocket (tarayıcı için) — gofiber websocket middleware veya `gorilla/websocket`.
- Serileştirme: Protobuf veya MessagePack tercih et (binary, düşük boyut). JSON sadece prototipte.
- Paket türleri:
  - Reliable ordered (game commands, chat, inventory ops) => WebSocket reliable channel yeter.
  - Real-time state (position updates) => optimize edilmiş binary deltas.
- Input sequencing: her client paketine `seq` ekle; server ack/feedback ile reconciliation mümkün.

---

## 11. UI / Frontend Notları
- Phaser kullanımı: sadece render, camera, sprite management, interpolation. Fizik client'ta olmayacak.
- `ldtk-ts` frontend parse: tilemap yükle, collider layer'ları görsel olarak render et; ama gerçek collision sunucu tarafında hesaplanır.
- UI: Angular container, menüler, inventory, skill bars, chat, map.
- Client-side prediction: optional. Eğer eklersen, reconciliation ile glitchleri düzelt.

---

## 12. Test, Monitoring ve Geliştirme Araçları
- Local dev ortam: `docker-compose` ile `postgres`/`redis` (ileriki aşama), go server, frontend container.
- Load testing: fake-player simülatörleri (goroutine tabanlı) ile 50/100/500 oyuncu senaryoları. AOI ve tick darboğazlarını ölç.
- Metrics: tick time, avg processing per tick, network bytes per client, DB latency. Prometheus + Grafana önerilir.
- CI: GitHub Actions + test suite (unit/integration).

---

## 13. Roadmap / Milestones (Öneri)
1. PoC: Auth + basit world + WebSocket, player movement (server authoritative) + basic snapshot (2 hafta)
2. LDtk integration (frontend) + basic tile colliders (1 hafta)
3. Headless physics engine v1: fixed timestep, kinematic movement, simple collision (2-3 hafta)
4. Combat ve NPC basics, quest örneği (2-3 hafta)
5. Persistence: SQLite -> migration path to Postgres, Redis cache (2 hafta)
6. Load testing & optimization (2 hafta)
7. Polish: UI/UX, art assets, audio (sürekli ilerleyen adımlar)

---

## 14. Asset & Art Notları
- Senin dediğin gibi özel konsept yoksa ilham almak için yukarıdaki arama sayfalarını kullan. Ücretsiz ve telifsiz art almayı düşünüyorsan `Unsplash`/`Pixabay`/`OpenGameArt` gibi kaynakları incele.
- Tile atlas ve sprite pack: 2D izometrik tileset bul veya kendin LDtk uyumlu tileset oluştur. LDtk ile uyumlu eksport formatlarını kontrol et.

Örnek asset arama:
- Unsplash fantasy: https://unsplash.com/s/photos/fantasy  
- OpenGameArt: https://opengameart.org/ (oyun varlıkları için uygun lisanslar bulunur)

---

## 15. Ek Notlar ve Öneriler
- Başlangıçta çok az hazır bağımlılık kullan; minimal, test edilebilir modüller oluştur. Bu sayede ileride gerektiğinde `ByteArena/box2d` veya `Rapier` gibi kütüphanelere geçiş yapabilirsin.
- `ldtk-ts` ile frontend tarafında LDtk dosyalarını parse ederken backend ile collider/tile ID eşleşmesine dikkat et — coordinate transform (pixel -> world units) konvansiyonunu netleştir.
- Hile engelleme: hareket limitleri, server-side reconciliation, anti-speed-hack logic, input rate limit ve salt doğrulamalar yeterli başlangıç katmanlarıdır.

---

## 16. Sonraki Adımlar (Benim tarafımdan yardım teklifleri)
Eğer istersen, minimal AI müdahalesi ile aşağıdakilerden birini hazırlarım:
- Basit `Message` `.proto` şeması (movement, snapshot, chat, auth) — client/server aynı şemayı kullansın.
- Headless physics motoru için adım-adım pseudo-algoritma ve veri akışı şeması.
- PoC sunucu iskeleti (Go + Fiber + WebSocket) — küçük bir repo iskeleti (senin minimal AI tercihini göz önünde tutarak kısa ve okunabilir).
- `LDtk` -> `ldtk-ts` frontend parse örneği açıklaması (Angular + Phaser entegrasyon notları).

Hangi yardımı istersin? Ben rehberlik ve taslak kod / şema üretiminde yardımcı olabilirim; senin tercihine göre daha az veya daha fazla otomatik içerik oluşturmam mümkün.

---

Not: Bu doküman başlangıç taslağıdır — senin değerlendirmelerin ve önceliklerine göre güncelleriz. Özellikle:
- hedef eşzamanlı oyuncu sayısını belirtirsen (örn. 50 / 500 / 5000), sunucu ve DB önerilerini daha kesin hale getiririm.
- gerçek asset lisans gereksinimi varsa bunu da planlarımıza ekleriz.
