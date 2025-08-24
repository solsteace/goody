-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE `users` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `nama` VARCHAR(255) NOT NULL,
    `kata_sandi` VARCHAR(255) NOT NULL,
    `no_telp` VARCHAR(255) NOT NULL UNIQUE,
    `tanggal_lahir` DATE NOT NULL,
    `jenis_kelamin` VARCHAR(255) NOT NULL,
    `tentang` LONGTEXT NOT NULL,
    `pekerjaan` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `id_provinsi` VARCHAR(255) NOT NULL,
    `id_kota` VARCHAR(255) NOT NULL,
    `is_admin` BOOLEAN NOT NULL,
    `updated_at` DATETIME NOT NULL 
        ON UPDATE CURRENT_TIMESTAMP
        DEFAULT CURRENT_TIMESTAMP,
    `created_at` DATETIME NOT NULL 
        DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE `alamat`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_user` INT UNSIGNED NOT NULL,

    `judul_alamat` VARCHAR(255) NOT NULL,
    `nama_penerima` VARCHAR(255) NOT NULL,
    `no_telp` VARCHAR(255) NOT NULL,
    `detail_alamat` VARCHAR(255) NOT NULL,
    `updated_at` DATE NOT NULL,
    `created_at` DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`id_user`)     
        REFERENCES `users`(`id`)
        ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS `alamat`;
DROP TABLE IF EXISTS `users`;