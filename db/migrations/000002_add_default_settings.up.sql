-- Insert default system settings
INSERT OR IGNORE INTO system_settings (key, value, value_type, description, updated_at)
VALUES 
    ('locale', 'en', 'string', 'locale of app', datetime('now')),
    ('web_interface_enabled', 'false', 'bool', 'Enable web interface for spectators', datetime('now')),
    ('ui_language', 'en', 'string', 'UI language (ru, en)', datetime('now')),
    ('db_path', '', 'string', 'Path to database file', datetime('now')),
    ('db_api_user', 'admin', 'string', 'API user', datetime('now')),
    ('db_api_password_hash', '', 'string', 'API password hash (bcrypt)', datetime('now')),
    ('hardware_com_port', '', 'string', 'COM port for Arduino', datetime('now')),
    ('hardware_reader_type', 'EM4095', 'string', 'RFID reader type', datetime('now')),
    ('hardware_debounce_ms', '2000', 'int', 'Debounce delay in ms', datetime('now')),
    ('log_retention_years', '2', 'int', 'Log retention period for manual cleanup', datetime('now')),
    ('last_log_cleanup_date', '', 'string', 'Last manual log cleanup date', datetime('now')),
    ('app_version', '1.0.0', 'string', 'Application version', datetime('now')),
    ('schema_version', '1', 'int', 'Database schema version', datetime('now'));
