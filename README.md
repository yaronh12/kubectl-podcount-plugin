# kubectl-podcount-yaron

A kubectl plugin that counts running pods in Kubernetes clusters. Supports both `kubectl` and OpenShift `oc` command-line tools.

## Features

- 🚀 Count running pods across all namespaces or specific namespaces
- 📊 Support for multiple namespaces (comma-separated)
- 🔧 Compatible with both `kubectl` and `oc` (OpenShift)
- 🔐 Uses standard Kubernetes authentication methods
- 🎯 Filters only "Running" phase pods
- 📋 Inherits all standard kubectl configuration flags

## Installation

### Prerequisites

- Go 1.19+ installed
- `kubectl` or `oc` command-line tool
- Access to a Kubernetes cluster

### Build and Install

1. **Clone the repository**:
   ```bash
   git clone <your-repo-url>
   cd count-pods-plugin
   ```

2. **Initialize Go modules and build**:
   ```bash
   go mod init kubectl-podcount-yaron
   go mod tidy
   go build -o kubectl-podcount_yaron kubectl-podcount-yaron.go
   ```

3. **Install for kubectl**:
   ```bash
   chmod +x kubectl-podcount_yaron
   sudo mv kubectl-podcount_yaron /usr/local/bin/
   ```

4. **Install for OpenShift oc (optional)**:
   ```bash
   sudo ln -s /usr/local/bin/kubectl-podcount_yaron /usr/local/bin/oc-podcount_yaron
   ```

5. **Verify installation**:
   ```bash
   kubectl plugin list
   # Should show kubectl-podcount_yaron
   ```

## Usage

### Basic Examples

```bash
# Count running pods in all namespaces
kubectl podcount_yaron

# Count running pods in a specific namespace
kubectl podcount_yaron -n kube-system

# Count running pods in multiple namespaces
kubectl podcount_yaron -n "default,kube-system,monitoring"
```

### OpenShift Examples

```bash
# Works with oc as well
oc podcount_yaron -n openshift-console
```

### Advanced Usage

```bash
# Use with specific kubeconfig
kubectl podcount_yaron --kubeconfig /path/to/config -n production

# Use with specific context
kubectl podcount_yaron --context my-cluster-context

# Connect to specific server
kubectl podcount_yaron --server https://my-k8s-api:6443 --token your-token
```

## Configuration

The plugin supports all standard kubectl configuration options:

| Flag | Description | Example |
|------|-------------|---------|
| `-n, --namespace` | Target namespace(s) | `-n "default,kube-system"` |
| `--kubeconfig` | Path to kubeconfig file | `--kubeconfig ~/.kube/config` |
| `--context` | Kubeconfig context to use | `--context production` |
| `--server` | Kubernetes API server URL | `--server https://api.k8s.io` |
| `--token` | Bearer token for authentication | `--token abc123...` |

### Cluster Connection Methods

1. **Default kubeconfig**: `~/.kube/config`
2. **KUBECONFIG environment variable**:
   ```bash
   export KUBECONFIG=/path/to/config
   kubectl podcount_yaron
   ```
3. **Explicit kubeconfig file**: `--kubeconfig /path/to/config`
4. **Context switching**: `--context my-context`
5. **Direct API server**: `--server https://api-server:6443`
6. **In-cluster**: Automatic when running inside Kubernetes pods

## Output

```bash
$ kubectl podcount_yaron -n kube-system
Running pod count in namespace: kube-system
Running pods: 12

$ kubectl podcount_yaron
Running pod count in all namespaces
Running pods: 45
```

## How It Works

1. **Authentication**: Uses the same authentication mechanisms as `kubectl`
2. **API Calls**: Connects to Kubernetes API server using the official Go client
3. **Filtering**: Lists all pods and filters by `Status.Phase == "Running"`
4. **Namespace Handling**: 
   - Empty namespace = all namespaces
   - Single namespace = specific namespace
   - Comma-separated = multiple namespaces (whitespace trimmed)

## Development

### Project Structure
```
