# 1.0.1

- Independent public app repository with its own builds, tests and releases.
- Removed unused installation-ACME code and dependencies.
- Configuration still contains only the pairing token.

# 1.0.0

- Jedyną opcją konfiguracji aplikacji HA jest token parowania.
- Usunięto opcje adresów, portów, ACME i certyfikatów z formularza oraz YAML.
- Stare ustawienia są automatycznie usuwane przez własny endpoint Supervisora.
- Zachowano parowanie, klucz urządzenia i certyfikat w danych aplikacji.
- Adres i port Core są wykrywane automatycznie; platforma i tunel nie wymagają
  ręcznej konfiguracji użytkownika.

# 0.2.0

- Automatyczny tunel mTLS bez zmian konfiguracji HTTPS w Core.
