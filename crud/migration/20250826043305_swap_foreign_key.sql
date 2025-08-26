-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE `detail_transaksi`
    DROP FOREIGN KEY `detail_transaksi_ibfk_3`;
ALTER TABLE `detail_transaksi`
    DROP COLUMN `id_log_produk`;

ALTER TABLE `log_produk`
    ADD COLUMN `id_detail_transaksi` INT UNSIGNED NOT NULL;

ALTER TABLE `log_produk`
    ADD CONSTRAINT `fk_log_produk_detail_transaksi`
    FOREIGN KEY (`id_detail_transaksi`)
        REFERENCES `detail_transaksi`(`id`);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE `log_produk`
    DROP FOREIGN KEY `fk_log_produk_detail_transaksi`;
ALTER TABLE `log_produk`
    DROP COLUMN `id_detail_transaksi`;

ALTER TABLE `detail_transaksi`
    ADD COLUMN `id_log_produk` INT UNSIGNED NULL;
ALTER TABLE `detail_transaksi`
    ADD CONSTRAINT `detail_transaksi_ibfk_3`
    FOREIGN KEY (`id_log_produk`)
        REFERENCES `log_produk`(`id`);
