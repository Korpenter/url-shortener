# url-shortener
url-shortener HTTP API application allowing users to shorten or delete shortened links in one by one or in batches with the ability to delete get the list of shortened urls using cookies.

## Usage

The application's configuration can be adjusted either by utilizing command-line flags or environmental variables. Below are the available configuration options:

- **Server Address (`-a` or `SERVER_ADDRESS`)**: Specifies the address and port where the server will be hosted. The default value is `localhost:8080`.

- **Base URL (`-b` or `BASE_URL`)**: Defines the base URL address for the application. The default value is `http://localhost:8080`.

- **File Storage Path (`-f` or `FILE_STORAGE_PATH`)**: Sets the storage path for files within the application.

- **Secret Key (`-k` or `URL_SHORTENER_KEY`)**: Provides the secret key required for cryptographic operations.

- **PostgreSQL Database URI (`-d` or `DATABASE_DSN`)**: Specifies the URI for connecting to the PostgreSQL database.

- **Enable HTTPS (`-s` or `ENABLE_HTTPS`)**: Determines whether HTTPS should be enabled. The default is set to `false`.

- **TLS Certificate File (`-l` or `TLS_CERT_FILE`)**: Specifies the path to the TLS certificate file.

- **TLS Key File (`-t` or `TLS_KEY_FILE`)**: Specifies the path to the TLS key file.

- **Config File (`-c`)**: Indicates the path to the configuration file.

### Docker
Build container:

```shell
docker buildx build -t shortener .
```
Run it:
```shell
docker run shortener
```
