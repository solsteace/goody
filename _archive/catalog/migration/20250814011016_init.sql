-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE `users`(
    `id` INT UNSIGNED PRIMARY KEY,
    `isAdmin` BOOLEAN
);

CREATE TABLE `kategori`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `nama_kategori` VARCHAR(255) NOT NULL,
    `created_at` DATE NOT NULL,
    `updated_at` DATE NOT NULL
);


CREATE TABLE `toko`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_user` INT UNSIGNED NOT NULL,
    `nama_toko` VARCHAR(255) NOT NULL,
    `url_foto` VARCHAR(255),
    `created_at` DATE NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATE NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`id_user`)
        REFERENCES `users`(`id`)
        ON DELETE CASCADE
);

CREATE TABLE `produk`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_toko` INT UNSIGNED NOT NULL,
    `id_kategori` INT UNSIGNED,
    `nama_produk` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(255) NOT NULL,
    `harga_reseller` INT UNSIGNED NOT NULL,
    `harga_konsumen` INT UNSIGNED NOT NULL,
    `stok` INT UNSIGNED NOT NULL,

    FOREIGN KEY (`id_kategori`)
        REFERENCES `kategori`(`id`)
        ON DELETE SET NULL,
    FOREIGN KEY (`id_toko`)
        REFERENCES`toko`(`id`)
        ON DELETE CASCADE
);

CREATE TABLE `foto_produk`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_produk` INT UNSIGNED NOT NULL,
    `url` VARCHAR(255),
    `created_at` DATE NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATE NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`id_produk`)
        REFERENCES `produk`(`id`)
        ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS `foto_produk`;
DROP TABLE IF EXISTS `produk`;
DROP TABLE IF EXISTS `toko`;
DROP TABLE IF EXISTS `kategori`;
DROP TABLE IF EXISTS `users`;
