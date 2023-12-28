# Go Broker

Go Broker is a lightweight and configurable data broker implemented in Go. It serves as a bridge between various data sources and destinations, allowing for easy data ingestion and transmission. This project is particularly useful for scenarios where data needs to be fetched from one or more sources, transformed, and then delivered to a specific destination.

## Features

- **Configurability**: Easily configure data sources and destinations using YAML files.
- **Support for Various Input Connectors**: Currently, Go Broker supports the HTTP connector for fetching data from remote sources. More connectors can be added in the future.
- **Flexible Schema Handling**: Define a schema to transform and filter data as it flows through the broker.
- **Elasticsearch Integration**: Send processed data to Elasticsearch for efficient storage and retrieval.

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/go-broker.git
   cd go-broker
    ```
2. **Create a configuration file (e.g., config.yaml) to specify your input and output settings:**
   ```yaml
    input:
      type: csv
      connector: http
      locations:
        - "https://raw.githubusercontent.com/mthcht/ThreatHunting-Keywords/main/signature_keyword.csv"
        #  schema:
        #    - "rickroll": "rickroll"
    output:
      connector: elasticsearch
      host: ELASTICSEARCH_URL
      username: ELASTICSEARCH_USERNAME
      password: ELASTICSEARCH_PASSWORD
      update: new
      store: go-broke
    ```
3. **Run the Go Broker:**
   ```bash
    go run main.go -config=config.yaml
    ```
# Configuration
The configuration file (config.yaml) is used to specify the input and output settings of the broker. Here's a breakdown of the available configuration options:

## Input Configuration
- **type**: Type of input data (e.g., csv).
- **connector**: Input connector to use (e.g., http).
- **locations**: List of URLs or paths to the input data.
- **schema**: (Optional) Define a schema for data transformation.

## Output Configuration
- **connector**: Output connector to use (e.g., elasticsearch).
- **host**: Elasticsearch server URL.
- **username**: Elasticsearch username (if applicable).
- **password**: Elasticsearch password (if applicable).
- **update**: Update strategy for Elasticsearch data (e.g., new, overwrite).
- **store**: Elasticsearch index to store the data.

# License
This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE.md) file for details.

