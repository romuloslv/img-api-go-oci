# img-api-go-oci

Usage

./image -a export $(terraform -chdir=/local/iac/terraform/info output -raw id)  
./image -a import $(terraform -chdir=/local/iac/terraform/info output -raw id)
