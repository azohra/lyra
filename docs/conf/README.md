# Lyra Configuration

See `lconf.toml` for an example of a lyra configuration file. This conf file lives inside `~/.lyra`. 

# Specs
Configuration is done via a toml file. The following are the exact specification of the toml file.

* `Vault` table:
  * `enable = false`
    * Enables vault, default set to false.
  * `server = "http://localhost:8200"` 
    * The url of the vault server. Default value is set to `http://localhost:8200`
  * `ca_cert = [""]` 
    * An array of CA certs (in PEM format) to be trusted. Usually this should be only filled if you are self signing your vault's cert. Lyra by defaults always trust your operating system's default trust store.