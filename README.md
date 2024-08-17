# SHIPPER

## Tools

- Azure CLI : for managing AKS with terminal : `curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash`

## Tasks :


### Setting up AKS --> (Then automating)

- Create a new AKS service
- Find a way to use in terminal
- Install dependencies manually 
  - shipwright 
  - Build strtegies
  - JQ
  - Tekton

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
    az role assignment create --assignee "c4af0fd7-4ef8-44f2-9bb2-1057c205a6ea" --role Owner --scope /subscriptions/"96bd1c1a-4586-4b61-af0b-b18f3eb919d5"


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