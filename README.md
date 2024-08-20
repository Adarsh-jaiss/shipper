# SHIPPER

## Tools

- Azure CLI : for managing AKS with terminal : `curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash`

## Tasks :


### Setting up AKS --> (Then automating) --> DONE!!!

[x] Create a new AKS service with terafform

[x] Find a way to use in terminal

[x] Install dependencies manually 
  - ubuntu
  - curl,wget, homebrew, azure cli
  - JQ
  - shipwright 
  - Build strtegies
  - Tekton

### Setting up build and build run

- Look for a way to parse the user request and add it to the required values in build and build run.
- Fix buildname issue in BuildRun,

### Building API's

- Look for shipwright API or create a new one which will execute this commmand in the cluster : `kubectl create -f ----`

### Setting up infrastructure

- Install Azure-cli

    ```bash
    brew install azure-cli
    az login
    ```

- Install terraform
  
  ```bash
  wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
  echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
  sudo apt update && sudo apt install terraform
  ```

  or 

  ```bash
  brew tap hashicorp/tap
  brew install hashicorp/tap/terraform

  ```


- Now create an **Active directory service principal account** for authentication to Azure. 
   
   ```bash
   az ad sp create-for-rbac --skip-assignment
   ```

**Note** : copy the subscription ID from the azure portal under subscriptions section.


### Accessing the Cluster from your terminal

- install Kubectl 

  ```bash
  az aks install-cli
  ```

- Log in to your Azure account with `az login`

- Set the subscription you're using (if you have multiple):

    ```bash
    az account set --subscription <your-subscription-id>

    az ad sp create-for-rbac --name <service-principal-name> --role Contributor --scopes /subscriptions/<subscription-id>
    ```

- Get the credentials for your AKS cluster (Alredy added in output.tf file)

  ```bash

  az aks get-credentials --resource-group <your-resource-group> --name <your-cluster-name>

  ```

- Then you can copy and run the outputted command. : `terraform output aks_credentials_command`

- Verify the connection :

  ```bash
    kubectl get nodes
    List all pods: kubectl get pods --all-namespaces
    Get cluster info: kubectl cluster-info
    View the configuration: kubectl config view 
  ```






# Things to keep in mind while deployment

- remove the hardcoded directory name from the `apply_k8s_resources.sh` file.