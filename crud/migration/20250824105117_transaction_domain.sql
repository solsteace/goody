-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE `transaksi`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_user` INT UNSIGNED NOT NULL,
    `alamat_pengiriman` INT UNSIGNED,
    `harga_total` INT UNSIGNED NOT NULL,
    `kode_invoice` VARCHAR(255) NOT NULL,
    `metode_bayar` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`alamat_pengiriman`)
        REFERENCES `alamat`(`id`)
        ON DELETE SET NULL,
    FOREIGN KEY (`id_user`)
        REFERENCES `users`(`id`)
        ON DELETE CASCADE
);

CREATE TABLE `log_produk`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_produk` INT UNSIGNED,
    `id_toko` INT UNSIGNED NOT NULL,
    `nama_produk` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(255) NOT NULL,
    `harga_reseller` INT UNSIGNED NOT NULL,
    `harga_konsumen` INT UNSIGNED NOT NULL,
    `deskripsi` TEXT,
    `created_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`id_produk`)
        REFERENCES `produk`(`id`)
        ON DELETE SET NULL,
    FOREIGN KEY (`id_toko`)
        REFERENCES `toko`(`id`)
        ON DELETE CASCADE
);

CREATE TABLE `detail_transaksi`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_transaksi` INT UNSIGNED NOT NULL,
    `id_log_produk` INT UNSIGNED,
    `id_toko` INT UNSIGNED NOT NULL,
    `kuantitas` INT UNSIGNED NOT NULL,
    `harga_total` INT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`id_toko`)
        REFERENCES `toko`(`id`)
        ON DELETE CASCADE,
    FOREIGN KEY (`id_transaksi`)
        REFERENCES `transaksi`(`id`)
        ON DELETE CASCADE,
    FOREIGN KEY (`id_log_produk`)
        REFERENCES `log_produk`(`id`)
        ON DELETE SET NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS `detail_transaksi`;
DROP TABLE IF EXISTS `log_produk`;
DROP TABLE IF EXISTS `transaksi`;