--changeset josephspurrier:1
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user_status;

--changeset josephspurrier:2
INSERT INTO `user_status` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);
--rollback TRUNCATE TABLE user_status;

--changeset josephspurrier:2
CREATE TABLE user (
    id VARCHAR(36) NOT NULL,
    
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    
    status_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
    
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    UNIQUE KEY (email),
    CONSTRAINT `f_user_status` FOREIGN KEY (`status_id`) REFERENCES `user_status` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user;