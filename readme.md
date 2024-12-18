## Prerequisites

*   **Go:** Go 1.23 or later.
*   **Docker:** Docker Desktop or Docker Engine.
*   **kind (Kubernetes IN Docker):** For local Kubernetes testing.
*   **kubectl:** Kubernetes command-line tool.

## Build and Run (Without Kubernetes)

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/](https://github.com/)dev-palkhe/filestore.git
    cd filestore
    ```

2.  **Build the server:**

    ```bash
    go build ./cmd/server
    ```

3.  **Build the client:**

    ```bash
    go build ./cmd/client
    ```

4.  **Run the server:**

    ```bash
    ./server
    ```

5.  **In a separate terminal, run the client commands:**

    ```bash
    ./store add test.txt       # Create test.txt first
    ./store ls
    ./store rm test.txt
        ./store update test.txt
    ./store wc
    ./store freq-words
    ```

## Build and Run (With Kubernetes using kind)

1.  **Create a kind cluster:**

    ```bash
    kind create cluster --name filestore-cluster
    ```

2.  **Build the Docker image:**

    ```bash
    docker build -t devasheesh22/filestore-app:latest . 
    ```

3.  **Load the Docker image into kind:**

    ```bash
    kind load docker-image devasheesh22/filestore-app:latest --name filestore-cluster 
    ```

4.  **Apply Kubernetes manifests:**

    ```bash
    kubectl apply -f k8s/manifests/deploy.yaml
    kubectl apply -f k8s/manifests/svc.yaml
    kubectl apply -f k8s/manifests/ingress.yaml 
    ```

5.  **Port-forward the service (Recommended for kind):**

    ```bash
    kubectl port-forward svc/store 8000:8000
    ```

6. **Change the default port for the client:**
    * Open `cmd/client/main.go`
    * Rebuild the client `go build ./cmd/client`

7.  **In a separate terminal, run the client commands (without the -s flag):**

    ```bash
    ./store add test.txt       
    ./store ls
    ./store rm test.txt
        ./store update test.txt
    ./store wc
    ./store freq-words
    ```
